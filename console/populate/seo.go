package populate

import (
	"context"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/seo"
	"time"

	"github.com/redis/go-redis/v9"
)

func url(ctx context.Context, pipe redis.Pipeliner, c seo.Content) {
	now := time.Now().Unix()

	pipe.HSet(ctx, "seo:"+c.Key,
		"title", db.Escape(c.Title),
		"description", db.Escape(c.Description),
		"url", db.Escape(c.URL),
		"key", c.Key,
		"lang", conf.DefaultLocale.String(),
		"created_at", now,
		"updated_at", now,
	)

	pipe.SAdd(ctx, "seo", c.Key)
}
