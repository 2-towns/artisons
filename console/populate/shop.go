package populate

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func shop(ctx context.Context, pipe redis.Pipeliner) {
	now := time.Now()

	pipe.HSet(ctx, "shop",
		"logo", "../web/images/logo",
		"active", "1",
		"guest", "1",
		"stock", "1",
		"name", "My Shop",
		"city", "Oran",
		"address", "Hay Yasmine",
		"zipcode", "31244",
		"phone", "0559682532",
		"updated_at", now.Unix(),
	)
}
