// Package carts manage the user cart
package carts

import (
	"artisons/addresses"
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/products"
	"artisons/shops"
	"artisons/users"
	"artisons/validators"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/maps"
)

type Cart struct {
	// The cart id
	ID int

	// "collect" or "home"
	Delivery string

	DeliveryFees float64

	Total float64

	Payment string

	Products []products.Product

	Address addresses.Address
}

// Exists is not linked as a cart method because
// of merge method.
func Exists(ctx context.Context, cid int) bool {
	l := slog.With(slog.Int("cid", cid))
	l.LogAttrs(ctx, slog.LevelInfo, "checking if the cart exists")

	if cid == 0 {
		l.LogAttrs(ctx, slog.LevelInfo, "the cid is empty")
		return false
	}

	ttl, err := db.Redis.TTL(ctx, fmt.Sprintf("cart:%d", cid)).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "error when retrieving ttl for cart", slog.String("error", err.Error()))
		return false
	}

	l.LogAttrs(ctx, slog.LevelInfo, "got ttl response", slog.Int64("ttl", ttl.Nanoseconds()))

	return ttl.Nanoseconds() > 0
}

func NewCartID(ctx context.Context) (int, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "generating a new cart id")

	cid := rand.Int()

	exists := Exists(ctx, cid)
	if exists {
		return NewCartID(ctx)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "new cart id generated", slog.Int("cid", cid))

	return cid, nil
}

// Add a product into a cart with its quantity
// Verify that the cart and the product exists.
func Add(ctx context.Context, cid int, pid string, quantity int) error {
	l := slog.With(slog.String("product_id", pid), slog.Int("quantity", quantity))
	l.LogAttrs(ctx, slog.LevelInfo, "adding a product to the cart")

	if !products.Available(ctx, pid) {
		return errors.New("oops the data is not found")
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HIncrBy(ctx, fmt.Sprintf("cart:%d", cid), pid, int64(quantity))
		rdb.Expire(ctx, fmt.Sprintf("cart:%d", cid), conf.CartDuration)
		return nil
	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot store the cart", slog.String("error", err.Error()))
		return err
	}

	l.LogAttrs(ctx, slog.LevelInfo, "product added in the cart")

	return nil
}

func (c Cart) Validate(ctx context.Context) error {
	l := slog.With()
	l.LogAttrs(ctx, slog.LevelInfo, "validating the cart")

	if !shops.IsValidDelivery(ctx, c.Delivery) {
		return errors.New("you are not authorized to process this request")
	}

	if !shops.IsValidPayment(ctx, c.Payment) {
		return errors.New("you are not authorized to process this request")
	}

	if len(c.Products) == 0 {
		l.LogAttrs(ctx, slog.LevelInfo, "the product list is empty")
		return errors.New("the cart is empty")
	}

	if err := validators.V.Struct(c.Address); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot validate the user", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input:%s", low)
	}

	pids := []string{}
	for _, p := range c.Products {
		pids = append(pids, p.ID)
	}

	if !products.Availables(ctx, pids) {
		l.LogAttrs(ctx, slog.LevelInfo, "no product is available")
		return errors.New("some products are not available anymore")
	}

	var amount float64 = 0

	for _, value := range c.Products {
		amount += float64(value.Quantity) * value.Price
	}

	min, err := shops.MinDelivery(ctx)
	if err != nil {
		return err
	}

	if amount < min {
		l.LogAttrs(ctx, slog.LevelInfo, "the minimum amount is not reached", slog.Float64("amount", amount), slog.Float64("min", shops.Data.Min))
		return errors.New("the minimum amount is not reached")
	}

	return nil
}

func (c Cart) SaveAddress(ctx context.Context, a addresses.Address) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "saving address")

	err := a.Save(ctx, fmt.Sprintf("cart:%d:info", c.ID))
	if err != nil {
		return err
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "address saved successfullly")

	return nil
}

func Delete(ctx context.Context, cid int, pid string, quantity int) error {
	l := slog.With(slog.String("product_id", pid), slog.Int("quantity", quantity))
	l.LogAttrs(ctx, slog.LevelInfo, "deleting a product to the cart")

	q, err := db.Redis.HGet(ctx, fmt.Sprintf("cart:%d", cid), pid).Result()
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

	if qty > int64(quantity) {
		_, err = db.Redis.HIncrBy(ctx, fmt.Sprintf("cart:%d", cid), pid, -int64(quantity)).Result()
	} else {
		_, err = db.Redis.HDel(ctx, fmt.Sprintf("cart:%d", cid), pid).Result()
	}

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot store the cart", slog.String("error", err.Error()))
		return err
	}

	l.LogAttrs(ctx, slog.LevelInfo, "product removed from the cart")

	return nil
}

func Get(ctx context.Context, cid int) (Cart, error) {
	l := slog.With(slog.Int("cid", cid))
	l.LogAttrs(ctx, slog.LevelInfo, "get the cart")

	if !Exists(ctx, cid) {
		slog.LogAttrs(ctx, slog.LevelInfo, "the cart is empty or does not exist")
		return Cart{}, nil
	}

	values, err := db.Redis.HGetAll(ctx, fmt.Sprintf("cart:%d:info", cid)).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the cart", slog.String("error", err.Error()))
		return Cart{}, errors.New("something went wrong")
	}

	qty, err := db.Redis.HGetAll(ctx, fmt.Sprintf("cart:%d", cid)).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the cart products", slog.String("error", err.Error()))
		return Cart{}, errors.New("something went wrong")
	}

	pids := maps.Keys(qty)

	data, err := products.FindAll(ctx, pids)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the products", slog.String("error", err.Error()))
		return Cart{}, errors.New("something went wrong")
	}

	var fees float64 = 0
	if values["delivery_fees"] != "" {
		fees, err = strconv.ParseFloat(values["delivery_fees"], 64)
		if err != nil {
			l.LogAttrs(ctx, slog.LevelError, "cannot parse the delivery fees", slog.String("error", err.Error()))
			return Cart{}, errors.New("something went wrong")
		}
	}

	var total float64 = 0
	if values["total"] != "" {
		total, err = strconv.ParseFloat(values["total"], 64)
		if err != nil {
			l.LogAttrs(ctx, slog.LevelError, "cannot parse the total", slog.String("error", err.Error()))
			return Cart{}, errors.New("something went wrong")
		}
	}

	pds := []products.Product{}

	for _, p := range data {
		q, err := strconv.ParseInt(qty[p.ID], 10, 32)
		if err != nil {
			l.LogAttrs(ctx, slog.LevelError, "cannot parse the quantity", slog.String("error", err.Error()))
			continue
		}

		p.Quantity = int(q)
		pds = append(pds, p)
	}

	l.LogAttrs(ctx, slog.LevelInfo, "got the cart with products", slog.Int("products", len(pds)))

	return Cart{
		ID:           cid,
		Delivery:     values["delivery"],
		DeliveryFees: fees,
		Payment:      values["payment"],
		Address: addresses.Address{
			Lastname:      values["lastname"],
			Firstname:     values["firstname"],
			Street:        values["street"],
			Complementary: values["complementary"],
			Zipcode:       values["zipcode"],
			City:          values["city"],
			Phone:         values["phone"],
		},
		Products: pds,
		Total:    total,
	}, nil
}

// UpdateDelivery update the delivery mode in Redis.
func (c Cart) UpdateDelivery(ctx context.Context, del string) error {
	l := slog.With(slog.String("delivery", del))
	l.LogAttrs(ctx, slog.LevelInfo, "updating the delivery")

	if !shops.IsValidDelivery(ctx, del) {
		return errors.New("you are not authorized to process this request")
	}

	if _, err := db.Redis.HSet(ctx, fmt.Sprintf("cart:%d:info", c.ID), "delivery", del).Result(); err != nil {
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

	if !shops.IsValidPayment(ctx, p) {
		return errors.New("you are not authorized to process this request")
	}

	if _, err := db.Redis.HSet(ctx, fmt.Sprintf("cart:%d:info", c.ID), "payment", p).Result(); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot update the payment", slog.String("err", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the payment is updated")

	return nil
}

// RefreshCID refreshes a cart ID (CID).
// If the CID does not exist, it will be created,
// with an expiration time.
func RefreshCID(ctx context.Context, cid int) error {
	l := slog.With(slog.Int("cid", cid))
	l.LogAttrs(ctx, slog.LevelInfo, "refreshing cart")

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Expire(ctx, fmt.Sprintf("cart:%d", cid), conf.CartDuration)
		rdb.Expire(ctx, fmt.Sprintf("cart:%d:info", cid), conf.CartDuration)

		return nil
	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot merge the cart into redis", slog.String("error", err.Error()))
		return err
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the cart is refreshed")

	return nil
}

func Merge(ctx context.Context, cid int) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "merging cart")

	u, ok := ctx.Value(contexts.User).(users.User)
	if !ok {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot get the user id")
		return errors.New("you are not authorized to process this request")
	}

	acart, err := db.Redis.HGetAll(ctx, fmt.Sprintf("cart:%d", cid)).Result()
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

		rdb.Del(ctx, fmt.Sprintf("cart:%d", cid), fmt.Sprintf("cart:%d:info", cid))
		return nil
	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot merge the cart into redis", slog.String("error", err.Error()))
		return err
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "cart merged")

	return nil
}

func (c Cart) CalculateTotal(ctx context.Context) (float64, error) {
	var total float64 = 0

	for _, value := range c.Products {
		total += float64(value.Quantity) * value.Price
	}

	var fees float64 = 0

	free, err := shops.DeliveryFreeFees(ctx)
	if err != nil {
		return 0, err
	}

	if c.Delivery != "collect" && total < free {
		del, err := shops.DeliveryFees(ctx)
		if err != nil {
			return 0, err
		}

		fees = del
	}

	total += fees

	_, err = db.Redis.HSet(ctx, fmt.Sprintf("cart:%d:info", c.ID), "total", total, "delivery_fees", fees).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot set the cart total")
		return 0, errors.New("something went wrong")
	}

	return total, nil
}
