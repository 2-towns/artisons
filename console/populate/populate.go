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

	product, err := Product(ctx, "test")
	if err != nil {
		return err
	}

	alive := true
	user, err := User(ctx, "test", alive)
	if err != nil {
		return err
	}

	_, err = Order(ctx, "test", user.ID, map[string]int64{product.PID: 1})
	if err != nil {
		return err
	}

	_, err = Cart(ctx, "test")
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
