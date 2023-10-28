// Package orders manage the order created on the application
package orders

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/products"
	"gifthub/string/stringutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slices"
)

var Status = []string{"created", "processing", "delivering", "delivered", "canceled"}

var Payments = []string{"cash", "wire", "bitcoin", "card"}

type Order struct {
	ID string

	// The user ID
	UID int64

	// "collect" or "home"
	Delivery string

	// "cash", "card", "bitcoin" or "wire"
	Payment string

	// "payment_validated", "payment_progress", "payment_refused"
	PaymentStatus string

	// "created", "processing", "delivering", "delivered", "canceled"
	Status string

	// The key contains the product id and the value
	// Is the quantity
	Products map[string]int64
}

// IsValidDelivery returns true if the delivery
// is valid. The values can be "collect" or "home".
// The "collect" value can be used only if it's allowed
// in the settings.
func IsValidDelivery(d string) bool {
	if d != "collect" && d != "home" {
		log.Printf("WARN: input_validation_fail: the delivery value is wrong %s", d)
		return false
	}

	if d == "home" && !conf.HasHomeDelivery {
		log.Printf("WARN: input_validation_fail: the home delivery is not activated %s", d)
		return false
	}

	return true
}

// IsValidPayment returns true if the payment
// is valid. The values can be "card", "cash", "bitcoin" or "wire".
// The payments can be enablee or disabled in the settings.
func IsValidPayment(p string) bool {
	if !slices.Contains(Payments, p) {
		log.Printf("WARN: input_validation_fail: the payment method is not activated %s", p)
		return false
	}

	return true
}

// Add create an order into Redis.
// The default order status is "created".
// The default payment_status is "payment_progress".
// The order ID is a random string and returned if it succeed.
// The products are stored as the cart, the key is the
// product id and the value is the quantity.
// An error occurs if the delivery or the payment values are invalid,
// if the product list is empty, or one of the product is not available.
func (o Order) Save() (string, error) {
	if !IsValidDelivery(o.Delivery) {
		log.Printf("WARN: input_validation_fail: the delivery value %s is wrong", o.Delivery)
		return "", errors.New("unauthorized")
	}

	if !IsValidPayment(o.Payment) {
		log.Printf("WARN: input_validation_fail: the payment value %s is wrong", o.Payment)
		return "", errors.New("unauthorized")
	}

	if len(o.Products) == 0 {
		log.Printf("input_validation_fail: the products are empty")
		return "", errors.New("cart_empty")
	}

	pids := []string{}
	for key := range o.Products {
		pids = append(pids, key)
	}

	if !products.Availables(pids) {
		return "", errors.New("cart_empty")
	}

	oid, err := stringutil.Random()
	if err != nil {
		log.Printf("sequence_fail: error when generating the order id %v", err)
		return "", errors.New("something_went_wrong")
	}

	now := time.Now()
	ctx := context.Background()
	if _, err = db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, "order:"+oid,
			"id", oid,
			"uid", o.UID,
			"delivery", o.Delivery,
			"payment", o.Payment,
			"payment_status", "payment_progress",
			"status", "created",
			"updated_at", now.Format(time.RFC3339),
			"created_at", now.Format(time.RFC3339),
		)

		for key, value := range o.Products {
			rdb.HSet(ctx, "order:"+oid, "product:"+key, value)
		}

		rdb.HSet(ctx, fmt.Sprintf("user:%d", o.UID), "order:"+oid, oid)

		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return "", errors.New("something_went_wrong")
	}

	return oid, nil
}

// UpdateStatus updates the order status.
// An error occurs if the status is not a correct value,
// or the order is not found.
// The full order is returned and an notification is expected
// to be sent to the customer.
func UpdateStatus(oid, status string) (Order, error) {
	if !slices.Contains(Status, status) {
		log.Printf("WARN: input_validation_fail: the status value is wrong %s", status)
		return Order{}, errors.New("unauthorized")
	}

	ctx := context.Background()
	data, err := db.Redis.HGetAll(ctx, "order:"+oid).Result()
	if err != nil {
		log.Printf("sequence_fail: error when getting the order %s %s", oid, err.Error())
		return Order{}, errors.New("something_went_wrong")
	}

	if data["id"] == "" {
		log.Printf("input_validation_fail: the order %s does not exist", oid)
		return Order{}, errors.New("order_not_found")
	}

	_, err = db.Redis.HSet(ctx, "order:"+oid, "status", status).Result()
	if err != nil {
		log.Printf("sequence_fail: error when setting the status order %s %s", oid, err.Error())
		return Order{}, errors.New("something_went_wrong")
	}

	o, err := parseOrder(data)
	if err != nil {
		log.Printf("sequence_fail: error when parsing the order %s %s", oid, err.Error())
		return Order{}, errors.New("something_went_wrong")
	}

	o.Status = status

	return o, nil
}

// AddNote create a new message attached to the order
func AddNote(oid, message string) error {
	if message == "" {
		log.Printf("input_validation_fail: the message is empty")
		return errors.New("order_note_required")
	}

	ctx := context.Background()

	rep, err := db.Redis.Exists(ctx, "order:"+oid).Result()
	if rep == 0 || err != nil {
		log.Printf("input_validation_fail: The order %s is not found", oid)
		return errors.New("order_not_found")
	}

	now := time.Now()
	timestamp := time.Now().UnixMilli()

	if _, err = db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		key := fmt.Sprintf("ordernote:%s%d", oid, timestamp)
		rdb.HSet(ctx, key, "message", message)
		rdb.HSet(ctx, key, "created_at", now.Format(time.RFC3339))
		rdb.SAdd(ctx, "ordernote:"+oid, timestamp)

		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return errors.New("something_went_wrong")
	}

	return nil
}

func parseOrder(m map[string]string) (Order, error) {
	id, err := strconv.ParseInt(m["uid"], 10, 32)
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when parsing uid %s %s", m["id"], err.Error())
		return Order{}, errors.New("something_went_wrong")
	}

	products := make(map[string]int64)
	for key, value := range m {
		if strings.HasPrefix("product:", key) {
			k := strings.Replace(key, "product:", "", 1)

			q, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				log.Printf("ERROR: sequence_fail: error when parsing quantity %s %s", value, err.Error())
				return Order{}, errors.New("something_went_wrong")
			}

			products[k] = q
		}
	}

	return Order{
		ID:            m["id"],
		UID:           id,
		Delivery:      m["delivery"],
		PaymentStatus: m["payment_status"],
		Payment:       m["payment"],
		Status:        m["status"],
		Products:      products,
	}, nil
}
