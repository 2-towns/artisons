// Package orders manage the order created on the application
package orders

import (
	"artisons/addresses"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/notifications/mails"
	"artisons/products"
	"artisons/string/stringutil"
	"artisons/users"
	"artisons/validators"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/maps"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type SearchResults struct {
	Total  int
	Orders []Order
}

type Order struct {
	ID string

	// The user ID
	UID int

	// "collect" or "home"
	Delivery string

	DeliveryFees float64

	// "cash", "card", "bitcoin" or "wire"
	Payment string

	// "payment_validated", "payment_progress", "payment_refused"
	PaymentStatus string

	// "created", "processing", "delivering", "delivered", "canceled"
	Status string

	// The order note added by the seller
	Notes []Note

	Address addresses.Address

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
	Keywords string
	UID      int
	Sorter   string
}

func (o *Order) AssignID(ctx context.Context) error {
	oid, err := stringutil.Random()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot generate the pid", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	o.ID = oid

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
func (o *Order) Save(ctx context.Context, cid int) error {
	l := slog.With()
	l.LogAttrs(ctx, slog.LevelInfo, "saving the order")

	now := time.Now()

	o.CreatedAt = now
	o.UpdatedAt = now

	u, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		o.UID = u.ID
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, "order:"+o.ID,
			"id", o.ID,
			"uid", o.UID,
			"delivery", o.Delivery,
			"payment", o.Payment,
			"payment_status", o.Status,
			"status", "created",
			"address_lastname", o.Address.Lastname,
			"address_firstname", o.Address.Firstname,
			"address_street", o.Address.Street,
			"address_city", o.Address.City,
			"address_complementary", o.Address.Complementary,
			"address_zipcode", o.Address.Zipcode,
			"address_phone", o.Address.Phone,
			"type", "order",
			"total", o.Total,
			"updated_at", now.Unix(),
			"created_at", now.Unix(),
		)

		if o.UID > 0 {
			rdb.HSet(ctx, "order:"+o.ID, "uid", o.UID)
		}

		for _, p := range o.Products {
			rdb.HSet(ctx, "order:"+o.ID+":products", p.ID, p.Quantity)
		}

		rdb.Del(ctx, fmt.Sprintf("cart:%d", cid), fmt.Sprintf("cart:%d:inf", cid))

		return nil
	}); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot save the order in redis", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the new order is created", slog.String("oid", o.ID))

	return nil
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

	l.LogAttrs(ctx, slog.LevelInfo, "the status is updated")

	return nil
}

func Find(ctx context.Context, oid string) (Order, error) {
	l := slog.With(slog.String("oid", oid))
	l.LogAttrs(ctx, slog.LevelInfo, "finding the order")

	if oid == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate empty id")
		return Order{}, errors.New("oops the data is not found")
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

	m, err := db.Redis.HGetAll(ctx, "order:"+o.ID+":products").Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot retrieve the order products", slog.String("error", err.Error()))
		return Order{}, errors.New("something went wrong")
	}

	pids := maps.Keys(m)
	pdts, err := products.FindAll(ctx, pids)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot retrieve the order products", slog.String("error", err.Error()))
		return Order{}, errors.New("something went wrong")
	}

	for _, pdt := range pdts {
		qty := m[pdt.ID]
		q, err := strconv.ParseInt(qty, 10, 32)
		if err != nil {
			l.LogAttrs(ctx, slog.LevelError, "cannot parse the quantity", slog.String("quantity", qty), slog.String("error", err.Error()))
			return Order{}, errors.New("something went wrong")
		}

		pdt.Quantity = int(q)

		o.Products = append(o.Products, pdt)
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

	return Order{
		ID:            m["id"],
		UID:           int(uid),
		Delivery:      m["delivery"],
		PaymentStatus: m["payment_status"],
		Payment:       m["payment"],
		Status:        m["status"],
		Address: addresses.Address{
			Lastname:      m["address_lastname"],
			Firstname:     m["address_firstname"],
			City:          m["address_city"],
			Street:        m["address_street"],
			Complementary: m["address_complementary"],
			Zipcode:       m["address_zipcode"],
			Phone:         m["address_phone"],
		},
		Notes:     []Note{},
		CreatedAt: time.Unix(createdAt, 0),
		UpdatedAt: time.Unix(updatedAt, 0),
		Total:     total,
	}, nil
}

func Search(ctx context.Context, q Query, offset, num int) (SearchResults, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "searching orders")

	qs := fmt.Sprintf("FT.SEARCH %s @type:{order}", db.OrderIdx)

	if q.Keywords != "" {
		k := db.SearchValue(q.Keywords)
		qs += fmt.Sprintf("(@id:{%s})|(@status:{%s})|(@delivery:{%s})|(@payment:{%s})", k, k, k, k)
	}

	if q.UID != 0 {
		qs += fmt.Sprintf("(@uid:{%d})", q.UID)
	}

	sorter := "updated_at"
	if q.Sorter == "created_at" {
		sorter = "created_at"
	}

	qs += fmt.Sprintf(" SORTBY %s desc LIMIT %d %d DIALECT 2", sorter, offset, num)

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
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the order", slog.Any("order", data), slog.String("error", err.Error()))
			continue
		}

		orders = append(orders, order)
	}

	return SearchResults{
		Total:  int(total),
		Orders: orders,
	}, nil
}
