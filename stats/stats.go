package stats

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/contexts"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

type View struct {
	Key   string
	Value int
}

func zcount(c context.Context, key string, days int) ([]int64, error) {
	l := slog.With(slog.String("key", key), slog.Int("days", days))
	l.LogAttrs(c, slog.LevelInfo, "get count statistics")

	ctx := context.Background()
	pipe := db.Redis.Pipeline()
	now := time.Now()

	for i := 0; i < days; i++ {
		t := now.AddDate(0, 0, -i)
		min := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		max := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
		k := key + ":" + t.Format("20060102")

		pipe.ZCount(ctx, k, fmt.Sprintf("%d", min.Unix()), fmt.Sprintf("%d", max.Unix()))
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get statistics", slog.String("key", key), slog.String("error", err.Error()))
		return []int64{}, err
	}

	s := []int64{}
	for _, cmd := range cmds {
		value := cmd.(*redis.IntCmd).Val()
		s = append(s, value)
	}

	slog.LogAttrs(c, slog.LevelInfo, "got statistics")

	return s, nil
}

func zsum(c context.Context, key string, days int) ([]float32, error) {
	l := slog.With(slog.String("key", key), slog.Int("days", days))
	l.LogAttrs(c, slog.LevelInfo, "get sum statistics")

	ctx := context.Background()
	pipe := db.Redis.Pipeline()
	now := time.Now()

	for i := 0; i < days; i++ {
		t := now.AddDate(0, 0, -i)
		min := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		max := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
		k := key + ":" + t.Format("20060102")

		pipe.ZRangeByScore(ctx, k, &redis.ZRangeBy{
			Min:   fmt.Sprintf("%d", min.Unix()),
			Max:   fmt.Sprintf("%d", max.Unix()),
			Count: -1,
		})
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get statistics", slog.String("key", key), slog.String("error", err.Error()))
		return []float32{}, err
	}

	s := []float32{}
	for _, cmd := range cmds {
		value := cmd.(*redis.StringSliceCmd).Val()

		var sum float32 = 0
		for _, v := range value {
			parts := strings.Split(v, ":")
			total, err := strconv.ParseFloat(parts[0], 32)
			if err != nil {
				slog.LogAttrs(c, slog.LevelError, "cannot parse the value", slog.String("value", v), slog.String("error", err.Error()))
				continue
			}

			sum += float32(total)
		}

		s = append(s, sum)
	}

	slog.LogAttrs(c, slog.LevelInfo, "got statistics")

	return s, nil
}

func zadd(c context.Context, key string, data interface{}) error {
	l := slog.With(slog.String("key", key))
	l.LogAttrs(c, slog.LevelInfo, "adding statistics", slog.Any("data", data))
	k := key + ":" + time.Now().Format("20060102")

	if _, err := db.Redis.TxPipelined(context.Background(), func(p redis.Pipeliner) error {
		p.ZAdd(context.Background(), k, redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: data,
		})
		p.Expire(context.Background(), k, conf.StatisticsDuration)

		return nil
	}); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot set revenue statistiscs", slog.String("err", err.Error()))
		return errors.New("error_http_general")
	}

	l.LogAttrs(c, slog.LevelInfo, "statistics added")

	return nil
}

func zgroup(c context.Context, key string, days int) ([]View, error) {
	l := slog.With(slog.String("key", key), slog.Int("days", days))
	l.LogAttrs(c, slog.LevelInfo, "get group statistics")

	ctx := context.Background()
	pipe := db.Redis.Pipeline()
	now := time.Now()

	for i := 0; i < days; i++ {
		t := now.AddDate(0, 0, -i)
		min := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		max := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
		k := key + ":" + t.Format("20060102")

		pipe.ZRangeByScore(ctx, k, &redis.ZRangeBy{
			Min:   fmt.Sprintf("%d", min.Unix()),
			Max:   fmt.Sprintf("%d", max.Unix()),
			Count: -1,
		})
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get statistics", slog.String("key", key), slog.String("error", err.Error()))
		return []View{}, err
	}

	p := map[string]int{}
	for _, cmd := range cmds {
		value := cmd.(*redis.StringSliceCmd).Val()
		for _, v := range value {
			parts := strings.Split(v, ":")

			p[parts[0]]++
		}
	}

	keys := make([]string, 0, len(p))
	for key := range p {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return p[keys[i]] > p[keys[j]]
	})

	views := []View{}
	for _, key := range keys {
		views = append(views, View{
			Key:   key,
			Value: p[key],
		})
	}

	slog.LogAttrs(c, slog.LevelInfo, "got statistics")

	return views, nil
}

func Visit(c context.Context) error {
	rid := c.Value(middleware.RequestIDKey).(string)
	if err := zadd(c, "statvisits", rid); err != nil {
		return err
	}

	cid := c.Value(contexts.Cart).(string)
	return zadd(c, "statuniquevisits", cid)
}

func Visits(c context.Context) ([]int64, error) {
	return zcount(c, "statvisits", 30)
}

func UniqueVisits(c context.Context) ([]int64, error) {
	return zcount(c, "statuniquevisits", 30)
}

func Users(c context.Context, days int) ([]int64, error) {
	return zcount(c, "statnewusers", days)
}

func NewUser(c context.Context, id int64) error {
	return zadd(c, "statnewusers", id)
}

func ActiveUsers(c context.Context, days int) ([]int64, error) {
	return zcount(c, "statactiveusers", days)
}

func ActiveUser(c context.Context, id string) error {
	return zadd(c, "statactiveusers", id)
}

func Order(c context.Context, id string) error {
	return zadd(c, "statorders", id)
}

func Orders(c context.Context, days int) ([]int64, error) {
	return zcount(c, "statorders", days)
}

func Revenue(c context.Context, oid string, total float32) error {
	return zadd(c, "statrevenues", fmt.Sprintf("%f:%s", total, oid))
}

func Revenues(c context.Context, days int) ([]float32, error) {
	return zsum(c, "statrevenues", days)
}

func SoldProduct(c context.Context, oid, id string, quantity int) error {
	for i := 0; i < quantity; i++ {
		err := zadd(c, "statsoldproducts", fmt.Sprintf("%s:%d:%s", id, quantity, oid))
		if err != nil {
			return err
		}
	}

	return nil
}

func SoldProducts(c context.Context, days int) ([]int64, error) {
	return zcount(c, "statsoldproducts", days)
}

func ProductView(c context.Context, id string) error {
	rid := c.Value(middleware.RequestIDKey).(string)
	return zadd(c, "statvisitproduct", id+":"+rid)
}

func ProductUniqueView(c context.Context, id string) error {
	cid := c.Value(contexts.Cart).(string)
	return zadd(c, "statuniquevisitproduct", id+":"+cid)
}

func ProductViews(c context.Context, days int) ([]View, error) {
	return zgroup(c, "statvisitproduct", days)
}
