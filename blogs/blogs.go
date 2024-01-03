// blogs contains the function related to the content creation for a blog
package blogs

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/string/stringutil"
	"gifthub/validators"
	"log/slog"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type Article struct {
	ID          int64
	Title       string `validate:"required"`
	Slug        string
	Description string `validate:"required"`

	// The image path
	Image string `validate:"required"`

	CreatedAt int64
	UpdatedAt int64
}

func (a Article) Save(c context.Context) error {
	slog.LogAttrs(c, slog.LevelInfo, "creating a blog article")

	if err := validators.V.Struct(a); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot validate the article", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input_%s_required", low)
	}

	ctx := context.Background()
	id, err := db.Redis.Incr(ctx, "article_next_id").Result()
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get the next id", slog.String("error", err.Error()))
		return errors.New("error_http_general")
	}

	slug := stringutil.Slugify(a.Title)

	now := time.Now()

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		image := path.Join(conf.ImgProxy.Path, "blog", fmt.Sprintf("%d", id))

		err := os.Rename(a.Image, image)
		if err != nil {
			slog.LogAttrs(c, slog.LevelError, "cannot move the temporary file", slog.String("file", a.Image), slog.String("error", err.Error()))
			return err
		}

		rdb.HSet(ctx, fmt.Sprintf("article:%d", id), "id", id,
			"title", a.Title,
			"description", a.Description,
			"image", a.Image,
			"slug", slug,
			"created_at", time.Now().Unix(),
			"updated_at", time.Now().Unix(),
		)
		rdb.ZAdd(ctx, "articles", redis.Z{
			Score:  float64(now.Unix()),
			Member: id,
		})

		return nil
	}); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return errors.New("error_http_general")
	}

	slog.LogAttrs(c, slog.LevelInfo, "blog article created", slog.Int64("id", id))

	return nil
}

func List(c context.Context, page int) ([]Article, error) {
	l := slog.With(slog.Int("page", page))
	l.LogAttrs(c, slog.LevelInfo, "listing blog articles")

	start := int64(page * conf.ItemsPerPage)
	end := start + conf.ItemsPerPage
	ctx := context.Background()

	ids, err := db.Redis.ZRange(ctx, "articles", start, end).Result()
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return []Article{}, errors.New("error_http_general")
	}

	cmds, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for _, value := range ids {
			rdb.HGetAll(ctx, "article:"+value)
		}

		return nil
	})

	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot get the articles", slog.String("error", err.Error()))
		return []Article{}, errors.New("error_http_general")
	}

	articles := []Article{}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil {
			slog.LogAttrs(c, slog.LevelError, "cannot get the article", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		data := cmd.(*redis.MapStringStringCmd).Val()

		id, err := strconv.ParseInt(data["id"], 10, 64)
		if err != nil {
			slog.LogAttrs(c, slog.LevelError, "cannot parse the id", slog.String("id", data["id"]), slog.String("error", err.Error()))
			continue
		}

		createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
		if err != nil {
			l.LogAttrs(c, slog.LevelError, "cannot parse the created at date", slog.String("error", err.Error()), slog.Int64("id", id), slog.String("created_at", data["created_at"]))
			continue
		}

		updatedAt, err := strconv.ParseInt(data["updated_at"], 10, 64)
		if err != nil {
			l.LogAttrs(c, slog.LevelError, "cannot parse the created at date", slog.String("error", err.Error()), slog.Int64("id", id), slog.String("updated_at", data["updated_at"]))
			continue
		}

		image := path.Join(conf.ImgProxy.Path, "blog", fmt.Sprintf("%d", id))

		a := Article{
			ID:          id,
			Title:       data["title"],
			Description: data["description"],
			Slug:        data["slug"],
			Image:       image,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		articles = append(articles, a)
	}

	slog.LogAttrs(c, slog.LevelInfo, "blog article list done", slog.Int("length", len(articles)))

	return articles, nil
}

func Delete(c context.Context, id int64) error {
	l := slog.With(slog.Int64("id", id))
	l.LogAttrs(c, slog.LevelInfo, "deleting blog article")

	ctx := context.Background()

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		image := path.Join(conf.ImgProxy.Path, "blog", fmt.Sprintf("%d", id))

		err := os.Remove(image)
		if err != nil {
			slog.LogAttrs(c, slog.LevelError, "cannot remove the temporary file", slog.String("file", image), slog.String("error", err.Error()))
			return err
		}

		rdb.Del(ctx, fmt.Sprintf("article:%d", id))
		rdb.ZRem(ctx, "articles", id)

		return nil
	}); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot delete the data", slog.String("error", err.Error()))
		return errors.New("error_http_general")
	}

	return nil
}
