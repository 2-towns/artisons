package populate

import (
	"context"
	"gifthub/db"
	"gifthub/products"
	"gifthub/string/stringutil"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
)

var titles []string = []string{
	"Une belle paire de claquette",
	"Un joli pull tricoté par ma maman",
	"Une paire de chaussette en cuir",
}

func Product(ctx context.Context, pid string) (products.Product, error) {
	sku, err := stringutil.Random()
	if err != nil {
		return products.Product{}, err
	}

	idx := rand.Intn(len(titles) - 1)
	title := titles[idx]
	description := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. "
	tags := []string{"gift"}
	links := []string{"http://google.fr"}
	meta := map[string]string{"color": "blue"}
	key := "product:" + pid

	_, err = db.Redis.HSet(ctx, key,
		"id", pid,
		"sku", sku,
		"title", title,
		"description", description,
		"slug", stringutil.Slugify(title),
		"length", rand.Intn(4),
		"currency", "EUR",
		"price", rand.Float32(),
		"quantity", rand.Intn(10),
		"status", "online",
		"weight", rand.Float32(),
		"mid", faker.Phonenumber(),
		"updated_at", time.Now().Format(time.RFC3339),
	).Result()
	if err != nil {
		return products.Product{}, err
	}

	_, err = db.Redis.SAdd(ctx, key+":links", links).Result()
	if err != nil {
		return products.Product{}, err
	}

	_, err = db.Redis.SAdd(ctx, key+":tags", tags).Result()
	if err != nil {
		return products.Product{}, err
	}

	_, err = db.Redis.HSet(ctx, key+":options", meta).Result()

	return products.Product{
		ID: pid,
	}, err
}
