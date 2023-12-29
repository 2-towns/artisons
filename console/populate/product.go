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
	filenames := []string{}
	n := rand.Intn(3) + 1

	for i := 0; i < n; i++ {
		filenames = append(filenames, fmt.Sprintf("%d%s", time.Now().UnixNano(), ".jpg"))
	}

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
		"image_1", fmt.Sprintf("%d%s", time.Now().UnixNano(), ".jpg"),
		"image_2", fmt.Sprintf("%d%s", time.Now().UnixNano(), ".jpg"),
		"links", "",
		"meta", products.SerializeMeta(ctx, map[string]string{"color": "blue"}, ";"),
		"created_at", time.Now().Unix(),
		"updated_at", time.Now().Unix(),
	)

	pipe.ZAdd(ctx, "products", redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: pid,
	})
}
