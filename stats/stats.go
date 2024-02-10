package stats

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/http/referer"
	"artisons/string/stringutil"
	"artisons/users"
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
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
	Value []int
	Sum   int
}

type Data []Count

type VisitData struct {
	URL     string
	Referer string
}

func generateDemoData(ctx context.Context) error {
	pipe := db.Redis.Pipeline()
	now := time.Now()
	referers := []string{"Google", "Unknown", "Yandex", "DuckDuckGo"}
	products := []string{"PDT1", "PDT2", "PDT3", "PDT4", "PDT5"}
	browsers := []string{"Chrome", "Safari", "Firefox", "Edge"}
	systems := []string{"Windows", "Android", "iOS", "Linux"}
	urls := []string{"/", "/super-article-du-blog", "/PDT2-sweat-a-capuche-uniforme", "/cgv", "/panier", "/coucou"}

	for i := 0; i < 30; i++ {
		i := rand.Intn(30)
		score := now.AddDate(0, 0, -i)
		exists, err := db.Redis.Exists(ctx, "demo:stats:pageviews:"+score.Format("20060102")).Result()
		if exists > 0 && err == nil {
			slog.LogAttrs(ctx, slog.LevelInfo, "the demo stat exists")
			continue
		}

		for count := 0; count < 10; count++ {

			urli := rand.Intn(len(urls))
			slug := urls[urli]
			visits := rand.Intn(1000) + 100
			amount := rand.Intn(5000) + 100
			count := rand.Intn(10) + 1
			uniques := rand.Intn(visits)
			rrand := rand.Intn(len(referers))
			brand := rand.Intn(len(browsers))
			srand := rand.Intn(len(systems))
			prand := rand.Intn(len(products))

			pipe.ZIncrBy(ctx, "demo:stats:pageviews:"+score.Format("20060102"), 1, slug)
			pipe.ZIncrBy(ctx, "demo:stats:products:most:"+score.Format("20060102"), 1, products[prand])
			pipe.ZIncrBy(ctx, "demo:stats:products:shared:"+score.Format("20060102"), 1, products[prand])
			pipe.ZIncrBy(ctx, "demo:stats:browsers:"+score.Format("20060102"), 1, browsers[brand])
			pipe.ZIncrBy(ctx, "demo:stats:referers:"+score.Format("20060102"), 1, referers[rrand])
			pipe.ZIncrBy(ctx, "demo:stats:systems:"+score.Format("20060102"), 1, systems[srand])
			pipe.Set(ctx, "demo:stats:visits:"+score.Format("20060102"), visits, 0)
			pipe.Set(ctx, "demo:stats:orders:revenues:"+score.Format("20060102"), amount, 0)
			pipe.Set(ctx, "demo:stats:orders:count:"+score.Format("20060102"), count, 0)
			pipe.Set(ctx, "demo:stats:visits:unique:"+score.Format("20060102"), uniques, 0)
			pipe.Set(ctx, "demo:stats:pageviews:all:"+score.Format("20060102"), visits*2, 0)
		}
	}

	_, err := pipe.Exec(ctx)

	return err
}

func getPrefix(ctx context.Context) string {
	u, ok := ctx.Value(contexts.User).(users.User)

	if ok && u.Demo {
		slog.LogAttrs(ctx, slog.LevelInfo, "demo is activated")
		return "demo:"
	}

	return ""

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
func MostValues(ctx context.Context, days int) ([][]MostValue, error) {
	prefix := getPrefix(ctx)
	values := [][]MostValue{}

	if prefix == "demo:" {
		if err := generateDemoData(ctx); err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get most values", slog.String("error", err.Error()))
			return values, errors.New("something went wrong")
		}
	}

	keys := []string{
		prefix + "stats:pageviews",
		prefix + "stats:browsers",
		prefix + "stats:referers",
		prefix + "stats:systems",
		prefix + "stats:products:most",
		prefix + "stats:products:shared",
	}

	l := slog.With(slog.Any("key", keys), slog.Int("days", days))
	l.LogAttrs(ctx, slog.LevelInfo, "get range statistics")

	now := time.Now()
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
		slog.LogAttrs(ctx, slog.LevelError, "cannot get statistics", slog.String("error", err.Error()))
		return [][]MostValue{}, err
	}

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
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the product names", slog.String("error", err.Error()))
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

func sum(array []int) int {
	result := 0

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
func GetAll(ctx context.Context, days int) (Data, error) {
	values := Data{
		{Value: []int{}},
		{Value: []int{}},
		{Value: []int{}},
		{Value: []int{}},
		{Value: []int{}},
		{Value: []int{}},
	}
	prefix := getPrefix(ctx)
	if prefix == "demo:" {
		if err := generateDemoData(ctx); err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get most values", slog.String("error", err.Error()))
			return values, errors.New("something went wrong")
		}
	}

	l := slog.With(slog.Int("days", days))
	l.LogAttrs(ctx, slog.LevelInfo, "get all statistics")

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
		slog.LogAttrs(ctx, slog.LevelError, "cannot get statistics", slog.String("error", err.Error()))
		return Data{}, err
	}

	row = 0

	for i, cmd := range cmds {
		if i > 0 && i%days == 0 {
			values[row].Sum = sum(values[row].Value)
			row++
		}

		s := cmd.(*redis.StringCmd).Val()

		if s == "" {
			values[row].Value = append([]int{0}, values[row].Value...)
			continue
		}

		f, err := strconv.ParseFloat(s, 64)

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot convert the value to int", slog.String("value", s), slog.String("error", err.Error()))
			values[row].Value = append([]int{0}, values[row].Value...)
		} else {
			value := int(f)
			values[row].Value = append([]int{value}, values[row].Value...)
		}
	}

	values[row].Sum = sum(values[row].Value)

	for i := 0; i < days; i++ {
		if values[1].Value[i] == 0 {
			values[len(keys)].Value = append(values[len(keys)].Value, 0)
		} else {
			p := float64(values[1].Value[i]) / float64(values[0].Value[i]) * 100
			values[len(keys)].Value = append(values[len(keys)].Value, int(p))
		}
	}

	if values[1].Sum == 0 {
		values[len(keys)].Sum = 0
	} else {
		values[len(keys)].Sum = int(float64(values[1].Sum) / float64(values[0].Sum) * 100)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "got statistics")

	return values, nil
}

func Visit(ctx context.Context, ua useragent.UserAgent, data VisitData) error {
	now := time.Now().Format("20060102")
	did := ctx.Value(contexts.Device).(string)
	hasVisited, err := db.Redis.SIsMember(ctx, "stats:visits:members:"+now, did).Result()

	if err != nil && err.Error() != "redis: nil" {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get is member", slog.String("error", err.Error()))
		return err
	}

	pipe := db.Redis.Pipeline()

	r := referer.Parse(data.Referer)

	prefix := getPrefix(ctx)

	if r.Referer != "" {
		pipe.Incr(ctx, prefix+"stats:visits:"+now)

		r := referer.Parse(data.Referer)
		pipe.ZIncrBy(ctx, prefix+"stats:referers:"+now, 1, r.Referer)
	}

	pipe.ZIncrBy(ctx, prefix+"stats:pageviews:"+now, 1, data.URL)

	if ua.Name != "" {
		pipe.ZIncrBy(ctx, prefix+"stats:browsers:"+now, 1, ua.Name)
	}

	if ua.OS != "" {
		pipe.ZIncrBy(ctx, prefix+"stats:systems:"+now, 1, ua.Name)
	}

	if !hasVisited {
		pipe.Incr(ctx, prefix+"stats:visits:unique:"+now)
		pipe.SAdd(ctx, prefix+"stats:visits:members:"+now, did)
		pipe.Expire(ctx, prefix+"stats:visits:members:"+now, time.Hour*24)
	}

	pipe.Incr(ctx, prefix+"stats:pageviews:all:"+now)

	_, err = pipe.Exec(ctx)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot add statistics", slog.String("error", err.Error()))
	} else {
		slog.Log(ctx, slog.LevelInfo, "visit stat added to redis")
	}

	return nil
}

func Order(ctx context.Context, id string, quantites map[string]int, total float64) error {
	l := slog.With(slog.String("id", id), slog.Float64("total", total))
	l.LogAttrs(ctx, slog.LevelInfo, "store order statistics")

	now := time.Now().Format("20060102")
	pipe := db.Redis.Pipeline()

	for pid, q := range quantites {
		pipe.ZIncrBy(ctx, "stats:products:most:"+now, float64(q), pid)
	}

	pipe.IncrByFloat(ctx, "stats:orders:revenues:"+now, total)
	pipe.Incr(ctx, "stats:orders:count:"+now)

	_, err := pipe.Exec(ctx)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot add statistics", slog.String("error", err.Error()))
	}

	return nil
}
