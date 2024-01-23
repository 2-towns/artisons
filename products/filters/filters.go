package filters

import (
	"context"
	"errors"
	"fmt"
	"gifthub/db"
	"gifthub/validators"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type Filter struct {
	Key      string `validate:"required,alphanum"`
	Editable bool
	Label    string `validate:"required"`
	Score    int
	Values   []string
	Active   bool

	UpdatedAt time.Time
}

type ListResults struct {
	Total   int
	Filters []Filter
}

func (f Filter) Validate(ctx context.Context) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "validating a filter")

	if err := validators.V.Struct(f); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot validate the filter", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input:%s", low)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "filter validated")

	return nil
}

func Exists(ctx context.Context, key string) (bool, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "checking existence", slog.String("key", key))

	exists, err := db.Redis.Exists(ctx, "filter:"+key).Result()

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot check filter existence")
		return false, errors.New("something went wrong")
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "filter existence", slog.String("key", key), slog.Int64("exists", exists))

	return exists > 0, nil
}

func Editable(ctx context.Context, key string) (bool, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "checking editable", slog.String("key", key))

	editable, err := db.Redis.HGet(ctx, "filter:"+key, "editable").Result()

	if err != nil && err.Error() != "redis: nil" {
		slog.LogAttrs(ctx, slog.LevelError, "cannot check filter editable")
		return false, errors.New("something went wrong")
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "filter is editable", slog.Bool("eligible", editable == "1"))

	return editable == "1", nil
}

func (f Filter) Save(ctx context.Context) (string, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "creating a filter")

	now := time.Now().Unix()

	editable, err := Editable(ctx, f.Key)
	if err != nil {
		return "", err
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		k := fmt.Sprintf("filter:%s", f.Key)
		rdb.HSet(ctx, k,
			"label", f.Label,
			"updated_at", now,
		)

		if editable {
			rdb.HSet(ctx, k, "values", strings.Join(f.Values, ";"))
		}

		rdb.HSetNX(ctx, k, "key", f.Key)
		rdb.HSetNX(ctx, k, "editable", "1")

		rdb.ZAdd(ctx, "filters", redis.Z{
			Score:  float64(now),
			Member: f.Key,
		})

		if f.Active {
			rdb.ZAdd(ctx, "filters:active", redis.Z{
				Score:  float64(f.Score),
				Member: f.Key,
			})
		} else {
			rdb.ZRem(ctx, "filters:active", f.Key)
		}

		return nil
	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "filter created", slog.String("key", f.Key))

	return f.Key, nil
}

func parse(ctx context.Context, data map[string]string) (Filter, error) {
	updatedAt, err := strconv.ParseInt(data["updated_at"], 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the filter updated at", slog.String("updated_at", data["updated_at"]))
		return Filter{}, errors.New("input:updated_at")
	}

	filter := Filter{
		Key:       data["key"],
		Editable:  data["editable"] == "1",
		Label:     data["label"],
		Values:    strings.Split(data["values"], ";"),
		UpdatedAt: time.Unix(updatedAt, 0),
	}

	return filter, nil
}

func Actives(ctx context.Context) ([]Filter, error) {
	l := slog.With()
	l.LogAttrs(ctx, slog.LevelInfo, "looking for filters")

	actives, err := db.Redis.ZRevRange(ctx, "filters:active", 0, 9999).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the active filter keys")
		return []Filter{}, errors.New("something went wrong")
	}

	cmds, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for _, k := range actives {
			rdb.HGetAll(ctx, "filter:"+k)
		}

		return nil
	})

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the filters")
		return []Filter{}, errors.New("something went wrong")
	}

	filters := []Filter{}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the filter", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		val := cmd.(*redis.MapStringStringCmd).Val()

		filter, err := parse(ctx, val)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the filter", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		filter.Active = slices.Contains(actives, filter.Key)

		filters = append(filters, filter)
	}

	return filters, nil
}

func List(ctx context.Context, offset, num int) (ListResults, error) {
	l := slog.With()
	l.LogAttrs(ctx, slog.LevelInfo, "looking for filters")

	keys, err := db.Redis.ZRevRange(ctx, "filters", int64(offset), int64(num)).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the filter keys")
		return ListResults{}, errors.New("something went wrong")
	}

	actives, err := db.Redis.ZRange(ctx, "filters:active", 0, 9999).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the active filter keys")
		return ListResults{}, errors.New("something went wrong")
	}

	cmds, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for _, k := range keys {
			rdb.HGetAll(ctx, "filter:"+k)
		}

		return nil
	})

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the filters")
		return ListResults{}, errors.New("something went wrong")
	}

	filters := []Filter{}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the filter", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		val := cmd.(*redis.MapStringStringCmd).Val()

		filter, err := parse(ctx, val)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the filter", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		filter.Active = slices.Contains(actives, filter.Key)

		filters = append(filters, filter)
	}

	total, err := db.Redis.ZCount(ctx, "filters", "-inf", "+inf").Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the filters count")
		return ListResults{}, errors.New("something went wrong")
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "found filters", slog.Int("length", len(filters)))

	return ListResults{
		Total:   int(total),
		Filters: filters,
	}, nil
}

func Find(ctx context.Context, key string) (Filter, error) {
	l := slog.With(slog.String("id", key))
	l.LogAttrs(ctx, slog.LevelInfo, "looking for filter")

	if key == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate empty filter key")
		return Filter{}, errors.New("input:id")
	}

	if exists, err := db.Redis.Exists(ctx, "filter:"+key).Result(); exists == 0 || err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot find the filter")
		return Filter{}, errors.New("oops the data is not found")
	}

	data, err := db.Redis.HGetAll(ctx, "filter:"+key).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot find the filter", slog.String("error", err.Error()))
		return Filter{}, err
	}

	filter, err := parse(ctx, data)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot parse the filter", slog.String("error", err.Error()))
		return Filter{}, err
	}

	score, err := db.Redis.ZScore(ctx, "filters:active", key).Result()
	if err == nil {
		filter.Active = true
		filter.Score = int(score)
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the filter is found")

	return filter, nil
}

func Delete(ctx context.Context, key string) error {
	l := slog.With(slog.String("filter", key))
	l.LogAttrs(ctx, slog.LevelInfo, "deleting filter")

	if key == "" {
		slog.LogAttrs(ctx, slog.LevelInfo, "the key cannot be empty")
		return errors.New("input:key")
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Del(ctx, "filter:"+key)
		rdb.ZRem(ctx, "filters", key)

		return nil

	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot delete the data", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "filter deleted successfully")

	return nil
}
