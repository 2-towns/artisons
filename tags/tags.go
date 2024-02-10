package tags

import (
	"artisons/db"
	"artisons/tags/tree"
	"artisons/validators"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type Tag struct {
	Key       string `validate:"required,alphanum"`
	Label     string `validate:"required"`
	Image     string
	Children  []string
	Root      bool
	Score     int
	UpdatedAt time.Time
}

type ListResults struct {
	Total int
	Tags  []Tag
}

func (p Tag) Validate(ctx context.Context) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "validating a tag")

	if err := validators.V.Struct(p); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot validate the tag", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input:%s", low)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "tag validated")

	return nil
}

func Exists(ctx context.Context, key string) (bool, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "checking existence", slog.String("key", key))

	exists, err := db.Redis.Exists(ctx, "tag:"+key).Result()

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot check tags existence")
		return false, errors.New("something went wrong")
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "tag existence", slog.String("key", key), slog.Int64("exists", exists))

	return exists > 0, nil
}

func AreEligible(ctx context.Context, keys []string) (bool, error) {
	l := slog.With(slog.Any("tag", keys))
	l.LogAttrs(ctx, slog.LevelInfo, "looking if keys are root tags")

	roots, err := db.Redis.ZRange(ctx, "tags:root", 0, 9999).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the root tags", slog.String("error", err.Error()))
		return false, errors.New("something went wrong")
	}

	for _, val := range keys {
		if slices.Contains(roots, val) {
			slog.LogAttrs(ctx, slog.LevelInfo, "the tag is not eligible", slog.String("tag", val))
			return false, nil
		}
	}

	return true, nil

}

func (t Tag) Save(ctx context.Context) (string, error) {
	l := slog.With(slog.String("tag", t.Key))
	l.LogAttrs(ctx, slog.LevelInfo, "adding a new tag")

	children := strings.Join(t.Children, ";")
	now := time.Now()

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, "tag:"+t.Key,
			"image", t.Image,
			"children", children,
			"label", t.Label,
			"updated_at", now.Unix(),
		)

		rdb.HSetNX(ctx, "tag:"+t.Key, "key", t.Key)

		rdb.ZAdd(ctx, "tags", redis.Z{
			Score:  float64(t.UpdatedAt.Unix()),
			Member: t.Key,
		})

		if t.Root {
			rdb.ZAdd(ctx, "tags:root", redis.Z{
				Score:  float64(t.Score),
				Member: t.Key,
			})
		}

		tree.Build(ctx)

		return nil

	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "tag saved successfully")

	return t.Key, nil
}

func parse(ctx context.Context, data map[string]string) (Tag, error) {
	var updatedAt int64 = 0
	var err error

	if data["created_at"] != "" {
		updatedAt, err = strconv.ParseInt(data["created_at"], 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the updated at", slog.String("error", err.Error()), slog.String("updated_at", data["updated_at"]))
		}
	}

	tag := Tag{
		Key:       data["key"],
		Label:     data["label"],
		Image:     data["image"],
		Children:  strings.Split(data["children"], ";"),
		UpdatedAt: time.Unix(updatedAt, 0),
	}

	return tag, nil
}

func Find(ctx context.Context, key string) (Tag, error) {
	l := slog.With(slog.String("id", key))
	l.LogAttrs(ctx, slog.LevelInfo, "looking for tag")

	if key == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate empty tag key")
		return Tag{}, errors.New("input:id")
	}

	if exists, err := db.Redis.Exists(ctx, "tag:"+key).Result(); exists == 0 || err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot find the tag")
		return Tag{}, errors.New("oops the data is not found")
	}

	data, err := db.Redis.HGetAll(ctx, "tag:"+key).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot find the tag", slog.String("error", err.Error()))
		return Tag{}, err
	}

	tag, err := parse(ctx, data)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot parse the tag", slog.String("error", err.Error()))
		return Tag{}, err
	}

	score, err := db.Redis.ZScore(ctx, "tags:root", key).Result()
	if err == nil {
		tag.Root = true
		tag.Score = int(score)
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the tag is found")

	return tag, nil
}

func List(ctx context.Context, offset, num int) (ListResults, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "listing tags")

	keys, err := db.Redis.ZRevRange(ctx, "tags", int64(offset), int64(num)).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the tags", slog.String("error", err.Error()))
		return ListResults{}, errors.New("something went wrong")
	}

	roots, err := db.Redis.ZRange(ctx, "tags:root", 0, 9999).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the root tags", slog.String("error", err.Error()))
		return ListResults{}, errors.New("something went wrong")
	}

	tags := []Tag{}

	cmds, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for _, val := range keys {
			rdb.HGetAll(ctx, "tag:"+val)
		}

		return nil
	})

	if err != nil && err.Error() != "redis: nil" {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the tag meta data", slog.String("error", err.Error()))
		return ListResults{}, errors.New("something went wrong")
	}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil && cmd.Err().Error() != "redis: nil" {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the tag meta data", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		val := cmd.(*redis.MapStringStringCmd).Val()

		tag, err := parse(ctx, val)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the tag", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		tag.Root = slices.Contains(roots, tag.Key)

		tags = append(tags, tag)
	}

	total, err := db.Redis.ZCount(ctx, "tags", "-inf", "+inf").Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the tags count")
		return ListResults{}, errors.New("something went wrong")
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "found tags", slog.Int("length", len(tags)))

	return ListResults{
		Total: int(total),
		Tags:  tags,
	}, nil
}

func Delete(ctx context.Context, key string) error {
	l := slog.With(slog.String("tag", key))
	l.LogAttrs(ctx, slog.LevelInfo, "deleting tag")

	if key == "" {
		slog.LogAttrs(ctx, slog.LevelInfo, "the key cannot be empty")
		return errors.New("input:key")
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HDel(ctx, "tag", key)
		rdb.Del(ctx, "tag:"+key)
		rdb.ZRem(ctx, "tags", key)

		return nil

	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot delete the data", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "tag deleted successfully")

	return nil
}
