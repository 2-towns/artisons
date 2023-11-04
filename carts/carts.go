// Package carts manage the user cart
package carts

import (
	"context"
	"errors"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/orders"
	"gifthub/products"
	"gifthub/string/stringutil"
	"log"
	"log/slog"
	"strings"

	"github.com/redis/go-redis/v9"
)

type Cart struct {
	// The cart id
	ID string

	// "collect" or "home"
	Delivery string

	// "cash", "card ", "bitcoin" or "wire"
	Payment string

	Products []products.Product
}

func cartExists(c context.Context, cid string) bool {
	l := slog.With(slog.String("cid", cid))
	l.LogAttrs(c, slog.LevelInfo, "checking if the car exists")

	ctx := context.Background()
	ttl, err := db.Redis.TTL(ctx, "cart:"+cid).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "error when retrieving ttl for cart", slog.String("error", err.Error()))
		return false
	}

	slog.LogAttrs(c, slog.LevelInfo, "got ttl response", slog.Int64("ttl", ttl.Nanoseconds()))

	return ttl.Nanoseconds() > 0
}

// Add a product into a cart with its quantity
// Verify that the cart and the product exists.
func Add(c context.Context, cid, pid string, quantity int) error {
	l := slog.With(slog.String("cid", cid), slog.String("pid", pid), slog.Int("quantity", quantity))
	l.LogAttrs(c, slog.LevelInfo, "adding a product to the cart")

	if !cartExists(c, cid) {
		return errors.New("cart_not_found")
	}

	if !products.Available(c, pid) {
		return errors.New("product_not_found")
	}

	ctx := context.Background()
	if _, err := db.Redis.HIncrBy(ctx, "cart:"+cid, "product:"+pid, int64(quantity)).Result(); err != nil {
		l.LogAttrs(c, slog.LevelError, " cannot store the product", slog.String("error", err.Error()))
		return errors.New("something_went_wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "product added in the cart")

	return nil
}

// Get the full session cart.
// TODO: Get the product detail from the products package
func Get(c context.Context, cid string) (Cart, error) {
	l := slog.With(slog.String("cid", cid))
	l.LogAttrs(c, slog.LevelInfo, "get the cart")

	if !cartExists(c, cid) {
		return Cart{}, errors.New("cart_not_found")
	}

	ctx := context.Background()
	values, err := db.Redis.HGetAll(ctx, "cart:"+cid).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot get the cart", slog.String("error", err.Error()))
		return Cart{}, errors.New("something_went_wrong")
	}

	products := []products.Product{}
	//pipe := db.Redis.Pipeline()
	for key := range values {
		if !strings.HasPrefix(key, "product:") {
			continue
		}
		k := strings.Replace(key, "product:", "", 1)

		// TODO: Get the product detail from the products package
		log.Println(k)
	}

	l.LogAttrs(c, slog.LevelInfo, "got the cart with products", slog.Int("products", len(products)))

	return Cart{
		ID:       cid,
		Delivery: values["delivery"],
		Payment:  values["payment"],
		Products: products,
	}, nil
}

// UpdateDelivery update the delivery mode in Redis.
func (c Cart) UpdateDelivery(co context.Context, d string) error {
	l := slog.With(slog.String("cid", c.ID), slog.String("delivery", d))
	l.LogAttrs(co, slog.LevelInfo, "updating the delivery")

	if !orders.IsValidDelivery(co, d) {
		return errors.New("unauthorized")
	}

	ctx := context.Background()
	if _, err := db.Redis.HSet(ctx, "cart:"+c.ID, "delivery", d).Result(); err != nil {
		l.LogAttrs(co, slog.LevelError, "cannot update the delivery", slog.String("err", err.Error()))
		return errors.New("something_went_wrong")
	}

	l.LogAttrs(co, slog.LevelInfo, "the delivery is updated")

	return nil
}

// UpdatePayment update the payment mode in Redis.
func (c Cart) UpdatePayment(co context.Context, p string) error {
	l := slog.With(slog.String("cid", c.ID), slog.String("payment", p))
	l.LogAttrs(co, slog.LevelInfo, "updating the payment")

	if !orders.IsValidPayment(co, p) {
		return errors.New("unauthorized")
	}

	ctx := context.Background()
	if _, err := db.Redis.HSet(ctx, "cart:"+c.ID, "payment", p).Result(); err != nil {
		l.LogAttrs(co, slog.LevelError, "cannot update the payment", slog.String("err", err.Error()))
		return errors.New("something_went_wrong")
	}

	l.LogAttrs(co, slog.LevelInfo, "the payment is updated")

	return nil
}

// RefreshCID refreshes a cart ID (CID).
// If the CID does not exist, it will be created,
// with an expiration time.
func RefreshCID(c context.Context, s string) (string, error) {
	cid := s

	if cid == "" {
		id, err := stringutil.Random()
		if err != nil {
			slog.LogAttrs(c, slog.LevelError, "cannot generate a new id", slog.String("error", err.Error()))
			return "", errors.New("something_went_wrong")
		}
		cid = id
	}

	l := slog.With(slog.String("cid", s))

	ctx := context.Background()
	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, "cart:"+cid, "cid", cid)
		rdb.Expire(ctx, "cart:"+cid, conf.CartDuration)

		return nil
	}); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot refresh the cart", slog.String("err", err.Error()))
		return "", errors.New("something_went_wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "the cart is refreshed")

	return cid, nil
}
