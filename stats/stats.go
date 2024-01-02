package stats

import (
	"cmp"
	"context"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/contexts"
	"gifthub/http/referer"
	"gifthub/string/stringutil"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mileusna/useragent"
	"github.com/redis/go-redis/v9"
	"golang.org/x/text/language"
)

type MostValue struct {
	Key     string
	Value   float64
	Percent float64
	// Label   string
	URL  string
	Lang language.Tag
}

type Count struct {
	Value []int64
	Sum   int64
	// Label string
}

type Data []Count

type VisitData struct {
	URL     string
	Referer string
}

// MostValues returns the most values statistics.
// The keys available are:
// - stats:pageviews - the most visited pages
// - stats:browsers - the most used browsers
// - stats:referers - the most used referers
// - stats:systems - the most used systems
// - stats:products:most - the most sold products
// - stats:products:shared - the most shared products
// For each statistics keys, a subset of keys is generated to retrieve the data
// for the specified days interval. So if the days are 7, 7 keys will be added to the subset:
// stats:pageviews:20060102, stats:pageviews:20060103 ...
// The most values are returned by Redis in using ZUnionWithScores.
// Then the sum is calculated, the results are ordered and the percent representation is added
// for each entry.
// The returned value is an array containing the most values described above, respecting the
// available keys order. So stats:pageviews is the index 0, stats:browsers 1...
// The product titles of the most sold products and most share are loaded from redis
// at the end of the function.
func MostValues(c context.Context, days int) ([][]MostValue, error) {
	prefix := ""
	demo, ok := c.Value(contexts.Demo).(bool)
	if demo && ok {
		prefix = "demo:"
	}

	keys := []string{
		prefix + "stats:pageviews",
		prefix + "stats:browsers",
		prefix + "stats:referers",
		prefix + "stats:systems",
		prefix + "stats:products:most",
		prefix + "stats:products:shared",
	}

	l := slog.With(slog.Any("key", keys), slog.Int("days", days), slog.Bool("demo", demo))
	l.LogAttrs(c, slog.LevelInfo, "get range statistics")

	now := time.Now()
	ctx := context.Background()
	pipe := db.Redis.Pipeline()

	for _, key := range keys {
		cKeys := []string{}

		for i := 0; i < days; i++ {
			t := now.AddDate(0, 0, -i)
			k := key + ":" + t.Format("20060102")
			cKeys = append(cKeys, k)
		}

		pipe.ZUnionWithScores(ctx, redis.ZStore{
			Keys: cKeys,
		})
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil && err.Error() != "redis: nil" {
		slog.LogAttrs(c, slog.LevelError, "cannot get statistics", slog.String("error", err.Error()))
		return [][]MostValue{}, err
	}

	values := [][]MostValue{}
	pids := []string{}

	for _, cmd := range cmds {
		var total float64 = 0
		cur := []MostValue{}

		v := cmd.(*redis.ZSliceCmd).Val()

		for _, val := range v {
			total += val.Score
		}

		slices.SortFunc(v, func(a, b redis.Z) int {
			return cmp.Compare(b.Score, a.Score)
		})

		key := cmd.(*redis.ZSliceCmd).Args()[2].(string)

		for i, value := range v {
			val := value.Member.(string)
			cur = append(cur, MostValue{
				Key:     val,
				Value:   value.Score,
				Percent: value.Score / total * 100,
				URL:     val,
			})

			if i < conf.DashboardMostItems && strings.HasPrefix(key, prefix+"stats:products") {
				pids = append(pids, val)
			}
		}

		values = append(values, cur)
	}

	pipe = db.Redis.Pipeline()
	for _, pid := range pids {
		pipe.HGet(ctx, "product:"+pid, "title")
	}

	cmds, err = pipe.Exec(ctx)
	if err != nil && err.Error() != "redis: nil" {
		slog.LogAttrs(c, slog.LevelError, "cannot get the product names", slog.String("error", err.Error()))
		return values, nil
	}

	pnames := map[string]string{}

	for idx, cmd := range cmds {
		v := cmd.(*redis.StringCmd).Val()
		pid := pids[idx]
		pnames[pid] = v
	}

	for i := 0; i < conf.DashboardMostItems; i++ {
		if len(values[4]) > i {
			val := values[4][i]
			name := db.Unescape(pnames[val.Key])
			values[4][i].Key = name
			slug := stringutil.Slugify(name)
			values[4][i].URL = fmt.Sprintf("%s-%s.html", val.Key, slug)
		}

		if len(values[5]) > i {
			val := values[5][i]
			name := db.Unescape(pnames[val.Key])
			values[5][i].Key = name
			slug := stringutil.Slugify(name)
			values[5][i].URL = fmt.Sprintf("%s-%s.html", val.Key, slug)
		}
	}

	return values, nil
}

func sum(array []int64) int64 {
	var result int64 = 0

	for _, v := range array {
		result += v
	}
	return result
}

// GetAll returns the count statistics values.
// The keys available are:
// - stats:visits - the visits
// - stats:visits:unique - the unique visits
// - stats:pageview - the page views
// - stats:orders - the order total amount
// - stats:orders:count - the order total count
// For each statistics keys, a subset of keys is generated to retrieve the data
// for the specified days interval.
// To use only one loop, the stop number is the multiplication result between the total of keys
// and the days. Then the current row is depending of the the dividend of this operation, which is
// also used as the current indice. The result are keys like this :
// stats:visits:20060102, stats:visits:20060103 ...
// If the value does not exist, 0 will be added.
// For each statistics, the sum is calculated and added to the struct.
// Finalylly a bounce rate between the visits and the unique visits is calculated and added
// in the same format than the other statistics.
func GetAll(c context.Context, days int) (Data, error) {
	prefix := ""
	demo, ok := c.Value(contexts.Demo).(bool)
	if demo && ok {
		prefix = "demo:"
	}

	l := slog.With(slog.Int("days", days), slog.Bool("demo", demo))
	l.LogAttrs(c, slog.LevelInfo, "get all statistics")

	ctx := context.Background()
	now := time.Now()
	pipe := db.Redis.Pipeline()
	keys := []string{
		prefix + "stats:visits:",
		prefix + "stats:visits:unique:",
		prefix + "stats:pageviews:all:",
		prefix + "stats:orders:revenues:",
		prefix + "stats:orders:count:",
	}
	row := 0

	for i := 0; i < days*len(keys); i++ {
		if i > 0 && i%days == 0 {
			row++
		}

		t := now.AddDate(0, 0, -(i % days))
		k := keys[row] + t.Format("20060102")
		pipe.Get(ctx, k)
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil && err.Error() != "redis: nil" {
		slog.LogAttrs(c, slog.LevelError, "cannot get statistics", slog.String("error", err.Error()))
		return Data{}, err
	}

	values := Data{
		{Value: []int64{}},
		{Value: []int64{}},
		{Value: []int64{}},
		{Value: []int64{}},
		{Value: []int64{}},
		{Value: []int64{}},
	}
	row = 0

	for i, cmd := range cmds {
		if i > 0 && i%days == 0 {
			values[row].Sum = sum(values[row].Value)
			row++
		}

		s := cmd.(*redis.StringCmd).Val()

		if s == "" {
			values[row].Value = append([]int64{0}, values[row].Value...)
			continue
		}

		f, err := strconv.ParseFloat(s, 64)

		if err != nil {
			slog.LogAttrs(c, slog.LevelError, "cannot convert the value to int", slog.String("value", s), slog.String("error", err.Error()))
			values[row].Value = append([]int64{0}, values[row].Value...)
		} else {
			value := int64(f)
			values[row].Value = append([]int64{value}, values[row].Value...)
		}
	}

	values[row].Sum = sum(values[row].Value)

	for i := 0; i < days; i++ {
		if values[1].Value[i] == 0 {
			values[len(keys)].Value = append(values[len(keys)].Value, 0)
		} else {
			p := float64(values[1].Value[i]) / float64(values[0].Value[i]) * 100
			values[len(keys)].Value = append(values[len(keys)].Value, int64(p))
		}
	}

	if values[1].Sum == 0 {
		values[len(keys)].Sum = 0
	} else {
		values[len(keys)].Sum = int64(float64(values[1].Sum) / float64(values[0].Sum) * 100)
	}

	slog.LogAttrs(c, slog.LevelInfo, "got statistics")

	return values, nil
}

func Visit(c context.Context, ua useragent.UserAgent, data VisitData) error {
	l := slog.With()
	l.LogAttrs(c, slog.LevelInfo, "get range statistics")

	ctx := context.Background()
	now := time.Now().Format("20060102")

	cid := c.Value(contexts.Cart).(string)
	hasVisited, err := db.Redis.SIsMember(ctx, "stats:visits:members:"+now, cid).Result()
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get is member", slog.String("error", err.Error()))
		return err
	}

	pipe := db.Redis.Pipeline()

	r := referer.Parse(data.Referer)

	if r.Referer != "" {
		pipe.Incr(ctx, "stats:visits:"+now)

		l.LogAttrs(ctx, slog.LevelInfo, "visits stat added to pipe")

		r := referer.Parse(data.Referer)
		pipe.ZIncrBy(ctx, "stats:referers:"+now, 1, r.Referer)
		l.LogAttrs(ctx, slog.LevelInfo, "referer stat added to pipe", slog.String("referer", r.Referer))
	}

	pipe.ZIncrBy(ctx, "stats:pageviews:"+now, 1, data.URL)
	l.LogAttrs(ctx, slog.LevelInfo, "pageviews stat added to pipe", slog.String("url", data.URL))

	if ua.Name != "" {
		pipe.ZIncrBy(ctx, "stats:browsers:"+now, 1, ua.Name)
		l.LogAttrs(ctx, slog.LevelInfo, "browsers stat added to pipe", slog.String("browser", ua.Name))
	}

	if ua.OS != "" {
		pipe.ZIncrBy(ctx, "stats:systems:"+now, 1, ua.Name)
		l.LogAttrs(ctx, slog.LevelInfo, "systems stat added to pipe", slog.String("system", ua.OS))
	}

	if !hasVisited {
		pipe.Incr(ctx, "stats:visits:unique:"+now)
		pipe.SAdd(ctx, "stats:visits:members:"+now, cid)
		pipe.Expire(ctx, "stats:visits:members:"+now, time.Hour*24)
		l.LogAttrs(ctx, slog.LevelInfo, "unique visite stat added to pipe", slog.String("cid", cid))
	}

	pipe.Incr(ctx, "stats:pageviews:all:"+now)
	l.LogAttrs(ctx, slog.LevelInfo, "pageview stats added to pipe")

	_, err = pipe.Exec(ctx)
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot add statistics", slog.String("error", err.Error()))
	}

	return nil
}

func Order(c context.Context, id string, quantites map[string]int, total float64) error {
	l := slog.With(slog.String("id", id), slog.Float64("total", total))
	l.LogAttrs(c, slog.LevelInfo, "store order statistics")

	ctx := context.Background()
	now := time.Now().Format("20060102")
	pipe := db.Redis.Pipeline()

	for pid, q := range quantites {
		pipe.ZIncrBy(ctx, "stats:products:most:"+now, float64(q), pid)
	}

	pipe.IncrByFloat(ctx, "stats:orders:revenues:"+now, total)
	pipe.Incr(ctx, "stats:orders:count:"+now)

	_, err := pipe.Exec(ctx)
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot add statistics", slog.String("error", err.Error()))
	}

	return nil
}
