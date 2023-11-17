// Package populate provide script to populate date into Redis
package populate

import (
	"context"
	"gifthub/db"
)

// Run the populate script. It will flush the database first
func Run() error {
	ctx := context.Background()

	db.Redis.FlushDB(ctx)

	pid := "test"
	sku := "skutest"
	price := float32(100.5)
	product, err := Product(ctx, pid, sku, price)
	if err != nil {
		return err
	}

	alive := true
	user, err := User(ctx, "test", alive)
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
	user, err = User(ctx, "expired", expired)
	if err != nil {
		return err
	}

	if _, err = Shop(ctx); err != nil {
		return err
	}

	if err = Tag(ctx); err != nil {
		return err
	}

	err = Stats(ctx)
	if err != nil {
		return err
	}

	return err
}
