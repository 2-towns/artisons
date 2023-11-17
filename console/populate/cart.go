package populate

import (
	"context"
	"gifthub/carts"
	"gifthub/conf"
	"gifthub/db"
)

func Cart(ctx context.Context, cid string, uid int64) (carts.Cart, error) {
	if _, err := db.Redis.HSet(ctx, "cart:"+cid, "cid", "test").Result(); err != nil {
		return carts.Cart{}, err
	}

	_, err := db.Redis.Set(ctx, "cart:"+cid+":user", uid, conf.CartDuration).Result()

	return carts.Cart{ID: cid}, err
}
