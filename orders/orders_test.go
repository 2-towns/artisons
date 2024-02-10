package orders

import (
	"artisons/conf"
	"artisons/products"
	"artisons/tests"
	"artisons/users"
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var order Order = Order{
	ID:       "ORD1",
	UID:      1,
	Delivery: "collect",
	Payment:  "cash",
	Total:    105.5,
	Quantities: map[string]int{
		"PDT1": 1,
	},
	DeliveryCost: 5,
	Address: users.Address{
		Firstname:     "Arnaud",
		Lastname:      "Arnaud",
		City:          "Oran",
		Street:        "Hay Yasmine",
		Complementary: "Hay Salam",
		Zipcode:       "31244",
		Phone:         "0559682532",
	},
	CreatedAt: time.Unix(1699628645, 0),
}

var cur string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	cur = path.Dir(filename) + "/"
}

func TestIsValidDelivery(t *testing.T) {
	ctx := tests.Context()

	var tests = []struct {
		name  string
		value string
		valid bool
	}{
		{"collect", "collect", true},
		{"idontexist", "idontexist", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if valid := IsValidDelivery(ctx, tt.value); valid != tt.valid {
				t.Fatalf(`valid = %v, want %v`, valid, tt.valid)
			}
		})
	}
}

func TestIsValidPayment(t *testing.T) {
	ctx := tests.Context()

	var tests = []struct {
		name  string
		value string
		valid bool
	}{
		{"cash", "cash", true},
		{"idontexist", "idontexist", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if valid := IsValidPayment(ctx, tt.value); valid != tt.valid {
				t.Fatalf(`valid = %v, want %v`, valid, tt.valid)
			}
		})
	}
}

func TestIsValidate(t *testing.T) {
	ctx := tests.Context()

	var tests = []struct {
		name  string
		field string
		value string
		err   error
	}{
		{"delivery=idontexist", "Delivery", "idontexist", errors.New("you are not authorized to process this request")},
		{"payment=idontexist", "Payment", "idontexist", errors.New("you are not authorized to process this request")},
		{"payment=idontexist", "Payment", "idontexist", errors.New("you are not authorized to process this request")},
		{"quantities=", "Quantities", "", errors.New("the cart is empty")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := order

			if tt.field == "Quantities" {
				o.Quantities = map[string]int{}
			} else if tt.field != "" {
				reflect.ValueOf(&o).Elem().FieldByName(tt.field).SetString(tt.value)
			}

			if err := o.Validate(ctx); fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.err) {
				t.Fatalf(`err = %v, want %s`, err, tt.err)
			}
		})
	}
}

func TestUpdateStatus(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/orders.redis")

	var tests = []struct {
		name  string
		id    string
		value string
		err   error
	}{
		{"status=processing", order.ID, "processing", nil},
		{"status=", "", order.ID, errors.New("input:status")},
		{"status=idontexist", order.ID, "idontexist", errors.New("input:status")},
		{"id=idontexist", "idontexist", "processing", errors.New("oops the data is not found")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateStatus(ctx, tt.id, tt.value); fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.err) {
				t.Fatalf(`err = %v, want %v`, err, tt.err)
			}
		})
	}
}

func TestFind(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/orders.redis")

	var tests = []struct {
		name string
		id   string
		err  error
	}{
		{"id=ORD1", order.ID, nil},
		{"id=", "", errors.New("oops the data is not found")},
		{"id=idontexist", "idontexist", errors.New("oops the data is not found")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Find(ctx, tt.id); fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.err) {
				t.Fatalf(`err = %v, want %v`, err, tt.err)
			}
		})
	}
}

func TestAddNote(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/orders.redis")

	var tests = []struct {
		name  string
		id    string
		value string
		err   error
	}{
		{"id=", "", "Useless", errors.New("oops the data is not found")},
		{"id=idontexist", "idontexist", "Useless", errors.New("oops the data is not found")},
		{"id=ORD1", "ORD1", "Useless", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddNote(ctx, tt.id, tt.value); fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.err) {
				t.Fatalf(`err = %v, want %v`, err, tt.err)
			}
		})
	}
}

func TestSendConfirmationEmail(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/orders.redis")

	message.SetString(language.English, "email_order_confirmation", "Hi %s,\nWoo hoo! Your order is on its way. Your order details can be found below.\n\n")
	message.SetString(language.English, "email_order_confirmationid", "Order ID: %s\n")
	message.SetString(language.English, "email_order_confirmationdate", "Order date: %s\n")
	message.SetString(language.English, "email_order_confirmationtotal", "Order total: %.2f\n\n")
	message.SetString(language.English, "title", "Title")
	message.SetString(language.English, "quality", "Quantity")
	message.SetString(language.English, "price", "Price")
	message.SetString(language.English, "total", "Total")
	message.SetString(language.English, "link", "Link")
	message.SetString(language.English, "email_order_confirmationfooter", "\nSee you around,\nThe Customer Experience Team at artisons shop")

	o, _ := order.WithProducts(ctx)
	o = o.WithTotal()

	tpl, err := o.SendConfirmationEmail(ctx)
	if err != nil {
		t.Fatalf(`err %v, want  nil`, err)
	}

	expected, err := os.ReadFile(cur + "testdata/mail.txt")
	if err != nil {
		t.Fatalf(`err %v, want  nil`, err)
	}

	got := strings.Join(strings.Fields(string(tpl)), "")
	exp := strings.Join(strings.Fields(string(expected)), "")

	if got != exp {
		t.Fatalf(`email = \n%s, want \n%s`, tpl, expected)
	}
}

func TestTotal(t *testing.T) {
	var tests = []struct {
		name     string
		products []products.Product
		total    float64
	}{
		{"total=65", []products.Product{{Quantity: 1, Price: 11}, {Quantity: 2, Price: 24.5}}, 65},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := order
			o.Products = tt.products
			o = o.WithTotal()

			if o.Total != tt.total {
				t.Fatalf(`total = %f, want %f`, o.Total, tt.total)
			}
		})
	}
}

func TestSearchReturnsOrdersWhenStatusIsFound(t *testing.T) {
	ctx := tests.Context()

	tests.Del(ctx, "order")
	tests.ImportData(ctx, cur+"testdata/orders.redis")

	var tests = []struct {
		name  string
		q     Query
		count int
	}{
		{"keywords=created", Query{Keywords: "created"}, 1},
		{"keywords=card", Query{Keywords: "card"}, 1},
		{"keywords=home", Query{Keywords: "home"}, 1},
		{"uid=1", Query{UID: 1}, 1},
		{"keywords=idontexist", Query{Keywords: "idontexist"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := Search(ctx, tt.q, 0, conf.DashboardMostItems)

			if err != nil {
				t.Fatalf(`err = %v, want nil`, err.Error())
			}

			if p.Total != tt.count {
				t.Fatalf(`total = %d, want %d`, p.Total, tt.count)
			}
		})
	}

}
