// Package orders manage the order created on the application
package orders

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/locales"
	"gifthub/notifications/mails"
	"gifthub/products"
	"gifthub/string/stringutil"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slices"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

	// The order note added by the seller
	Notes []Note
}

type Note struct {
	Note      string
	CreatedAt time.Time
}

// IsValidDelivery returns true if the delivery
// is valid. The values can be "collect" or "home".
// The "collect" value can be used only if it's allowed
// in the settings.
func IsValidDelivery(c context.Context, d string) bool {
	l := slog.With(slog.String("delivery", d))
	l.LogAttrs(c, slog.LevelInfo, "checking delivery validity")

	if d != "collect" && d != "home" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the delivery")
		return false
	}

	if d == "home" && !conf.HasHomeDelivery {
		l.LogAttrs(c, slog.LevelInfo, "cannot continue while the home is not activated")
		return false
	}

	l.LogAttrs(c, slog.LevelInfo, "the delivery is valid")

	return true
}

// IsValidPayment returns true if the payment
// is valid. The values can be "card", "cash", "bitcoin" or "wire".
// The payments can be enablee or disabled in the settings.
func IsValidPayment(c context.Context, p string) bool {
	l := slog.With(slog.String("payment", p))
	l.LogAttrs(c, slog.LevelInfo, "checking payment validity")

	if !slices.Contains(Payments, p) {
		l.LogAttrs(c, slog.LevelInfo, "cannot continue while the payment method is not activated")
		return false
	}

	l.LogAttrs(c, slog.LevelInfo, "the payment is valid")

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
func (o Order) Save(c context.Context) (string, error) {
	l := slog.With(slog.String("oid", o.ID))
	l.LogAttrs(c, slog.LevelInfo, "saving the order")

	if !IsValidDelivery(c, o.Delivery) {
		return "", errors.New("unauthorized")
	}

	if !IsValidPayment(c, o.Payment) {
		return "", errors.New("unauthorized")
	}

	if len(o.Products) == 0 {
		l.LogAttrs(c, slog.LevelInfo, "the product list is empty")
		return "", errors.New("cart_empty")
	}

	pids := []string{}
	for key := range o.Products {
		pids = append(pids, key)
	}

	if !products.Availables(c, pids) {
		l.LogAttrs(c, slog.LevelInfo, "no product is available")
		return "", errors.New("cart_empty")
	}

	oid, err := stringutil.Random()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot generate the pid", slog.String("error", err.Error()))
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
		l.LogAttrs(c, slog.LevelError, "cannot save the order in redis", slog.String("error", err.Error()))
		return "", errors.New("something_went_wrong")
	}

	email, err := db.Redis.HGet(ctx, fmt.Sprintf("user:%d", o.UID), "email").Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelWarn, "cannot send an email ", slog.Int64("uid", o.UID))
	} else {
		lang := c.Value(locales.ContextKey).(language.Tag)
		p := message.NewPrinter(lang)
		// todo add more detail about the order
		msg := p.Sprintf("order_created_email", o.ID)
		go mails.Send(ctx, email, msg)
	}

	l.LogAttrs(c, slog.LevelInfo, "the new order is created", slog.String("oid", oid))

	return oid, nil
}

// UpdateStatus updates the order status.
// An error occurs if the status is not a correct value,
// or the order is not found.
// The full order is returned and an notification is expected
// to be sent to the customer.
func UpdateStatus(c context.Context, oid, status string) error {
	l := slog.With(slog.String("oid", oid), slog.String("status", status))
	l.LogAttrs(c, slog.LevelInfo, "updating the order")

	if !slices.Contains(Status, status) {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the status")
		return errors.New("order_bad_status")
	}

	ctx := context.Background()

	if exists, err := db.Redis.Exists(ctx, "order:"+oid).Result(); exists == 0 || err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot find the order")
		return errors.New("order_not_found")
	}

	_, err := db.Redis.HSet(ctx, "order:"+oid, "status", status).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, " error when setting the status order", slog.String("error", err.Error()))
		return errors.New("something_went_wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "the status is updated")

	return nil
}

func Find(c context.Context, oid string) (Order, error) {
	l := slog.With(slog.String("oid", oid))
	l.LogAttrs(c, slog.LevelInfo, "finding the order")

	ctx := context.Background()

	if exists, err := db.Redis.Exists(ctx, "order:"+oid).Result(); exists == 0 || err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot find the order")
		return Order{}, errors.New("order_not_found")
	}

	data, err := db.Redis.HGetAll(ctx, "order:"+oid).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot get the order from redis", slog.String("error", err.Error()))
		return Order{}, errors.New("something_went_wrong")
	}

	o, err := parseOrder(c, data)
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot parse the order", slog.String("error", err.Error()))
		return Order{}, errors.New("something_went_wrong")
	}

	ids, err := db.Redis.SMembers(ctx, "order:"+oid+":notes").Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot parse the order note ids", slog.String("error", err.Error()))
		return Order{}, errors.New("something_went_wrong")
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		for _, id := range ids {
			key := "order:" + oid + ":note:" + id
			n, err := db.Redis.HGetAll(ctx, key).Result()
			if err != nil {
				l.LogAttrs(c, slog.LevelError, "cannot get the order note", slog.String("error", err.Error()), slog.String("id", id))
				continue
			}

			createdAt, err := time.Parse(time.RFC3339, n["created_at"])
			if err != nil {
				l.LogAttrs(c, slog.LevelError, "cannot parse the created at date", slog.String("error", err.Error()), slog.String("id", id), slog.String("created_at", n["created_at"]))
				continue
			}

			o.Notes = append(o.Notes, Note{
				Note:      n["note"],
				CreatedAt: createdAt,
			})
		}

		return nil
	}); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot get the order notes", slog.String("error", err.Error()))
		return Order{}, errors.New("something_went_wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "got the order with notes", slog.Int("notes", len(o.Notes)))

	return o, nil
}

// AddNote create a new note attached to the order
func AddNote(c context.Context, oid, note string) error {
	l := slog.With(slog.String("oid", oid))
	l.LogAttrs(c, slog.LevelInfo, "adding a note")

	ctx := context.Background()

	if note == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the note")
		return errors.New("order_note_required")
	}

	rep, err := db.Redis.Exists(ctx, "order:"+oid).Result()
	if rep == 0 || err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot find the order")
		return errors.New("order_not_found")
	}

	now := time.Now()
	timestamp := time.Now().UnixMilli()

	if _, err = db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		key := fmt.Sprintf("order:%s:note:%d", oid, timestamp)
		rdb.HSet(ctx, key, "note", note)
		rdb.HSet(ctx, key, "created_at", now.Format(time.RFC3339))
		rdb.SAdd(ctx, "order:"+oid+":notes", timestamp)

		return nil
	}); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot store the note", slog.String("error", err.Error()))
		return errors.New("something_went_wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "note added")

	return nil
}

func parseOrder(c context.Context, m map[string]string) (Order, error) {
	l := slog.With(slog.String("user_id", m["uid"]))
	l.LogAttrs(c, slog.LevelInfo, "parsing the order")

	uid, err := strconv.ParseInt(m["uid"], 10, 64)
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot parse the uid", slog.String("error", err.Error()))
		return Order{}, errors.New("something_went_wrong")
	}

	products := make(map[string]int64)
	for key, value := range m {
		if strings.HasPrefix("product:", key) {
			k := strings.Replace(key, "product:", "", 1)

			q, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				l.LogAttrs(c, slog.LevelError, "cannot parse the quantity", slog.String("quantity", value), slog.String("error", err.Error()))
				return Order{}, errors.New("something_went_wrong")
			}

			products[k] = q
		}
	}

	slog.LogAttrs(c, slog.LevelInfo, "order parsed", slog.String("id", m["id"]))

	return Order{
		ID:            m["id"],
		UID:           uid,
		Delivery:      m["delivery"],
		PaymentStatus: m["payment_status"],
		Payment:       m["payment"],
		Status:        m["status"],
		Products:      products,
		Notes:         []Note{},
	}, nil
}
