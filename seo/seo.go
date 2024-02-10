package seo

import (
	"artisons/db"
	"artisons/seo/urls"
	"artisons/validators"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type Translation struct {
	URL         string
	Title       string
	Description string
}

type Content struct {
	Key         string `validate:"required"`
	URL         string `validate:"required"`
	Title       string `validate:"required"`
	Description string `validate:"required"`
}

type ListResults struct {
	Total   int
	Content []Content
}

func (c Content) Validate(ctx context.Context) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "validating a seo")

	if err := validators.V.Struct(c); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot validate the translation", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input:%s", low)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "seo validated")

	return nil
}

func (c Content) Save(ctx context.Context) (string, error) {
	l := slog.With(slog.String("key", c.Key))
	l.LogAttrs(ctx, slog.LevelInfo, "saving a seo")

	key := "seo:" + c.Key
	now := time.Now()

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, key,
			"title", db.Escape(c.Title),
			"description", db.Escape(c.Description),
			"url", db.Escape(c.URL),
			"updated_at", now.Unix(),
		)

		rdb.SAdd(ctx, "seo", c.Key)

		return nil
	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot store the seo", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	urls.Set(c.Key, "title", c.Title)
	urls.Set(c.Key, "description", c.Description)
	urls.Set(c.Key, "url", c.URL)

	l.LogAttrs(ctx, slog.LevelInfo, "seo saved")

	return c.Key, nil
}

func Find(ctx context.Context, key string) (Content, error) {
	l := slog.With(slog.String("key", key))
	l.LogAttrs(ctx, slog.LevelInfo, "looking for seo")

	if key == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate empty seo id")
		return Content{}, errors.New("input:id")
	}

	c, err := db.Redis.HGetAll(ctx, key).Result()
	if err != nil || len(c) == 0 {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot get the seo", slog.Any("error", err))
		return Content{}, errors.New("the data is not found")
	}

	return Content{
		Key:         c["key"],
		URL:         db.Unescape(c["url"]),
		Title:       db.Unescape(c["title"]),
		Description: db.Unescape(c["description"]),
	}, nil
}

func List(ctx context.Context, offset, num int) (ListResults, error) {
	keys, err := db.Redis.SMembers(ctx, "seo").Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the seo keys", slog.String("error", err.Error()))
		return ListResults{}, errors.New("something went wrong")
	}

	cmds, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for _, key := range keys {
			key := "seo:" + key
			rdb.HGetAll(ctx, key)
		}

		return nil
	})

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the seo", slog.String("error", err.Error()))
		return ListResults{}, errors.New("something went wrong")
	}

	content := []Content{}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the seo", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		val := cmd.(*redis.MapStringStringCmd).Val()

		c := Content{
			Key:         val["key"],
			URL:         db.Unescape(val["url"]),
			Title:       db.Unescape(val["title"]),
			Description: db.Unescape(val["description"]),
		}

		content = append(content, c)
	}

	o := math.Min(float64(offset), float64(len(content)))
	n := math.Min(float64(num), float64(len(content)))

	return ListResults{
		Total:   len(content),
		Content: content[int(o):int(n)],
	}, nil
}
