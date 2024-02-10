// Package carts manage the user cart
package carts

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/orders"
	"artisons/products"
	"artisons/string/stringutil"
	"artisons/users"
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

// Exists is not linked as a cart method because
// of merge method.
func Exists(ctx context.Context, cid string) bool {
	l := slog.With(slog.String("cid", cid))
	l.LogAttrs(ctx, slog.LevelInfo, "checking if the cart exists")

	if cid == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "the cid is empty")
		return false
	}

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
func (c Cart) Add(ctx context.Context, pid string, quantity int) (string, error) {
	l := slog.With(slog.String("product_id", pid), slog.Int("quantity", quantity))
	l.LogAttrs(ctx, slog.LevelInfo, "adding a product to the cart")

	cid := c.ID
	if cid == "" {
		id, err := stringutil.Random()
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot generate random string")
			return "", err
		}

		cid = id

		slog.LogAttrs(ctx, slog.LevelInfo, "new cart id generated", slog.String("cid", cid))
	}

	if !products.Available(ctx, pid) {
		return "", errors.New("oops the data is not found")
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HIncrBy(ctx, "cart:"+cid, pid, int64(quantity))
		rdb.Expire(ctx, "cart:"+cid, conf.CartDuration)

		return nil
	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot store the cart", slog.String("error", err.Error()))
		return "", err
	}

	l.LogAttrs(ctx, slog.LevelInfo, "product added in the cart")

	if c.ID == "" {
		return cid, nil
	}

	return "", nil
}

func (c Cart) Delete(ctx context.Context, pid string, quantity int) error {
	l := slog.With(slog.String("product_id", pid), slog.Int("quantity", quantity))
	l.LogAttrs(ctx, slog.LevelInfo, "deleting a product to the cart")

	q, err := db.Redis.HGet(ctx, "cart:"+c.ID, pid).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the quantity", slog.String("error", err.Error()))

		if err.Error() == "redis: nil" {
			return errors.New("oops the data is not found")
		}

		return err
	}

	qty, err := strconv.ParseInt(q, 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the quantity", slog.String("error", err.Error()))
		return err
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		if qty > int64(quantity) {
			rdb.HIncrBy(ctx, "cart:"+c.ID, pid, -int64(quantity))
		} else {
			rdb.HDel(ctx, "cart:"+c.ID, pid)
		}

		rdb.Expire(ctx, "cart:"+c.ID, conf.CartDuration)

		return nil
	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot store the cart", slog.String("error", err.Error()))
		return err
	}

	l.LogAttrs(ctx, slog.LevelInfo, "product removed from the cart")

	return nil
}

// Get the full session cart.
func (c Cart) Get(ctx context.Context) (Cart, error) {
	l := slog.With(slog.String("cid", c.ID))
	l.LogAttrs(ctx, slog.LevelInfo, "get the cart")

	if c.ID == "" {
		return Cart{}, nil
	}

	if !Exists(ctx, c.ID) {
		slog.LogAttrs(ctx, slog.LevelInfo, "the cart is empty or does not exist")
		return Cart{}, nil
	}

	values, err := db.Redis.HGetAll(ctx, "cart:"+c.ID).Result()
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
		ID:       c.ID,
		Delivery: values["delivery"],
		Payment:  values["payment"],
		Products: pds,
	}, nil
}

// UpdateDelivery update the delivery mode in Redis.
func (c Cart) UpdateDelivery(ctx context.Context, del string) error {
	l := slog.With(slog.String("delivery", del))
	l.LogAttrs(ctx, slog.LevelInfo, "updating the delivery")

	if !orders.IsValidDelivery(ctx, del) {
		return errors.New("you are not authorized to process this request")
	}

	if _, err := db.Redis.HSet(ctx, "cart:"+c.ID, "delivery", del).Result(); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot update the delivery", slog.String("err", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the delivery is updated")

	return nil
}

// UpdatePayment update the payment mode in Redis.
func (c Cart) UpdatePayment(ctx context.Context, p string) error {
	l := slog.With(slog.String("payment", p))
	l.LogAttrs(ctx, slog.LevelInfo, "updating the payment")

	if !orders.IsValidPayment(ctx, p) {
		return errors.New("you are not authorized to process this request")
	}

	if _, err := db.Redis.HSet(ctx, "cart:"+c.ID, "payment", p).Result(); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot update the payment", slog.String("err", err.Error()))
		return errors.New("something went wrong")
	}

	// TODO move to the payment route
	// if conf.EnableTrackingLog {
	// 	tra := map[string]string{
	// 		"payment": p,
	// 	}

	// 	go tracking.Log(ctx, "cart_payment", tra)
	// }

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

func Merge(ctx context.Context, cid string) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "merging cart")

	u, ok := ctx.Value(contexts.User).(users.User)
	if !ok {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot get the user id")
		return errors.New("you are not authorized to process this request")
	}

	acart, err := db.Redis.HGetAll(ctx, "cart:"+cid).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot get the anonymous cart items")
		return errors.New("something went wrong")
	}

	ucart, err := db.Redis.HGetAll(ctx, fmt.Sprintf("cart:%d", u.ID)).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot get the anonymous cart items")
		return errors.New("something went wrong")
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		for key, val := range acart {
			a, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot parse the anonymous quantity", slog.String("product", key), slog.String("quantity", val))
				continue
			}

			if ucart[key] != "" {
				b, err := strconv.ParseInt(ucart[key], 10, 64)
				if err != nil {
					slog.LogAttrs(ctx, slog.LevelError, "cannot parse the existing quantity", slog.String("quantity", val))
					continue
				}

				rdb.HSet(ctx, fmt.Sprintf("cart:%d", u.ID), key, a+b)
			} else {
				rdb.HSet(ctx, fmt.Sprintf("cart:%d", u.ID), key, a)
			}
		}

		rdb.Del(ctx, "cart:"+cid)
		return nil
	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot merge the cart into redis", slog.String("error", err.Error()))
		return err
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "cart  merged")

	return nil
}
