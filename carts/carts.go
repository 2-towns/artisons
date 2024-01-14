// Package carts manage the user cart
package carts

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/contexts"
	"gifthub/orders"
	"gifthub/products"
	"gifthub/tracking"
	"log/slog"
	"strconv"
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

func cartExists(c context.Context) bool {
	cid := c.Value(contexts.Cart).(string)
	l := slog.With(slog.String("cid", cid))
	l.LogAttrs(c, slog.LevelInfo, "checking if the car exists")

	ctx := context.Background()
	ttl, err := db.Redis.TTL(ctx, "cart:"+cid+":user").Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "error when retrieving ttl for cart", slog.String("error", err.Error()))
		return false
	}

	slog.LogAttrs(c, slog.LevelInfo, "got ttl response", slog.Int64("ttl", ttl.Nanoseconds()))

	return ttl.Nanoseconds() > 0
}

// Add a product into a cart with its quantity
// Verify that the cart and the product exists.
func Add(c context.Context, pid string, quantity int) error {
	cid := c.Value(contexts.Cart).(string)
	l := slog.With(slog.String("cid", cid), slog.String("product_id", pid), slog.Int("quantity", quantity))
	l.LogAttrs(c, slog.LevelInfo, "adding a product to the cart")

	if !cartExists(c) {
		return errors.New("the session is expired")
	}

	if !products.Available(c, pid) {
		return errors.New("product_not_found")
	}

	ctx := context.Background()
	if _, err := db.Redis.HIncrBy(ctx, "cart:"+cid, pid, int64(quantity)).Result(); err != nil {
		l.LogAttrs(c, slog.LevelError, " cannot store the product", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	tra := map[string]string{
		"pid":      pid,
		"quantity": fmt.Sprintf("%d", quantity),
	}

	go tracking.Log(c, "cart_add", tra)

	l.LogAttrs(c, slog.LevelInfo, "product added in the cart")

	return nil
}

// Get the full session cart.
func Get(c context.Context) (Cart, error) {
	cid := c.Value(contexts.Cart).(string)
	l := slog.With(slog.String("cid", cid))
	l.LogAttrs(c, slog.LevelInfo, "get the cart")

	if !cartExists(c) {
		return Cart{}, errors.New("the session is expired")
	}

	ctx := context.Background()
	values, err := db.Redis.HGetAll(ctx, "cart:"+cid).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot get the cart", slog.String("error", err.Error()))
		return Cart{}, errors.New("something went wrong")
	}

	pds := []products.Product{}
	for key, value := range values {
		product, err := products.Find(ctx, key)

		if err != nil {
			l.LogAttrs(c, slog.LevelError, "cannot get the product", slog.String("id", key), slog.String("error", err.Error()))
			continue
		}

		q, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			l.LogAttrs(c, slog.LevelError, "cannot parse the quantity", slog.String("error", err.Error()))
			continue
		}

		product.Quantity = int(q)

		pds = append(pds, product)
	}

	l.LogAttrs(c, slog.LevelInfo, "got the cart with products", slog.Int("products", len(pds)))

	return Cart{
		ID:       cid,
		Delivery: values["delivery"],
		Payment:  values["payment"],
		Products: pds,
	}, nil
}

// UpdateDelivery update the delivery mode in Redis.
func (c Cart) UpdateDelivery(co context.Context, d string) error {
	l := slog.With(slog.String("cid", c.ID), slog.String("delivery", d))
	l.LogAttrs(co, slog.LevelInfo, "updating the delivery")

	if !orders.IsValidDelivery(co, d) {
		return errors.New("your are not authorized to process this request")
	}

	ctx := context.Background()
	if _, err := db.Redis.HSet(ctx, "cart:"+c.ID, "delivery", d).Result(); err != nil {
		l.LogAttrs(co, slog.LevelError, "cannot update the delivery", slog.String("err", err.Error()))
		return errors.New("something went wrong")
	}

	tra := map[string]string{
		"delivery": d,
	}

	go tracking.Log(co, "cart_delivery", tra)

	l.LogAttrs(co, slog.LevelInfo, "the delivery is updated")

	return nil
}

// UpdatePayment update the payment mode in Redis.
func (c Cart) UpdatePayment(co context.Context, p string) error {
	l := slog.With(slog.String("cid", c.ID), slog.String("payment", p))
	l.LogAttrs(co, slog.LevelInfo, "updating the payment")

	if !orders.IsValidPayment(co, p) {
		return errors.New("your are not authorized to process this request")
	}

	ctx := context.Background()
	if _, err := db.Redis.HSet(ctx, "cart:"+c.ID, "payment", p).Result(); err != nil {
		l.LogAttrs(co, slog.LevelError, "cannot update the payment", slog.String("err", err.Error()))
		return errors.New("something went wrong")
	}

	tra := map[string]string{
		"payment": p,
	}

	go tracking.Log(co, "cart_payment", tra)

	l.LogAttrs(co, slog.LevelInfo, "the payment is updated")

	return nil
}

// RefreshCID refreshes a cart ID (CID).
// If the CID does not exist, it will be created,
// with an expiration time.
func RefreshCID(c context.Context, s string, uid int64) (string, error) {
	cid := c.Value(contexts.Cart).(string)
	l := slog.With(slog.String("cid", s), slog.Int64("uid", uid))
	l.LogAttrs(c, slog.LevelInfo, "refreshing cart")

	ctx := context.Background()
	if _, err := db.Redis.Set(ctx, "cart:"+cid+":user", uid, conf.CartDuration).Result(); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot refresh the cart", slog.String("err", err.Error()))
		return "", errors.New("something went wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "the cart is refreshed")

	return cid, nil
}
