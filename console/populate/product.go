package populate

import (
	"context"
	"fmt"
	"gifthub/db"
	"gifthub/products"
	"gifthub/string/stringutil"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/redis/go-redis/v9"
)

func product(ctx context.Context, pipe redis.Pipeliner, title, description, pid, sku string) {
	now := time.Now().Unix()

	pipe.HSet(ctx, "product:"+pid,
		"id", pid,
		"sku", sku,
		"title", db.Escape(title),
		"description", db.Escape(description),
		"slug", stringutil.Slugify(db.Escape(title)),
		"currency", "EUR",
		"price", 100.5,
		"quantity", rand.Intn(10),
		"status", "online",
		"weight", rand.Float32(),
		"mid", faker.Phonenumber(),
		"tags", "clothes",
		"image_1", fmt.Sprintf("%s%s", pid, ".jpeg"),
		"image_2", fmt.Sprintf("%s%s", pid, ".jpeg"),
		"links", "",
		"meta", products.SerializeMeta(ctx, map[string]string{"color": "blue"}, ";"),
		"created_at", now,
		"updated_at", now,
	)

	pipe.ZAdd(ctx, "products", redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: pid,
	})
}
