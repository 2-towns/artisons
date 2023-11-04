package populate

import (
	"context"
	"fmt"
	"gifthub/db"
	"gifthub/orders"
	"time"
)

func Order(ctx context.Context, oid string, uid int64, ids map[string]int64) (orders.Order, error) {
	now := time.Now()

	_, err := db.Redis.HSet(ctx, "order:"+oid,
		"id", oid,
		"uid", uid,
		"delivery", "home",
		"payment", "card",
		"payment_status", "payment_progress",
		"status", "created",
		"updated_at", now.Format(time.RFC3339),
		"created_at", now.Format(time.RFC3339),
	).Result()

	if err != nil {
		return orders.Order{}, err
	}

	for key, value := range ids {
		_, err := db.Redis.HSet(ctx, "order:"+oid, "product:"+key, value).Result()
		if err != nil {
			return orders.Order{}, err
		}
	}

	_, err = db.Redis.HSet(ctx, fmt.Sprintf("user:%d", uid), "order:"+oid, oid).Result()

	return orders.Order{}, err
}
