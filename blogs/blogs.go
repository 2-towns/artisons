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
)

type Article struct {
	ID          int64
	Title       string `validate:"required"`
	Slug        string
	Description string `validate:"required"`
	Status      string `redis:"status" validate:"oneof=online offline"`

	// The image path
	Image string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type SearchResults struct {
	Total    int64
	Articles []Article
}

type Query struct {
	Keywords string
	Lang     string
}

func (p Article) Validate(c context.Context) error {
	slog.LogAttrs(c, slog.LevelInfo, "validating a article")

	if err := validators.V.Struct(p); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot validate the article", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input:%s", low)
	}

	slog.LogAttrs(c, slog.LevelInfo, "article validated")

	return nil
}

func NextID(ctx context.Context) (int64, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "getting the next article id")

	id, err := db.Redis.Incr(ctx, "blog_next_id").Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the next id", slog.String("error", err.Error()))
		return 0, errors.New("something went wrong")
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "next id generated", slog.Int64("id", id))

	return id, nil
}

func (a Article) Save(c context.Context) error {
	slog.LogAttrs(c, slog.LevelInfo, "creating a blog article")

	ctx := context.Background()
	slug := stringutil.Slugify(a.Title)
	now := time.Now().Unix()

	if _, err := db.Redis.HSet(ctx, fmt.Sprintf("blog:%d", a.ID), "id", a.ID,
		"title", a.Title,
		"description", a.Description,
		"image", a.Image,
		"slug", slug,
		"status", a.Status,
		"created_at", now,
		"updated_at", now,
	).Result(); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	slog.LogAttrs(c, slog.LevelInfo, "blog article created", slog.Int64("id", a.ID))

	return nil
}

func parse(ctx context.Context, data map[string]string) (Article, error) {
	id, err := strconv.ParseInt(data["id"], 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.String("id", data["id"]), slog.String("error", err.Error()))
		return Article{}, errors.New("input:id")
	}

	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the created at", slog.String("error", err.Error()), slog.Int64("id", id), slog.String("created_at", data["created_at"]))
		return Article{}, errors.New("input:created_at")
	}

	updatedAt, err := strconv.ParseInt(data["updated_at"], 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the updated at", slog.String("error", err.Error()), slog.Int64("id", id), slog.String("updated_at", data["updated_at"]))
		return Article{}, errors.New("input:updated_at")
	}

	image := path.Join(conf.ImgProxy.Path, "blog", fmt.Sprintf("%d", id))

	a := Article{
		ID:          id,
		Title:       db.Unescape(data["title"]),
		Description: db.Unescape(data["description"]),
		Slug:        data["slug"],
		Status:      data["status"],
		Image:       image,
		CreatedAt:   time.Unix(createdAt, 0),
		UpdatedAt:   time.Unix(updatedAt, 0),
	}

	return a, nil
}

func Search(c context.Context, q Query, offset, num int) (SearchResults, error) {
	slog.LogAttrs(c, slog.LevelInfo, "searching products")

	qs := "@status:{online} "

	if q.Keywords != "" {
		k := db.Escape(q.Keywords)
		qs += fmt.Sprintf("(@title:'*%s*')|(@description:'*%s*')|(@id:'{%s})'", k, k, k)
	}

	if q.Lang != "" {
		qs += fmt.Sprintf("(@lang:'{%s})'", q.Lang)
	}

	slog.LogAttrs(c, slog.LevelInfo, "preparing redis request", slog.String("query", qs))

	ctx := context.Background()
	cmds, err := db.Redis.Do(
		ctx,
		"FT.SEARCH",
		db.BlogIdx,
		qs,
		"LIMIT",
		fmt.Sprintf("%d", offset),
		fmt.Sprintf("%d", num),
		"SORTBY",
		"updated_at",
		"desc",
	).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot run the search query", slog.String("query", qs), slog.String("error", err.Error()))
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

		product, err := parse(c, data)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the blog", slog.Any("blog", data), slog.String("error", err.Error()))
			continue
		}

		articles = append(articles, product)
	}

	slog.LogAttrs(c, slog.LevelInfo, "search done", slog.Int64("results", total))

	return SearchResults{
		Total:    total,
		Articles: articles,
	}, nil
}

func Delete(c context.Context, id int64) error {
	l := slog.With(slog.Int64("id", id))
	l.LogAttrs(c, slog.LevelInfo, "deleting blog article")

	ctx := context.Background()

	if _, err := db.Redis.Del(ctx, fmt.Sprintf("blog:%d", id)).Result(); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot delete the data", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	image := path.Join(conf.ImgProxy.Path, "blog", fmt.Sprintf("%d", id))
	err := os.Remove(image)
	if err != nil {
		slog.LogAttrs(c, slog.LevelWarn, "cannot remove the image", slog.String("file", image), slog.String("error", err.Error()))
		return nil
	}

	l.LogAttrs(c, slog.LevelInfo, "the article is deleted successfuly")

	return nil
}

// Find looks for a blog by its id
func Find(c context.Context, id int64) (Article, error) {
	l := slog.With(slog.Int64("id", id))
	l.LogAttrs(c, slog.LevelInfo, "looking for article")

	if id == 0 {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate empty article id")
		return Article{}, errors.New("input:id")
	}

	ctx := context.Background()

	if exists, err := db.Redis.Exists(ctx, fmt.Sprintf("blog:%d", id)).Result(); exists == 0 || err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot find the blog")
		return Article{}, errors.New("oops the data is not found")
	}

	data, err := db.Redis.HGetAll(ctx, fmt.Sprintf("blog:%d", id)).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot find the article", slog.String("error", err.Error()))
		return Article{}, err
	}

	a, err := parse(c, data)

	if err != nil {
		l.LogAttrs(c, slog.LevelInfo, "the article is found")
	}

	return a, err
}
