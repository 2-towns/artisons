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
}

func (p Article) Validate(c context.Context) error {
	slog.LogAttrs(c, slog.LevelInfo, "validating a product")

	if err := validators.V.Struct(p); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot validate the article", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input_%s_invalid", low)
	}

	return nil
}

func NextID(ctx context.Context) (int64, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "getting the next article id")

	id, err := db.Redis.Incr(ctx, "blog_next_id").Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the next id", slog.String("error", err.Error()))
		return 0, errors.New("error_http_general")
	}

	return id, nil
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
	slug := stringutil.Slugify(a.Title)
	now := time.Now()

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, fmt.Sprintf("blog:%d", a.ID), "id", a.ID,
			"title", a.Title,
			"description", a.Description,
			"image", a.Image,
			"slug", slug,
			"status", a.Status,
			"created_at", time.Now().Unix(),
			"updated_at", time.Now().Unix(),
		)
		rdb.ZAdd(ctx, "blog", redis.Z{
			Score:  float64(now.Unix()),
			Member: a.ID,
		})

		return nil
	}); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return errors.New("error_http_general")
	}

	slog.LogAttrs(c, slog.LevelInfo, "blog article created", slog.Int64("id", a.ID))

	return nil
}

func parse(ctx context.Context, data map[string]string) (Article, error) {
	id, err := strconv.ParseInt(data["id"], 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.String("id", data["id"]), slog.String("error", err.Error()))
		return Article{}, errors.New("input_id_invalid")
	}

	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the created at date", slog.String("error", err.Error()), slog.Int64("id", id), slog.String("created_at", data["created_at"]))
		return Article{}, errors.New("input_created_at_invalid")
	}

	updatedAt, err := strconv.ParseInt(data["updated_at"], 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the created at date", slog.String("error", err.Error()), slog.Int64("id", id), slog.String("updated_at", data["updated_at"]))
		return Article{}, errors.New("input_updated_at_invalid")
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

	slog.LogAttrs(c, slog.LevelInfo, "search done", slog.Int64("results", total))

	results := res["results"].([]interface{})
	articles := []Article{}

	for _, value := range results {
		m := value.(map[interface{}]interface{})
		attributes := m["extra_attributes"].(map[interface{}]interface{})
		data := db.ConvertMap(attributes)

		product, err := parse(c, data)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the product", slog.Any("product", data), slog.String("error", err.Error()))
			continue
		}

		articles = append(articles, product)
	}

	return SearchResults{
		Total:    total,
		Articles: articles,
	}, nil
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

		rdb.Del(ctx, fmt.Sprintf("blog:%d", id))
		rdb.ZRem(ctx, "blog", id)

		return nil
	}); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot delete the data", slog.String("error", err.Error()))
		return errors.New("error_http_general")
	}

	return nil
}

// Find looks for a blog by its id
func Find(c context.Context, id int64) (Article, error) {
	l := slog.With(slog.Int64("id", id))
	l.LogAttrs(c, slog.LevelInfo, "looking for article")

	if id == 0 {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate empty article id")
		return Article{}, errors.New("input_id_required")
	}

	ctx := context.Background()

	if exists, err := db.Redis.Exists(ctx, fmt.Sprintf("blog:%d", id)).Result(); exists == 0 || err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot find the blog")
		return Article{}, errors.New("error_http_blognotfound")
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

// func List(c context.Context, page int) ([]Article, error) {
// 	l := slog.With(slog.Int("page", page))
// 	l.LogAttrs(c, slog.LevelInfo, "listing blog articles")

// 	start := int64(page * conf.ItemsPerPage)
// 	end := start + conf.ItemsPerPage
// 	ctx := context.Background()

// 	ids, err := db.Redis.ZRange(ctx, "articles", start, end).Result()
// 	if err != nil {
// 		slog.LogAttrs(c, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
// 		return []Article{}, errors.New("error_http_general")
// 	}

// 	cmds, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
// 		for _, value := range ids {
// 			rdb.HGetAll(ctx, "blog:"+value)
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		l.LogAttrs(c, slog.LevelError, "cannot get the articles", slog.String("error", err.Error()))
// 		return []Article{}, errors.New("error_http_general")
// 	}

// 	articles := []Article{}

// 	for _, cmd := range cmds {
// 		key := fmt.Sprintf("%s", cmd.Args()[1])

// 		if cmd.Err() != nil {
// 			slog.LogAttrs(c, slog.LevelError, "cannot get the article", slog.String("key", key), slog.String("error", err.Error()))
// 			continue
// 		}

// 		data := cmd.(*redis.MapStringStringCmd).Val()

// 		id, err := strconv.ParseInt(data["id"], 10, 64)
// 		if err != nil {
// 			slog.LogAttrs(c, slog.LevelError, "cannot parse the id", slog.String("id", data["id"]), slog.String("error", err.Error()))
// 			continue
// 		}

// 		createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
// 		if err != nil {
// 			l.LogAttrs(c, slog.LevelError, "cannot parse the created at date", slog.String("error", err.Error()), slog.Int64("id", id), slog.String("created_at", data["created_at"]))
// 			continue
// 		}

// 		updatedAt, err := strconv.ParseInt(data["updated_at"], 10, 64)
// 		if err != nil {
// 			l.LogAttrs(c, slog.LevelError, "cannot parse the created at date", slog.String("error", err.Error()), slog.Int64("id", id), slog.String("updated_at", data["updated_at"]))
// 			continue
// 		}

// 		image := path.Join(conf.ImgProxy.Path, "blog", fmt.Sprintf("%d", id))

// 		a := Article{
// 			ID:          id,
// 			Title:       data["title"],
// 			Description: data["description"],
// 			Slug:        data["slug"],
// 			Image:       image,
// 			CreatedAt:   createdAt,
// 			UpdatedAt:   updatedAt,
// 		}

// 		articles = append(articles, a)
// 	}

// 	slog.LogAttrs(c, slog.LevelInfo, "blog article list done", slog.Int("length", len(articles)))

// 	return articles, nil
// }
