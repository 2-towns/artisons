// Package orders manage the order created on the application
package orders

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/notifications/mails"
	"artisons/products"
	"artisons/stats"
	"artisons/string/stringutil"
	"artisons/tracking"
	"artisons/users"
	"artisons/validators"
	"bytes"
	"context"
	"errors"
	"fmt"
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

type SearchResults struct {
	Total  int
	Orders []Order
}

var Payments = []string{"cash", "wire", "bitcoin", "card"}

type Order struct {
	ID string

	// The user ID
	UID int

	// "collect" or "home"
	Delivery string

	DeliveryCost float64

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

	CreatedAt time.Time
	UpdatedAt time.Time

	Products []products.Product

	Total float64
}

type Note struct {
	Note      string
	CreatedAt time.Time
}

type Query struct {
	Keyword string
}

// IsValidDelivery returns true if the delivery
// is valid. The values can be "collect" or "home".
// The "collect" value can be used only if it's allowed
// in the settings.
func IsValidDelivery(ctx context.Context, d string) bool {
	l := slog.With(slog.String("delivery", d))
	l.LogAttrs(ctx, slog.LevelInfo, "checking delivery validity")

	if err := validators.V.Var(d, "oneof=collect home"); err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate  the delivery", slog.String("error", err.Error()))
		return false
	}

	if d == "home" && !conf.HasHomeDelivery {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot continue while the home is not activated")
		return false
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the delivery is valid")

	return true
}

// IsValidPayment returns true if the payment
// is valid. The values can be "card", "cash", "bitcoin" or "wire".
// The payments can be enablee or disabled in the settings.
func IsValidPayment(ctx context.Context, p string) bool {
	l := slog.With(slog.String("payment", p))
	l.LogAttrs(ctx, slog.LevelInfo, "checking payment validity")

	if err := validators.V.Var(ctx, "oneof=cash wire bitcoin card"); err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate  te payment", slog.String("error", err.Error()))
		return false
	}

	if !slices.Contains(Payments, p) {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot continue while the payment method is not activated")
		return false
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the payment is valid")

	return true
}

func (o Order) Validate(ctx context.Context) error {
	l := slog.With()
	l.LogAttrs(ctx, slog.LevelInfo, "saving the order")

	if !IsValidDelivery(ctx, o.Delivery) {
		return errors.New("your are not authorized to process this request")
	}

	if !IsValidPayment(ctx, o.Payment) {
		return errors.New("your are not authorized to process this request")
	}

	if len(o.Quantities) == 0 {
		l.LogAttrs(ctx, slog.LevelInfo, "the product list is empty")
		return errors.New("the cart is empty")
	}

	if err := validators.V.Struct(o.Address); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot validate the user", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("address_%s_required", low)
	}

	return nil
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
func (o Order) Save(ctx context.Context) (string, error) {
	l := slog.With()
	l.LogAttrs(ctx, slog.LevelInfo, "saving the order")

	tra := map[string]string{
		"uid":      fmt.Sprintf("%d", o.UID),
		"delivery": o.Delivery,
		"payment":  o.Payment,
	}

	pids := []string{}
	for key, q := range o.Quantities {
		pids = append(pids, key)
		tra[key] = fmt.Sprintf("%d", q)
	}

	if !products.Availables(ctx, pids) {
		l.LogAttrs(ctx, slog.LevelInfo, "no product is available")
		return "", errors.New("the cart is empty")
	}

	oid, err := stringutil.Random()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot generate the pid", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	tra["oid"] = oid

	o, err = o.WithProducts(ctx)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelWarn, "cannot get the products", slog.Int("uid", o.UID))
		return "", err
	}

	o = o.WithTotal()
	now := time.Now()

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
			"total", o.Total,
			"updated_at", now.Unix(),
			"created_at", now.Unix(),
		)

		// Use for Redis Search in order to restrict the items
		rdb.HSetNX(ctx, "order:"+oid, "type", "order")

		for key, value := range o.Quantities {
			rdb.HSet(ctx, "order:"+oid+":products", key, value)
		}

		rdb.ZAdd(ctx, fmt.Sprintf("user:%d:orders", o.UID), redis.Z{
			Score:  float64(now.Unix()),
			Member: oid,
		})

		return nil
	}); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot save the order in redis", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	go stats.Order(ctx, oid, o.Quantities, o.Total)

	go o.SendConfirmationEmail(ctx)

	go tracking.Log(ctx, "order", tra)

	l.LogAttrs(ctx, slog.LevelInfo, "the new order is created", slog.String("oid", oid))

	return oid, nil
}

func (o Order) WithTotal() Order {
	total := o.DeliveryCost
	for _, value := range o.Products {
		total += float64(value.Quantity) * value.Price
	}

	o.Total = total

	return o
}

func (o Order) SendConfirmationEmail(ctx context.Context) (string, error) {
	l := slog.With(slog.String("oid", o.ID))
	l.LogAttrs(ctx, slog.LevelInfo, "sending confirmation email")

	email, err := db.Redis.HGet(ctx, fmt.Sprintf("user:%d", o.UID), "email").Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelWarn, "cannot get the email", slog.Int("uid", o.UID), slog.String("error", err.Error()))
		return "", err
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	p := message.NewPrinter(lang)

	msg := p.Sprintf("email_order_confirmation", o.Address.Firstname)
	msg += p.Sprintf("email_order_confirmationid", o.ID)
	msg += p.Sprintf("email_order_confirmationdate", o.CreatedAt.Format("Monday, January 1"))
	msg += p.Sprintf("email_order_confirmationtotal", o.Total)

	t := table.NewWriter()
	buf := new(bytes.Buffer)
	t.SetOutputMirror(buf)
	t.AppendHeader(table.Row{p.Sprintf("title"), p.Sprintf("quality"), p.Sprintf("price"), p.Sprintf("total"), p.Sprintf("link")})

	for _, value := range o.Products {
		t.AppendRow([]interface{}{value.Title, value.Quantity, value.Price, float64(value.Quantity) * value.Price, value.URL()})
	}

	t.Render()

	msg += buf.String()

	msg += p.Sprintf("email_order_confirmationfooter")

	err = mails.Send(ctx, email, p.Sprintf("email_order_subject", o.ID), msg)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelWarn, "cannot send the email", slog.String("error", err.Error()))
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
func UpdateStatus(ctx context.Context, oid, status string) error {
	l := slog.With(slog.String("oid", oid), slog.String("status", status))
	l.LogAttrs(ctx, slog.LevelInfo, "updating the order status")

	if err := validators.V.Var(status, "required,oneof=created processing delivering delivered canceled"); err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate the status", slog.String("error", err.Error()))
		return errors.New("input:status")
	}

	if exists, err := db.Redis.Exists(ctx, "order:"+oid).Result(); exists == 0 || err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot find the order")
		return errors.New("oops the data is not found")
	}

	_, err := db.Redis.HSet(ctx, "order:"+oid, "status", status).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, " cannot update the status order", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	tra := map[string]string{
		"oid":    oid,
		"status": status,
	}

	go tracking.Log(ctx, "order_status", tra)

	l.LogAttrs(ctx, slog.LevelInfo, "the status is updated")

	return nil
}

func (o Order) WithProducts(ctx context.Context) (Order, error) {
	l := slog.With(slog.String("oid", o.ID))

	m, err := db.Redis.HGetAll(ctx, "order:"+o.ID+":products").Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot retrieve the order products", slog.String("error", err.Error()))
		return Order{}, errors.New("something went wrong")
	}

	pds := []products.Product{}

	for key, value := range m {
		product, err := products.Find(ctx, key)
		if err != nil {
			l.LogAttrs(ctx, slog.LevelError, "cannot retrieve the product", slog.String("pid", key), slog.String("error", err.Error()))
			return Order{}, errors.New("something went wrong")
		}

		q, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			l.LogAttrs(ctx, slog.LevelError, "cannot parse the quantity", slog.String("quantity", value), slog.String("error", err.Error()))
			return Order{}, errors.New("something went wrong")
		}

		product.Quantity = int(q)

		pds = append(pds, product)
	}

	o.Products = pds

	return o, nil
}

func Find(ctx context.Context, oid string) (Order, error) {
	l := slog.With(slog.String("oid", oid))
	l.LogAttrs(ctx, slog.LevelInfo, "finding the order")

	if oid == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate empty id")
		return Order{}, errors.New("input:id")
	}

	if exists, err := db.Redis.Exists(ctx, "order:"+oid).Result(); exists == 0 || err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot find the order")
		return Order{}, errors.New("oops the data is not found")
	}

	data, err := db.Redis.HGetAll(ctx, "order:"+oid).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the order from redis", slog.String("error", err.Error()))
		return Order{}, errors.New("something went wrong")
	}

	o, err := parse(ctx, data)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot parse the order", slog.String("error", err.Error()))
		return Order{}, errors.New("something went wrong")
	}

	ids, err := db.Redis.SMembers(ctx, "order:"+oid+":notes").Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot parse the order note ids", slog.String("error", err.Error()))
		return Order{}, errors.New("something went wrong")
	}

	cmds, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for _, id := range ids {
			key := "order:" + oid + ":note:" + id
			rdb.HGetAll(ctx, key)
		}

		return nil
	})

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the order notes", slog.String("error", err.Error()))
		return o, errors.New("something went wrong")
	}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the tag links", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		val := cmd.(*redis.MapStringStringCmd).Val()

		createdAt, err := strconv.ParseInt(val["created_at"], 10, 64)
		if err != nil {
			l.LogAttrs(ctx, slog.LevelError, "cannot parse the created at date", slog.String("error", err.Error()), slog.String("created_at", val["created_at"]))
			continue
		}

		o.Notes = append(o.Notes, Note{
			Note:      val["note"],
			CreatedAt: time.Unix(createdAt, 0),
		})
	}

	l.LogAttrs(ctx, slog.LevelInfo, "got the order with notes", slog.Int("notes", len(o.Notes)))

	return o, nil
}

// AddNote create a new note attached to the order
// The keys are:
// - order:oid:note:nid => the note data
// - order:oid:notes => the note id list
func AddNote(ctx context.Context, oid, note string) error {
	l := slog.With(slog.String("oid", oid))
	l.LogAttrs(ctx, slog.LevelInfo, "adding a note")

	if note == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate the note")
		return errors.New("input:note")
	}

	rep, err := db.Redis.Exists(ctx, "order:"+oid).Result()
	if rep == 0 || err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot find the order")
		return errors.New("oops the data is not found")
	}

	now := time.Now()
	timestamp := time.Now().UnixMilli()

	if _, err = db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		key := fmt.Sprintf("order:%s:note:%d", oid, timestamp)
		rdb.HSet(ctx, key, "created_at", now.Unix(), "note", note)
		rdb.SAdd(ctx, "order:"+oid+":notes", timestamp)

		return nil
	}); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot store the note", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "note added")

	return nil
}

func parse(ctx context.Context, m map[string]string) (Order, error) {
	l := slog.With(slog.String("user_id", m["uid"]))
	l.LogAttrs(ctx, slog.LevelInfo, "parsing the order")

	uid, err := strconv.ParseInt(m["uid"], 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the uid", slog.String("error", err.Error()))
		return Order{}, errors.New("something went wrong")
	}

	createdAt, err := strconv.ParseInt(m["created_at"], 10, 64)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot parse the created_at", slog.String("created_at", m["created_at"]), slog.String("error", err.Error()))
		return Order{}, errors.New("something went wrong")
	}

	updatedAt, err := strconv.ParseInt(m["updated_at"], 10, 64)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot parse the updated_at", slog.String("updated_at", m["updated_at"]), slog.String("error", err.Error()))
		return Order{}, errors.New("something went wrong")
	}

	total, err := strconv.ParseFloat(m["total"], 64)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot parse the total", slog.String("total", m["total"]), slog.String("error", err.Error()))
		return Order{}, errors.New("something went wrong")
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "order parsed", slog.String("id", m["id"]))

	return Order{
		ID:            m["id"],
		UID:           int(uid),
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
		CreatedAt:  time.Unix(createdAt, 0),
		UpdatedAt:  time.Unix(updatedAt, 0),
		Total:      total,
	}, nil
}

func Search(ctx context.Context, q Query, offset, num int) (SearchResults, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "searching orders")

	qs := fmt.Sprintf("FT.SEARCH %s @type:{order}", db.OrderIdx)

	if q.Keyword != "" {
		k := db.SearchValue(q.Keyword)
		qs += fmt.Sprintf("(@id:{%s})|(@status:{%s})|(@delivery:{%s})|(@payment:{%s})", k, k, k, k)
	}

	qs += fmt.Sprintf(" SORTBY updated_at desc LIMIT %d %d DIALECT 2", offset, num)

	slog.LogAttrs(ctx, slog.LevelInfo, "preparing redis request", slog.String("query", qs))

	args, err := db.SplitQuery(ctx, qs)
	if err != nil {
		return SearchResults{}, err
	}

	cmds, err := db.Redis.Do(ctx, args...).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot run the search query", slog.String("error", err.Error()))
		return SearchResults{}, err
	}

	res := cmds.(map[interface{}]interface{})
	total := res["total_results"].(int64)

	slog.LogAttrs(ctx, slog.LevelInfo, "search done", slog.Int64("results", res["total_results"].(int64)))

	results := res["results"].([]interface{})
	orders := []Order{}

	for _, value := range results {
		m := value.(map[interface{}]interface{})
		attributes := m["extra_attributes"].(map[interface{}]interface{})
		data := db.ConvertMap(attributes)

		order, err := parse(ctx, data)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot order the product", slog.Any("order", data), slog.String("error", err.Error()))
			continue
		}

		orders = append(orders, order)
	}

	return SearchResults{
		Total:  int(total),
		Orders: orders,
	}, nil
}
