package populate

import (
	"context"
	"fmt"
	"gifthub/conf"
	"path"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/text/language"
)

func article(ctx context.Context, pipe redis.Pipeliner, id int64, online bool) {
	pipe.Incr(ctx, "blog_next_id").Result()
	now := time.Now()

	pipe.HSet(ctx, fmt.Sprintf("blog:%d", id),
		"id", id,
		"title", "Manger de l'ail c'est bon pour la santé",
		"slug", "manger-de-l-ail-c-est-bon-pour-la-santé",
		"status", "online",
		"lang", language.English.String(),
		"description", "C'est un antiseptique.",
		"image", path.Join(conf.WorkingSpace, "web", "images", "blog", fmt.Sprintf("%d.jpeg", id)),
		"online", fmt.Sprintf("%t", online),
		"updated_at", now.Unix(),
		"created_at", now.Unix(),
	)
}
