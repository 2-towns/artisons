package populate

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func order(ctx context.Context, pipe redis.Pipeliner, oid string, uid int64, ids map[string]int64) {
	now := time.Now()

	createdAt, _ := time.Parse(time.DateTime, "2023-11-10 15:04:05")
	updatedAt, _ := time.Parse(time.DateTime, "2023-11-10 15:04:05")

	pipe.HSet(ctx, "order:"+oid,
		"id", oid,
		"uid", uid,
		"delivery", "home",
		"payment", "card",
		"payment_status", "payment_progress",
		"status", "created",
		"total", "100.5",
		"type", "order",
		"address_lastname", "Arnaud",
		"address_firstname", "Arnaud",
		"address_city", "Lille",
		"address_street", "Rue du moulin",
		"address_complementary", "Appartement C",
		"address_zipcode", "59000",
		"address_phone", "3345668832",
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
