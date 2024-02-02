// Package carts manage the user cart
package carts

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/orders"
	"artisons/products"
	"artisons/string/stringutil"
	"artisons/tracking"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

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

func cartExists(ctx context.Context) bool {
	cid := ctx.Value(contexts.Cart).(string)
	l := slog.With(slog.String("cid", cid))
	l.LogAttrs(ctx, slog.LevelInfo, "checking if the car exists")

	ttl, err := db.Redis.TTL(ctx, "cart:"+cid).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "error when retrieving ttl for cart", slog.String("error", err.Error()))
		return false
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "got ttl response", slog.Int64("ttl", ttl.Nanoseconds()))

	return ttl.Nanoseconds() > 0
}

// Add a product into a cart with its quantity
// Verify that the cart and the product exists.
func Add(ctx context.Context, pid string, quantity int) (string, error) {
	l := slog.With(slog.String("product_id", pid), slog.Int("quantity", quantity))
	l.LogAttrs(ctx, slog.LevelInfo, "adding a product to the cart")

	cid, ok := ctx.Value(contexts.Cart).(string)
	if !ok {
		var err error
		cid, err = stringutil.Random()
		if err != nil {
			l.LogAttrs(ctx, slog.LevelError, "cannot generate random string")
		}

		ctx = context.WithValue(ctx, contexts.Cart, cid)
	}

	if !products.Available(ctx, pid) {
		return "", errors.New("product_not_found")
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HIncrBy(ctx, "cart:"+cid, pid, int64(quantity))
		rdb.Expire(ctx, "cart:"+cid, conf.CartDuration)

		return nil
	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot store the cart", slog.String("error", err.Error()))
		return "", err
	}

	if _, err := db.Redis.HIncrBy(ctx, "cart:"+cid, pid, int64(quantity)).Result(); err != nil {
		l.LogAttrs(ctx, slog.LevelError, " cannot store the product", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	if conf.EnableTrackingLog {
		tra := map[string]string{
			"pid":      pid,
			"quantity": fmt.Sprintf("%d", quantity),
		}

		go tracking.Log(ctx, "cart_add", tra)
	}

	l.LogAttrs(ctx, slog.LevelInfo, "product added in the cart")

	return cid, nil
}

// Get the full session cart.
func Get(ctx context.Context) (Cart, error) {
	cid := ctx.Value(contexts.Cart).(string)
	l := slog.With(slog.String("cid", cid))
	l.LogAttrs(ctx, slog.LevelInfo, "get the cart")

	if !cartExists(ctx) {
		return Cart{}, nil
	}

	values, err := db.Redis.HGetAll(ctx, "cart:"+cid).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the cart", slog.String("error", err.Error()))
		return Cart{}, errors.New("something went wrong")
	}

	pids := []string{}
	for key := range values {
		pids = append(pids, key)
	}

	tmps, err := products.FindAll(ctx, pids)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the products", slog.String("error", err.Error()))
		return Cart{}, errors.New("something went wrong")
	}

	pds := []products.Product{}

	for _, p := range tmps {
		q, err := strconv.ParseInt(values[p.ID], 10, 32)
		if err != nil {
			l.LogAttrs(ctx, slog.LevelError, "cannot parse the quantity", slog.String("error", err.Error()))
			continue
		}

		p.Quantity = int(q)
		pds = append(pds, p)
	}

	l.LogAttrs(ctx, slog.LevelInfo, "got the cart with products", slog.Int("products", len(tmps)))

	return Cart{
		ID:       cid,
		Delivery: values["delivery"],
		Payment:  values["payment"],
		Products: pds,
	}, nil
}

// UpdateDelivery update the delivery mode in Redis.
func (c Cart) UpdateDelivery(ctx context.Context, d string) error {
	l := slog.With(slog.String("cid", c.ID), slog.String("delivery", d))
	l.LogAttrs(ctx, slog.LevelInfo, "updating the delivery")

	if !orders.IsValidDelivery(ctx, d) {
		return errors.New("you are not authorized to process this request")
	}

	if _, err := db.Redis.HSet(ctx, "cart:"+c.ID, "delivery", d).Result(); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot update the delivery", slog.String("err", err.Error()))
		return errors.New("something went wrong")
	}

	if conf.EnableTrackingLog {
		tra := map[string]string{
			"delivery": d,
		}

		go tracking.Log(ctx, "cart_delivery", tra)
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the delivery is updated")

	return nil
}

// UpdatePayment update the payment mode in Redis.
func (c Cart) UpdatePayment(ctx context.Context, p string) error {
	l := slog.With(slog.String("cid", c.ID), slog.String("payment", p))
	l.LogAttrs(ctx, slog.LevelInfo, "updating the payment")

	if !orders.IsValidPayment(ctx, p) {
		return errors.New("you are not authorized to process this request")
	}

	if _, err := db.Redis.HSet(ctx, "cart:"+c.ID, "payment", p).Result(); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot update the payment", slog.String("err", err.Error()))
		return errors.New("something went wrong")
	}

	if conf.EnableTrackingLog {
		tra := map[string]string{
			"payment": p,
		}

		go tracking.Log(ctx, "cart_payment", tra)
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the payment is updated")

	return nil
}

// RefreshCID refreshes a cart ID (CID).
// If the CID does not exist, it will be created,
// with an expiration time.
func RefreshCID(ctx context.Context, cid string) (string, error) {
	l := slog.With(slog.String("cid", cid))
	l.LogAttrs(ctx, slog.LevelInfo, "refreshing cart")

	if _, err := db.Redis.Expire(ctx, "cart:"+cid, conf.CartDuration).Result(); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot refresh the cart", slog.String("err", err.Error()))
		return "", errors.New("something went wrong")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the cart is refreshed")

	return cid, nil
}
