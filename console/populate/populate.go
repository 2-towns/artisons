// Package populate provide script to populate date into Redis
package populate

import (
	"context"
	"gifthub/db"

	"github.com/redis/go-redis/v9"

	"gifthub/products"
	"gifthub/users"

	"github.com/go-faker/faker/v4"
)

func del(ctx context.Context, pipe redis.Pipeliner, pattern string) {
	keys, _ := db.Redis.Keys(ctx, pattern).Result()
	for _, key := range keys {
		pipe.Del(ctx, key).Err()
	}
}

// Run the populate script. It will flush the database first
func Run() error {
	ctx := context.Background()
	pipe := db.Redis.Pipeline()

	del(ctx, pipe, "user:*")
	del(ctx, pipe, "users:*")
	del(ctx, pipe, "stats:*")
	del(ctx, pipe, "product:*")
	del(ctx, pipe, "tag:*")
	del(ctx, pipe, "article:*")
	del(ctx, pipe, "order:*")
	del(ctx, pipe, "orders:*")
	del(ctx, pipe, "shop")
	del(ctx, pipe, "cart:*")
	del(ctx, pipe, "article_next_id")
	del(ctx, pipe, "user_next_id")

	_, err := pipe.Exec(ctx)

	if err != nil {
		return err
	}

	pipe = db.Redis.Pipeline()

	product(
		ctx,
		pipe,
		"T-shirt Tester c’est douter",
		"T-shirt développeur unisexe Tester c’est douter",
		"PDT1",
		"SKU1",
	)

	product(
		ctx,
		pipe,
		"T-shirt développeur unisexe JavaScript Park",
		" 100 % coton pour les couleurs unies",
		"PDT2",
		"SKU2",
	)

	product(
		ctx,
		pipe,
		"Bouteille en acier inoxydable",
		"En plus d'être canon, cette bouteille de 500 ml maintiendra votre boisson au chaud ou au froid pendant 6 heures.",
		"PDT3",
		"SKU3",
	)

	product(
		ctx,
		pipe,
		"Mug développeur",
		"Cet incroyable mug augmente les chances de réussite de vos mises en prod de 42%*",
		"PDT4",
		"SKU4",
	)

	product(
		ctx,
		pipe,
		"Sweat à capuche unisexe développeur",
		"Spécialement conçu pour vous réconforter pendant vos longues sessions de debug, ce sweat soutient également Les Joies du Code.",
		"PDT5",
		"PDT5",
	)

	alive := true
	user(ctx, pipe, "SES1", alive)
	alive := true
	user, err := User(ctx, "test", alive)
	err = products.Add(products.FakeProduct())
	if err != nil {
		return err
	}

	if _, err = Order(ctx, "test", user.ID, map[string]int64{product.ID: 1}); err != nil {
		return err
	}

	if _, err = Cart(ctx, "test", user.ID); err != nil {
		return err
	}

	online := true
	if _, err = Article(ctx, online); err != nil {
		return err
	}

	offline := false
	if _, err = Article(ctx, offline); err != nil {
		return err
	}

	if _, err = Article(ctx, offline); err != nil {
		return err
	}

	expired := false
	user(ctx, pipe, "expired", expired)

	pipe.SAdd(ctx, "admins", "hello@world.com", "lock@world.com")

	var uid int64 = 1
	order(ctx, pipe, "ORD1", uid, map[string]int64{"PDT1": 1})
	cart(ctx, pipe, "CAR1", uid)

	online := true
	article(ctx, pipe, 1, online)

	offline := false
	article(ctx, pipe, 2, offline)
	article(ctx, pipe, 3, offline)

	shop(ctx, pipe)

	tag(ctx, pipe)

	stats(ctx, pipe)

	_, err = pipe.Exec(ctx)

	return err
}
