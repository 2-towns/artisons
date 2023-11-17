package populate

import (
	"context"
	"fmt"
	"gifthub/db"
	"gifthub/orders"
	"time"

	"github.com/redis/go-redis/v9"
)

func Order(ctx context.Context, oid string, uid int64, ids map[string]int64) (orders.Order, error) {
	now := time.Now()

	createdAt, _ := time.Parse(time.DateTime, "2023-11-10 15:04:05")
	updatedAt, _ := time.Parse(time.DateTime, "2023-11-10 15:04:05")

	if _, err := db.Redis.HSet(ctx, "order:"+oid,
		"id", oid,
		"uid", uid,
		"delivery", "home",
		"payment", "card",
		"payment_status", "payment_progress",
		"status", "created",
		"updated_at", updatedAt.Unix(),
		"created_at", createdAt.Unix(),
	).Result(); err != nil {
		return orders.Order{}, err
	}

	for key, value := range ids {
		if _, err := db.Redis.HSet(ctx, "order:"+oid+":products", key, value).Result(); err != nil {
			return orders.Order{}, err
		}
	}

	_, err := db.Redis.ZAdd(ctx, fmt.Sprintf("user:%s:orders", oid), redis.Z{
		Score:  float64(now.Unix()),
		Member: oid,
	}).Result()

	return orders.Order{}, err
}
