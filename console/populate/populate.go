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

	_, err = Order(ctx, "test", user.ID, map[string]int64{product.ID: 1})
	if err != nil {
		return err
	}

	_, err = Cart(ctx, "test", user.ID)
	if err != nil {
		return err
	}

	online := true
	_, err = Article(ctx, online)
	if err != nil {
		return err
	}

	offline := false
	_, err = Article(ctx, offline)
	if err != nil {
		return err
	}

	_, err = Article(ctx, offline)
	if err != nil {
		return err
	}

	expired := false
	user, err = User(ctx, "expired", expired)
	if err != nil {
		return err
	}

	return err
}
