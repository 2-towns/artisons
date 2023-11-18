package populate

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func order(ctx context.Context, pipe redis.Pipeliner, oid string, uid int64, ids map[string]int64) {
	now := time.Now()

	createdAt, _ := time.Parse(time.RFC3339, "2023-11-10T15:04:05Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2023-11-10T15:04:05Z")

	pipe.HSet(ctx, "order:"+oid,
		"id", oid,
		"uid", uid,
		"delivery", "home",
		"payment", "card",
		"payment_status", "payment_progress",
		"status", "created",
		"total", "100.5",
		"updated_at", updatedAt.Unix(),
		"created_at", createdAt.Unix(),
	)

	for key, value := range ids {
		pipe.HSet(ctx, "order:"+oid+":products", key, value).Result()
	}

	pipe.ZAdd(ctx, fmt.Sprintf("user:%s:orders", oid), redis.Z{
		Score:  float64(now.Unix()),
		Member: oid,
	}).Result()
}
