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
		"slug", "manger-de-l-ail-c-est-bon-pour-la-santé",
		"address_firstname", "Arnaud",
		"address_lastname", "None",
		"address_city", "Oran",
		"address_street", "Hay Yasmine",
		"address_complementary", "Hay Salam",
		"address_zipcode", "31244",
		"address_phone", "0559682532",
		"updated_at", now.Unix(),
	)
}
