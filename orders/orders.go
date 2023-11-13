// Package orders manage the order created on the application
package orders

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/contexts"
	"gifthub/notifications/mails"
	"gifthub/products"
	"gifthub/stats"
	"gifthub/string/stringutil"
	"gifthub/users"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jedib0t/go-pretty/v6/table"
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

	DeliveryCost float32

	// "cash", "card", "bitcoin" or "wire"
	Payment string

	// "payment_validated", "payment_progress", "payment_refused"
	PaymentStatus string

	// "created", "processing", "delivering", "delivered", "canceled"
	Status string

	// The key contains the product id and the value is the quantity.
	// Quantities are only filled for the input data.
	// To retrieve the order products, use .Products method
	Quantities map[string]int

	// The order note added by the seller
	Notes []Note

	Address users.Address

	CreatedAt int64
	UpdateAt  int64
}

type Note struct {
	Note      string
	CreatedAt int64
}

type Query struct {
	Status   string
	Delivery string
	Payment  string
}

// IsValidDelivery returns true if the delivery
// is valid. The values can be "collect" or "home".
// The "collect" value can be used only if it's allowed
// in the settings.
func IsValidDelivery(c context.Context, d string) bool {
	l := slog.With(slog.String("delivery", d))
	l.LogAttrs(c, slog.LevelInfo, "checking delivery validity")

	v := validator.New()
	if err := v.Var(d, "oneof=collect home"); err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate  the delivery", slog.String("error", err.Error()))
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

	v := validator.New()
	if err := v.Var(c, "oneof=cash wire bitcoin card"); err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate  te payment", slog.String("error", err.Error()))
		return false
	}

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
// The data are stored like this:
// - order:ID => the order data
// - order:ID product:ID => the product quantity
// - user:ID:orders => the order id added in the set
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

	if len(o.Quantities) == 0 {
		l.LogAttrs(c, slog.LevelInfo, "the product list is empty")
		return "", errors.New("cart_empty")
	}

	v := validator.New()
	if err := v.Struct(o.Address); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot validate the user", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return "", fmt.Errorf("address_%s_required", low)
	}

	pids := []string{}
	for key := range o.Quantities {
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
			"address_lastname", o.Address.Lastname,
			"address_firstname", o.Address.Firstname,
			"address_street", o.Address.Street,
			"address_city", o.Address.City,
			"address_complementary", o.Address.Complementary,
			"address_zipcode", o.Address.Zipcode,
			"address_phone", o.Address.Phone,
			"updated_at", now.Unix(),
			"created_at", now.Unix(),
		)

		for key, value := range o.Quantities {
			rdb.HSet(ctx, "order:"+oid+":products", key, value)
		}

		rdb.ZAdd(ctx, fmt.Sprintf("user:%d:orders", o.UID), redis.Z{
			Score:  float64(now.Unix()),
			Member: oid,
		})

		return nil
	}); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot save the order in redis", slog.String("error", err.Error()))
		return "", errors.New("something_went_wrong")
	}

	go stats.Order(c, oid)

	for pid, quantity := range o.Quantities {
		go stats.SoldProduct(c, oid, pid, quantity)
	}

	go o.SendConfirmationEmail(c)

	l.LogAttrs(c, slog.LevelInfo, "the new order is created", slog.String("oid", oid))

	return oid, nil
}

func (o Order) Total(p []products.Product) float32 {
	total := o.DeliveryCost
	for _, value := range p {
		total += float32(value.Quantity) * value.Price
	}

	return total
}

func (o Order) SendConfirmationEmail(c context.Context) (string, error) {
	l := slog.With(slog.String("oid", o.ID))
	l.LogAttrs(c, slog.LevelInfo, "sending confirmation email")

	email, err := db.Redis.HGet(c, fmt.Sprintf("user:%d", o.UID), "email").Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelWarn, "cannot get the email", slog.Int64("uid", o.UID), slog.String("error", err.Error()))
		return "", err
	}

	lang := c.Value(contexts.Locale).(language.Tag)
	p := message.NewPrinter(lang)

	msg := p.Sprintf("order_created_email", o.Address.Firstname)
	msg += p.Sprintf("order_id_email", o.ID)
	msg += p.Sprintf("order_date_email", time.Unix(o.CreatedAt, 0).Format("Monday, January 1"))

	pds, err := o.Products(c)
	if err != nil {
		l.LogAttrs(c, slog.LevelWarn, "cannot get the products", slog.Int64("uid", o.UID))
		return "", err
	}

	total := o.Total(pds)
	msg += p.Sprintf("order_total_email", total)

	t := table.NewWriter()
	buf := new(bytes.Buffer)
	t.SetOutputMirror(buf)
	t.AppendHeader(table.Row{p.Sprintf("order_title"), p.Sprintf("order_quantity"), p.Sprintf("order_price"), p.Sprintf("order_total"), p.Sprintf("order_link")})

	for _, value := range pds {
		t.AppendRow([]interface{}{value.Title, value.Quantity, value.Price, float32(value.Quantity) * value.Price, value.URL()})
	}

	t.Render()

	msg += buf.String()

	msg += p.Sprintf("order_footer_email")

	err = mails.Send(c, email, msg)
	if err != nil {
		l.LogAttrs(c, slog.LevelWarn, "cannot send the email", slog.String("error", err.Error()))
		return "", err
	}

	return msg, nil
}

// UpdateStatus updates the order status.
// An error occurs if the status is not a correct value,
// or the order is not found.
// The full order is returned and an notification is expected
// to be sent to the customer.
// The keys are :
// - order:oid => the order data
func UpdateStatus(c context.Context, oid, status string) error {
	l := slog.With(slog.String("oid", oid), slog.String("status", status))
	l.LogAttrs(c, slog.LevelInfo, "updating the order")

	v := validator.New()
	if err := v.Var(c, "oneof=created processing delivering delivered canceled"); err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate  the delivery", slog.String("error", err.Error()))
		return errors.New("order_bad_status")
	}

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
		l.LogAttrs(c, slog.LevelError, " cannot update the status order", slog.String("error", err.Error()))
		return errors.New("something_went_wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "the status is updated")

	return nil
}

func (o Order) Products(c context.Context) ([]products.Product, error) {
	l := slog.With(slog.String("oid", o.ID))

	ctx := context.Background()

	m, err := db.Redis.HGetAll(ctx, "order:"+o.ID+":products").Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot retrieve the order products", slog.String("error", err.Error()))
		return []products.Product{}, errors.New("something_went_wrong")
	}

	op := []products.Product{}

	for key, value := range m {
		product, err := products.Find(c, key)
		if err != nil {
			l.LogAttrs(c, slog.LevelError, "cannot retrieve the product", slog.String("pid", key), slog.String("error", err.Error()))
			return []products.Product{}, errors.New("something_went_wrong")
		}

		q, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			l.LogAttrs(c, slog.LevelError, "cannot parse the quantity", slog.String("quantity", value), slog.String("error", err.Error()))
			return []products.Product{}, errors.New("something_went_wrong")
		}

		product.Quantity = int(q)

		op = append(op, product)
	}

	return op, nil
}

func Find(c context.Context, oid string) (Order, error) {
	l := slog.With(slog.String("oid", oid))
	l.LogAttrs(c, slog.LevelInfo, "finding the order")

	if oid == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate empty id")
		return Order{}, errors.New("order_id_required")
	}

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

	if _, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for _, id := range ids {
			key := "order:" + oid + ":note:" + id
			n, err := db.Redis.HGetAll(ctx, key).Result()
			if err != nil {
				l.LogAttrs(c, slog.LevelError, "cannot get the order note", slog.String("error", err.Error()), slog.String("id", id))
				continue
			}

			createdAt, err := strconv.ParseInt(n["created_at"], 10, 64)
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
// The keys are:
// - order:oid:note:nid => the note data
// - order:oid:notes => the note id list
func AddNote(c context.Context, oid, note string) error {
	l := slog.With(slog.String("oid", oid))
	l.LogAttrs(c, slog.LevelInfo, "adding a note")

	ctx := context.Background()

	if note == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the note")
		return errors.New("input_note_required")
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
		rdb.HSet(ctx, key, "created_at", now.Unix())
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

	createdAt, err := strconv.ParseInt(m["created_at"], 10, 64)
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot parse the created_at", slog.String("created_at", m["created_at"]), slog.String("error", err.Error()))
		return Order{}, errors.New("something_went_wrong")
	}

	updatedAt, err := strconv.ParseInt(m["updated_at"], 10, 64)
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot parse the updated_at", slog.String("updated_at", m["updated_at"]), slog.String("error", err.Error()))
		return Order{}, errors.New("something_went_wrong")
	}

	slog.LogAttrs(c, slog.LevelInfo, "order parsed", slog.String("id", m["id"]))

	return Order{
		ID:            m["id"],
		UID:           uid,
		Delivery:      m["delivery"],
		PaymentStatus: m["payment_status"],
		Payment:       m["payment"],
		Status:        m["status"],
		Address: users.Address{
			Lastname:      m["address_lastname"],
			Firstname:     m["address_firstname"],
			City:          m["address_city"],
			Street:        m["address_street"],
			Complementary: m["address_complementary"],
			Zipcode:       m["address_zipcode"],
			Phone:         m["address_phone"],
		},
		Quantities: map[string]int{},
		Notes:      []Note{},
		CreatedAt:  createdAt,
		UpdateAt:   updatedAt,
	}, nil
}

func Search(c context.Context, q Query) ([]Order, error) {
	slog.LogAttrs(c, slog.LevelInfo, "searching orders")

	qs := ""
	if q.Status != "" {
		qs += fmt.Sprintf("@status:{%s}", q.Status)
	}

	if q.Delivery != "" {
		qs += fmt.Sprintf("@delivery:{%s}", q.Delivery)
	}

	if q.Payment != "" {
		qs += fmt.Sprintf("@payment:{%s}", q.Payment)
	}

	slog.LogAttrs(c, slog.LevelInfo, "preparing redis request", slog.String("query", qs))

	ctx := context.Background()
	cmds, err := db.Redis.Do(
		ctx,
		"FT.SEARCH",
		db.OrderIdx,
		qs,
	).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot run the search query", slog.String("query", qs), slog.String("error", err.Error()))
		return []Order{}, err
	}

	res := cmds.(map[interface{}]interface{})
	slog.LogAttrs(c, slog.LevelInfo, "search done", slog.Int64("results", res["total_results"].(int64)))

	results := res["results"].([]interface{})
	orders := []Order{}

	for _, value := range results {
		m := value.(map[interface{}]interface{})
		attributes := m["extra_attributes"].(map[interface{}]interface{})
		data := db.ConvertMap(attributes)

		order, err := parseOrder(c, data)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot order the product", slog.Any("order", data), slog.String("error", err.Error()))
			continue
		}

		orders = append(orders, order)
	}

	return orders, nil
}
