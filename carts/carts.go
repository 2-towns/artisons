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

func cartExists(cid string) bool {
	ctx := context.Background()
	ttl, err := db.Redis.TTL(ctx, "cart:"+cid).Result()
	if err != nil {
		log.Printf("sequence_fail: error when retrieving ttl for cart:%s %s", cid, err.Error())
		return false
	}

	return ttl.Nanoseconds() > 0
}

// Add a product into a cart with its quantity
// Verify that the cart and the product exists.
func Add(cid, pid string, quantity int64) error {
	if !cartExists(cid) {
		return errors.New("cart_not_found")
	}

	if !products.Available(pid) {
		return errors.New("product_not_found")
	}

	ctx := context.Background()
	if _, err := db.Redis.HIncrBy(ctx, "cart:"+cid, "product:"+pid, quantity).Result(); err != nil {
		log.Printf("sequence_fail: error when storing product:%s in cart:%s %s", pid, cid, err.Error())
		return errors.New("something_went_wrong")
	}

	return nil
}

// Get the full session cart.
// TODO: Get the product detail from the products package
func Get(cid string) (Cart, error) {
	if !cartExists(cid) {
		return Cart{}, errors.New("cart_not_found")
	}

	ctx := context.Background()
	values, err := db.Redis.HGetAll(ctx, "cart:"+cid).Result()
	if err != nil {
		log.Printf("sequence_fail: error when getting the cart:%s %s", cid, err.Error())
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

	return Cart{
		ID:       cid,
		Delivery: values["delivery"],
		Payment:  values["payment"],
		Products: products,
	}, nil
}

// UpdateDelivery update the delivery mode in Redis.
func (c Cart) UpdateDelivery(d string) error {
	if !orders.IsValidDelivery(d) {
		return errors.New("unauthorized")
	}

	ctx := context.Background()
	if _, err := db.Redis.HSet(ctx, "cart:"+c.ID, "delivery", d).Result(); err != nil {
		log.Printf("sequence_fail: error when updating delivery %s in cart:%s %s", d, c.ID, err.Error())
		return errors.New("something_went_wrong")
	}

	return nil
}

// UpdatePayment update the payment mode in Redis.
func (c Cart) UpdatePayment(p string) error {
	if !orders.IsValidPayment(p) {
		log.Printf("WARN: input_validation_fail: the payment value is wrong %s", p)
		return errors.New("unauthorized")
	}

	ctx := context.Background()
	if _, err := db.Redis.HSet(ctx, "cart:"+c.ID, "payment", p).Result(); err != nil {
		log.Printf("sequence_fail: error when updating payment %s in cart:%s %s", p, c.ID, err.Error())
		return errors.New("something_went_wrong")
	}

	return nil
}

// RefreshCID refreshes a cart ID (CID).
// If the CID does not exist, it will be created,
// with an expiration time.
func RefreshCID(s string) (string, error) {
	cid := s

	if cid == "" {
		id, err := stringutil.Random()
		if err != nil {
			log.Printf("sequence_fail: error when generating cart:%s %s", cid, err.Error())
			return "", errors.New("something_went_wrong")
		}
		cid = id
	}

	ctx := context.Background()
	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, "cart:"+cid, "cid", cid)
		rdb.Expire(ctx, "cart:"+cid, conf.CartDuration)

		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return "", errors.New("something_went_wrong")
	}

	return cid, nil
}
