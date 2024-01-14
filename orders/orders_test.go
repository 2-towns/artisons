package orders

import (
	"gifthub/conf"
	"gifthub/products"
	"gifthub/string/stringutil"
	"gifthub/tests"
	"gifthub/users"
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

func TestIsValidDeliveryTrueWhenValid(t *testing.T) {
	ctx := tests.Context()
	if valid := IsValidDelivery(ctx, "collect"); !valid {
		t.Fatalf(`IsValidDelivery(ctx, "collect") = %v want true`, valid)
	}
}

func TestIsValidDeliveryFalseWhenInvalid(t *testing.T) {
	ctx := tests.Context()
	if valid := IsValidDelivery(ctx, "toto"); valid {
		t.Fatalf(`IsValidDelivery(ctx, "toto") = %v want false`, valid)
	}
}

func TestIsValidPaymentTrueWhenValid(t *testing.T) {
	ctx := tests.Context()
	if valid := IsValidPayment(ctx, "cash"); !valid {
		t.Fatalf(`IsValidPayment(ctx, "cash") = %v want true`, valid)
	}
}

func TestIsValidPaymentFalseWhenInvalid(t *testing.T) {
	ctx := tests.Context()
	if valid := IsValidPayment(ctx, "toto"); valid {
		t.Fatalf(`IsValidDelivery(ctx, "toto") = %v want false`, valid)
	}
}

func TestSaveReturnsNilWhenSuccess(t *testing.T) {
	o := order
	ctx := tests.Context()

	if oid, err := o.Save(ctx); err != nil || oid == "" {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want string, nil`, oid, err)
	}
}

func TestSaveReturnsErrorWhenDeliveryIsInvalid(t *testing.T) {
	o := order
	o.Delivery = "toto"
	ctx := tests.Context()

	if oid, err := o.Save(ctx); oid != "" || err == nil || err.Error() != "your are not authorized to process this request" {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want string, 'unauthorized'`, oid, err)
	}
}

func TestSaveReturnsErrorWhenPaymentIsInvalid(t *testing.T) {
	o := order
	o.Payment = "toto"
	ctx := tests.Context()

	if oid, err := o.Save(ctx); oid != "" || err == nil || err.Error() != "your are not authorized to process this request" {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want string, 'unauthorized'`, oid, err)
	}
}

func TestSaveReturnsErrorWhenProductsIsEmpty(t *testing.T) {
	o := order
	o.Quantities = map[string]int{}
	ctx := tests.Context()

	if oid, err := o.Save(ctx); oid != "" || err == nil || err.Error() != "The cart is empty." {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want string, 'The cart is empty.'`, oid, err)
	}
}

func TestSaveReturnsErrorWhenProductsAreUnavailable(t *testing.T) {
	o := order
	o.Quantities = map[string]int{"toto12": 1}
	ctx := tests.Context()

	if oid, err := o.Save(ctx); oid != "" || err == nil || err.Error() != "The cart is empty." {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want "", 'The cart is empty.'`, oid, err)
	}
}

func TestUpdateStatusReturnsNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	o := order
	if err := UpdateStatus(ctx, o.ID, "processing"); err != nil {
		t.Fatalf(`UpdateStatus(ctx, o.ID, "processing") = %v, nil`, err)
	}
}

func TestUpdateStatusReturnsErrorWhenStatusIsEmpty(t *testing.T) {
	ctx := tests.Context()
	o := order
	if err := UpdateStatus(ctx, o.ID, ""); err == nil || err.Error() != "input:status" {
		t.Fatalf(`UpdateStatus(ctx, o.ID, "") = '%v', want 'input:status'`, err)
	}
}

func TestUpdateStatusReturnsErrorWhenStatusIsInvalid(t *testing.T) {
	ctx := tests.Context()
	o := order
	if err := UpdateStatus(ctx, o.ID, "toto"); err == nil || err.Error() != "input:status" {
		t.Fatalf(`UpdateStatus(ctx, o.ID, "toto") = '%s', want 'input:status'`, err.Error())
	}
}

func TestUpdateStatusReturnsErrorWhenOrderDoesNotExist(t *testing.T) {
	oid, _ := stringutil.Random()
	ctx := tests.Context()

	if err := UpdateStatus(ctx, oid, "processing"); err == nil || err.Error() != "oops the data is not found" {
		t.Fatalf(`UpdateStatus(ctx, oid, "processing") = '%s', want 'oops the data is not found'`, err.Error())
	}
}

func TestFindReturnsOrderWhenSuccess(t *testing.T) {
	o := order
	ctx := tests.Context()

	if oo, err := Find(ctx, o.ID); err != nil || oo.ID == "" {
		t.Fatalf(`Find(ctx, o.ID) = %v, %v, nil`, oo, err)
	}
}

func TestFindReturnsErrorWhenOrderIsNotFound(t *testing.T) {
	oid, _ := stringutil.Random()
	ctx := tests.Context()

	if oo, err := Find(ctx, oid); err == nil || err.Error() != "oops the data is not found" {
		t.Fatalf(`Find(ctx, oid) = %v, %s, want Order{},'oops the data is not found'`, oo, err.Error())
	}
}

func TestAddNoteReturnsNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()

	if err := AddNote(ctx, "ORD1", "Ta commande, tu peux te la garder !"); err != nil {
		t.Fatalf(`AddNote(ctx, "ORD1", "Ta commande, tu peux te la garder !") = '%v', want nil`, err)
	}
}

func TestAddNoteReturnsErrorWhenNoteIsEmpty(t *testing.T) {
	ctx := tests.Context()

	if err := AddNote(ctx, "ORD1", ""); err == nil || err.Error() != "input:note" {
		t.Fatalf(`orders.AddNote(ctx, "ORD1", "") = '%s', want 'input:note'`, err.Error())
	}
}

func TestAddNoteReturnsErrorWhenOrderDoesNotExist(t *testing.T) {
	ctx := tests.Context()

	if err := AddNote(ctx, "123", "Useless"); err == nil || err.Error() != "oops the data is not found" {
		t.Fatalf(`orders.AddNote(ctx, "123", "Useless") = '%s', want 'oops the data is not found'`, err.Error())
	}
}

func TestProductsReturnsProductsWhenSuccess(t *testing.T) {
	ctx := tests.Context()

	o, err := order.WithProducts(ctx)
	if err != nil {
		t.Fatalf(`order.Products(ctx) = %v, %s, want []Products{}, nil`, o, err.Error())
	}

	if len(o.Products) == 0 {
		t.Fatalf(`len(p) = %d, want > 0`, len(o.Products))
	}

	p := o.Products[0]

	if p.Quantity == 0 {
		t.Fatalf(`op.Quantity = %d, want > 0`, p.Quantity)
	}

	if p.Slug == "" {
		t.Fatalf(`op.Slug = %s, want not empty`, p.Slug)
	}
}

func TestSendConfirmationEmailReturnsEmailContentWhenSuccess(t *testing.T) {
	ctx := tests.Context()

	message.SetString(language.English, "email_order_confirmation", "Hi %s,\nWoo hoo! Your order is on its way. Your order details can be found below.\n\n")
	message.SetString(language.English, "email_order_confirmationid", "Order ID: %s\n")
	message.SetString(language.English, "email_order_confirmationdate", "Order date: %s\n")
	message.SetString(language.English, "email_order_confirmationtotal", "Order total: %.2f\n\n")
	message.SetString(language.English, "title", "Title")
	message.SetString(language.English, "quality", "Quantity")
	message.SetString(language.English, "price", "Price")
	message.SetString(language.English, "total", "Total")
	message.SetString(language.English, "link", "Link")
	message.SetString(language.English, "email_order_confirmationfooter", "\nSee you around,\nThe Customer Experience Team at gifthub shop")

	o, _ := order.WithProducts(ctx)
	o = o.WithTotal()

	tpl, err := o.SendConfirmationEmail(ctx)
	if err != nil {
		t.Fatalf(`order.SendConfirmationEmail(ctx) = '%s', %v, want not empty, nil`, tpl, err)
	}

	expected := `Hi Arnaud,
Woo hoo! Your order is on its way. Your order details can be found below.

Order ID: ORD1
Order date: Friday, November 11
Order total: 105.50

+-----------------------------+----------+-------+-------+--------------------------------------------------------+
| TITLE                       | QUANTITY | PRICE | TOTAL | LINK                                                   |
+-----------------------------+----------+-------+-------+--------------------------------------------------------+
| T-shirt Tester c’est douter |        1 | 100.5 | 100.5 | http://localhost/PDT1-t-shirt-tester-c’est-douter.html |
+-----------------------------+----------+-------+-------+--------------------------------------------------------+

See you around,
The Customer Experience Team at gifthub shop`

	if tpl != expected {
		t.Fatalf(`tpl = \n%s, want \n%s`, tpl, expected)
	}
}

func TestTotalReturnsTheOrderTotalWhenSuccess(t *testing.T) {
	o := order
	o.Products = []products.Product{{Quantity: 1, Price: 11}, {Quantity: 2, Price: 24.5}}
	o = o.WithTotal()
	if o.Total != 65 {
		t.Fatalf(`order.Total(p)  = %f, want 65`, o.Total)
	}
}

func TestSearchReturnsOrdersWhenStatusIsFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keyword: "created"}, 0, conf.DashboardMostItems)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keyword: "created"}, 0, conf.DashboardMostItems) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if p.Orders[0].ID == "" {
		t.Fatalf(`p.Orders[0].ID = %s, want not empty`, p.Orders[0].ID)
	}
}

func TestSearchReturnsOrdersWhenPaymentIsFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keyword: "card"}, 0, conf.DashboardMostItems)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keyword: "card"}, 0, conf.DashboardMostItems) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if p.Orders[0].ID == "" {
		t.Fatalf(`p.Orders[0].ID = %s, want not empty`, p.Orders[0].ID)
	}
}

func TestSearchReturnsOrdersWhenDeliveryIsFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keyword: "home"}, 0, conf.ItemsPerPage)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keyword: "home"}, conf.ItemsPerPage)) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if p.Orders[0].ID == "" {
		t.Fatalf(`p.Orders[0].ID = %s, want not empty`, p.Orders[0].ID)
	}
}

func TestSearchReturnsNoOrdersWhenDeliveryIsCrazy(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keyword: "crazy"}, 0, conf.ItemsPerPage)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keyword: "crazy"}, 0, conf.ItemsPerPage) = %v, want nil`, err.Error())
	}

	if p.Total != 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}
}
