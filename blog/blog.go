// blogs contains the function related to the content creation for a blog
package blog

import (
	"context"
	"errors"
	"fmt"
	"artisons/conf"
	"artisons/db"
	"artisons/validators"
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
	ID          int
	Title       string `validate:"required"`
	Type        string
	Slug        string `validate:"required"`
	Description string `validate:"required"`
	Status      string `redis:"status" validate:"oneof=online offline"`

	// The image path
	Image string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type SearchResults struct {
	Total    int
	Articles []Article
}

type Query struct {
	Keywords string
	Type     string
}

func Deletable(ctx context.Context, id int) (bool, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "checking deletable", slog.Int("key", id))

	typ, err := db.Redis.HGet(ctx, fmt.Sprintf("blog:%d", id), "type").Result()

	if err != nil && err.Error() != "redis: nil" {
		slog.LogAttrs(ctx, slog.LevelError, "cannot check blog deletable")
		return false, errors.New("something went wrong")
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "filter is deletable", slog.Bool("deletable", typ == "blog"))

	return typ == "blog", nil
}

func (p Article) Validate(ctx context.Context) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "validating a article")

	if err := validators.V.Struct(p); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot validate the article", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input:%s", low)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "article validated")

	return nil
}

func GetIDFromSlug(ctx context.Context, slug string) (int, error) {
	exists, err := db.Redis.HExists(ctx, "bids", slug).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot verify the slug existence", slog.String("error", err.Error()))
		return 0, errors.New("something went wrong")
	}

	if !exists {
		return 0, nil
	}

	s, err := db.Redis.HGet(ctx, "bids", slug).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the slug", slog.String("error", err.Error()))
		return 0, errors.New("something went wrong")
	}

	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id ", slog.String("id", s), slog.String("error", err.Error()))
		return 0, errors.New("something went wrong")
	}

	return int(id), nil
}

func (a *Article) UpdateImage(key, value string) {
	a.Image = value
}

func (a Article) Save(ctx context.Context) (string, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "creating a blog article")

	if a.ID == 0 {
		id, err := db.Redis.Incr(ctx, "blog_next_id").Result()
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the next id", slog.String("error", err.Error()))
			return "", errors.New("something went wrong")
		}
		a.ID = int(id)
	}

	now := time.Now().Unix()

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		key := fmt.Sprintf("blog:%d", a.ID)

		rdb.HSet(ctx, key,
			"title", db.Escape(a.Title),
			"description", db.Escape(a.Description),
			"image", a.Image,
			"slug", db.Escape(a.Slug),
			"status", a.Status,
			"updated_at", now,
		)
		rdb.HSetNX(ctx, key, "id", a.ID)
		rdb.HSetNX(ctx, key, "created_at", now)
		rdb.HSet(ctx, "bids", a.Slug, a.ID)

		return nil
	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot store the product", slog.String("error", err.Error()))
		return "", err
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "blog article created", slog.Int("id", a.ID))

	return fmt.Sprintf("%d", a.ID), nil
}

func parse(ctx context.Context, data map[string]string) (Article, error) {
	id, err := strconv.ParseInt(data["id"], 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.String("id", data["id"]), slog.String("error", err.Error()))
		return Article{}, errors.New("input:id")
	}

	updatedAt, err := strconv.ParseInt(data["updated_at"], 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the updated at", slog.String("error", err.Error()), slog.Int64("id", id), slog.String("updated_at", data["updated_at"]))
		return Article{}, errors.New("input:updated_at")
	}

	a := Article{
		ID:          int(id),
		Title:       db.Unescape(data["title"]),
		Description: db.Unescape(data["description"]),
		Slug:        data["slug"],
		Status:      data["status"],
		Image:       data["image"],
		Type:        data["type"],
		UpdatedAt:   time.Unix(updatedAt, 0),
	}

	return a, nil
}

func Search(ctx context.Context, q Query, offset, num int) (SearchResults, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "searching articles", slog.Int("offset", offset), slog.Int("num", num))

	qs := fmt.Sprintf("FT.SEARCH %s @status:{online}", db.BlogIdx)

	if q.Keywords != "" {
		k := db.SearchValue(q.Keywords)
		qs += fmt.Sprintf("(@title:%s)|(@description:%s)|(@id:{%s})", k, k, k)
	}

	if q.Type != "" {
		qs += fmt.Sprintf("(@type:{%s})", q.Type)
	}

	qs += fmt.Sprintf(" SORTBY updated_at desc LIMIT %d %d DIALECT 2", offset, num)

	slog.LogAttrs(ctx, slog.LevelInfo, "preparing redis request", slog.String("query", qs))

	args, err := db.SplitQuery(ctx, qs)
	if err != nil {
		return SearchResults{}, err
	}

	cmds, err := db.Redis.Do(ctx, args...).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot run the search query", slog.String("error", err.Error()))
		return SearchResults{}, err
	}

	res := cmds.(map[interface{}]interface{})
	total := res["total_results"].(int64)
	results := res["results"].([]interface{})
	articles := []Article{}

	for _, value := range results {
		m := value.(map[interface{}]interface{})
		attributes := m["extra_attributes"].(map[interface{}]interface{})
		data := db.ConvertMap(attributes)

		product, err := parse(ctx, data)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the blog", slog.Any("blog", data), slog.String("error", err.Error()))
			continue
		}

		articles = append(articles, product)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "search done", slog.Int64("results", total))

	return SearchResults{
		Total:    int(total),
		Articles: articles,
	}, nil
}

func Delete(ctx context.Context, id int) error {
	l := slog.With(slog.Int("id", id))
	l.LogAttrs(ctx, slog.LevelInfo, "deleting blog article")

	if _, err := db.Redis.Del(ctx, fmt.Sprintf("blog:%d", id)).Result(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot delete the data", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	image := path.Join(conf.ImgProxy.Path, "blog", fmt.Sprintf("%d", id))
	err := os.Remove(image)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelWarn, "cannot remove the image", slog.String("file", image), slog.String("error", err.Error()))
		return nil
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the article is deleted successfuly")

	return nil
}

// Find looks for a blog by its id
func Find(ctx context.Context, id int) (Article, error) {
	l := slog.With(slog.Int("id", id))
	l.LogAttrs(ctx, slog.LevelInfo, "looking for article")

	if id == 0 {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate empty article id")
		return Article{}, errors.New("input:id")
	}

	if exists, err := db.Redis.Exists(ctx, fmt.Sprintf("blog:%d", id)).Result(); exists == 0 || err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot find the blog")
		return Article{}, errors.New("oops the data is not found")
	}

	data, err := db.Redis.HGetAll(ctx, fmt.Sprintf("blog:%d", id)).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot find the article", slog.String("error", err.Error()))
		return Article{}, err
	}

	a, err := parse(ctx, data)

	if err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "the article is found")
	}

	return a, err
}

func FindBySlug(ctx context.Context, slug string) (Article, error) {
	l := slog.With(slog.String("slug", slug))
	l.LogAttrs(ctx, slog.LevelInfo, "looking for article from slug")

	if slug == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot continue with empty slug")
		return Article{}, errors.New("the data is not found")
	}

	id, err := GetIDFromSlug(ctx, slug)
	if err != nil {
		return Article{}, nil
	}

	if id == 0 {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot continue with empty id")
		return Article{}, errors.New("the data is not found")
	}

	return Find(ctx, id)
}
