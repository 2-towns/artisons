package populate

import (
	"context"
	"gifthub/products"
	"gifthub/string/stringutil"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/redis/go-redis/v9"
)

func product(ctx context.Context, pipe redis.Pipeliner, title, description, pid, sku string) {
	pipe.HSet(ctx, "product:"+pid,
		"id", pid,
		"sku", sku,
		"title", title,
		"description", description,
		"slug", stringutil.Slugify(title),
		"length", rand.Intn(4)+1,
		"currency", "EUR",
		"price", 100.5,
		"quantity", rand.Intn(10),
		"status", "online",
		"weight", rand.Float32(),
		"mid", faker.Phonenumber(),
		"tags", "clothes",
		"links", "",
		"meta", products.SerializeMeta(ctx, map[string]string{"color": "blue"}, ";"),
		"created_at", time.Now().Unix(),
		"updated_at", time.Now().Unix(),
	)
}
