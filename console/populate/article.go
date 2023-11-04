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

	_, err = db.Redis.HSet(ctx, fmt.Sprintf("article:%d", id),
		"id", id,
		"title", "Manger de l'ail c'est bon pour la santé",
		"slug", "manger-de-l-ail-c-est-bon-pour-la-santé",
		"description", "C'est un antiseptique.",
		"image", "/tmp/hello",
		"online", fmt.Sprintf("%t", online),
		"updated_at", now.Format(time.RFC3339),
		"created_at", now.Format(time.RFC3339),
	).Result()

	if err != nil {
		return blogs.Article{}, err
	}

	_, err = db.Redis.ZAdd(ctx, "articles", redis.Z{
		Score:  float64(now.Unix()),
		Member: id,
	}).Result()

	if err != nil {
		return blogs.Article{}, err
	}

	return blogs.Article{ID: id}, err
}
