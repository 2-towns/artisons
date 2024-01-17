package seo

import (
	"context"
	"errors"
	"fmt"
	"gifthub/db"
	"gifthub/http/router"
	"gifthub/validators"
	"log"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"

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
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SearchResults struct {
	Total   int
	Content []Content
}

var URLs map[string]Content = map[string]Content{}

func BlockRobots(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Robots-Tag", "noindex")
		next.ServeHTTP(w, r)
	})
}

func init() {
	ctx := context.Background()
	keys, err := db.Redis.SMembers(ctx, "seo").Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the seo keys", slog.String("error", err.Error()))
		log.Panicln(err)
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
		log.Panicln((err))
	}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the seo", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		val := cmd.(*redis.MapStringStringCmd).Val()

		createdAt, err := strconv.ParseInt(val["created_at"], 10, 64)
		if err != nil {
			slog.Error("cannot parse the created at", slog.String("created_at", val["created_at"]))
			continue
		}

		updatedAt, err := strconv.ParseInt(val["updated_at"], 10, 64)
		if err != nil {
			slog.Error("cannot parse the product updated at", slog.String("updated_at", val["updated_at"]))
			continue
		}

		k := val["key"]
		c := URLs[k]
		c.Key = k
		c.Title = val["title"]
		c.URL = val["url"]
		c.Description = val["description"]
		c.CreatedAt = time.Unix(createdAt, 0)
		c.UpdatedAt = time.Unix(updatedAt, 0)
		URLs[k] = c

	}
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

func (c Content) Save(ctx context.Context) error {
	l := slog.With(slog.String("key", c.Key))
	l.LogAttrs(ctx, slog.LevelInfo, "saving a seo")

	key := "seo:" + c.Key
	now := time.Now()
	prv := URLs[c.Key]
	routes := router.R.Routes()

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, key,
			"title", db.Escape(c.Title),
			"description", db.Escape(c.Description),
			"url", db.Escape(c.URL),
			"updated_at", now.Unix(),
		)

		for _, route := range routes {
			if route.Pattern == prv.URL {

				handler, ok := route.Handlers["GET"].(http.HandlerFunc)
				if !ok {
					slog.LogAttrs(ctx, slog.LevelError, "cannot make type asserting for the handler", slog.String("url", prv.URL))
					return errors.New("something went wrong")
				}

				router.R.Get(prv.URL, http.NotFound)
				router.R.Get(c.URL, handler)

				l.LogAttrs(ctx, slog.LevelInfo, "previous route is replace", slog.String("previous", prv.URL), slog.String("url", c.URL))
			}
		}

		rdb.SAdd(ctx, "seo", c.Key)

		return nil
	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot store the seo", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	c.UpdatedAt = now
	URLs[c.Key] = c

	l.LogAttrs(ctx, slog.LevelInfo, "seo saved")

	return nil
}

func Find(ctx context.Context, key string) (Content, error) {
	l := slog.With(slog.String("key", key))
	l.LogAttrs(ctx, slog.LevelInfo, "looking for seo")

	if key == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate empty seo id")
		return Content{}, errors.New("input:id")
	}

	c := URLs[key]

	if c.Title == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot find the seo")
		return Content{}, errors.New("the data is not found")
	}

	return c, nil
}

func List(ctx context.Context, offset, num int) SearchResults {
	o := math.Min(float64(offset), float64(len(URLs)))
	n := math.Min(float64(num), float64(len(URLs)))

	return SearchResults{
		Total:   len(URLs),
		Content: maps.Values(URLs)[int(o):int(n)],
	}
}
