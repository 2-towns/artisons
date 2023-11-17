package populate

import (
	"context"
	"fmt"
	"gifthub/blogs"
	"gifthub/db"
	"time"

	"github.com/redis/go-redis/v9"
)

func Article(ctx context.Context, online bool) (blogs.Article, error) {
	id, err := db.Redis.Incr(ctx, "article_next_id").Result()
	if err != nil {
		return blogs.Article{}, err
	}

	now := time.Now()

	if _, err = db.Redis.HSet(ctx, fmt.Sprintf("article:%d", id),
		"id", id,
		"title", "Manger de l'ail c'est bon pour la santé",
		"slug", "manger-de-l-ail-c-est-bon-pour-la-santé",
		"description", "C'est un antiseptique.",
		"image", "/tmp/hello",
		"online", fmt.Sprintf("%t", online),
		"updated_at", now.Unix(),
		"created_at", now.Unix(),
	).Result(); err != nil {
		return blogs.Article{}, err
	}

	if _, err = db.Redis.ZAdd(ctx, "articles", redis.Z{
		Score:  float64(now.Unix()),
		Member: id,
	}).Result(); err != nil {
		return blogs.Article{}, err
	}

	return blogs.Article{ID: id}, err
}
