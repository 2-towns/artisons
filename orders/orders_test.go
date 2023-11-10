package orders

import (
	"gifthub/products"
	"gifthub/string/stringutil"
	"gifthub/tests"
	"gifthub/users"
	"log"
	"testing"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var order Order = Order{
	ID:       "test",
	UID:      1,
	Delivery: "collect",
	Payment:  "cash",
	Quantities: map[string]int{
		"test": 1,
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
	CreatedAt: time.Now(),
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

	if oid, err := o.Save(ctx); oid != "" || err == nil || err.Error() != "unauthorized" {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want string, 'unauthorized'`, oid, err)
	}
}

func TestSaveReturnsErrorWhenPaymentIsInvalid(t *testing.T) {
	o := order
	o.Payment = "toto"
	ctx := tests.Context()

	if oid, err := o.Save(ctx); oid != "" || err == nil || err.Error() != "unauthorized" {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want string, 'unauthorized'`, oid, err)
	}
}

func TestSaveReturnsErrorWhenProductsIsEmpty(t *testing.T) {
	o := order
	o.Quantities = map[string]int{}
	ctx := tests.Context()

	if oid, err := o.Save(ctx); oid != "" || err == nil || err.Error() != "cart_empty" {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want string, 'cart_empty'`, oid, err)
	}
}

func TestSaveReturnsErrorWhenProductsAreUnavailable(t *testing.T) {
	o := order
	o.Quantities = map[string]int{"toto12": 1}
	ctx := tests.Context()

	if oid, err := o.Save(ctx); oid != "" || err == nil || err.Error() != "cart_empty" {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want "", 'cart_empty'`, oid, err)
	}
}

func TestUpdateStatusReturnsNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	o := order
	if err := UpdateStatus(ctx, o.ID, "processing"); err != nil {
		t.Fatalf(`UpdateStatus(ctx, o.ID, "processing") = %v, nil`, err)
	}
}

func TestUpdateStatusReturnsErrorWhenStatusIsInvalid(t *testing.T) {
	ctx := tests.Context()
	o := order
	if err := UpdateStatus(ctx, o.ID, "toto"); err == nil || err.Error() != "order_bad_status" {
		t.Fatalf(`UpdateStatus(ctx, o.ID, "toto") = '%s', want 'order_bad_status'`, err.Error())
	}
}

func TestUpdateStatusReturnsErrorWhenOrderDoesNotExist(t *testing.T) {
	oid, _ := stringutil.Random()
	ctx := tests.Context()

	if err := UpdateStatus(ctx, oid, "processing"); err == nil || err.Error() != "order_not_found" {
		t.Fatalf(`UpdateStatus(ctx, oid, "processing") = '%s', want 'order_not_found'`, err.Error())
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

	if oo, err := Find(ctx, oid); err == nil || err.Error() != "order_not_found" {
		t.Fatalf(`Find(ctx, oid) = %v, %s, want Order{},'order_not_found'`, oo, err.Error())
	}
}

func TestAddNoteReturnsNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()

	if err := AddNote(ctx, "test", "Ta commande, tu peux te la garder !"); err != nil {
		t.Fatalf(`AddNote(ctx, "test", "Ta commande, tu peux te la garder !") = '%v', want nil`, err)
	}
}

func TestAddNoteReturnsErrorWhenNoteIsEmpty(t *testing.T) {
	ctx := tests.Context()

	if err := AddNote(ctx, "test", ""); err == nil || err.Error() != "order_note_required" {
		t.Fatalf(`orders.AddNote(ctx, "test", "") = '%s', want 'order_note_required'`, err.Error())
	}
}

func TestAddNoteReturnsErrorWhenOrderDoesNotExist(t *testing.T) {
	ctx := tests.Context()

	if err := AddNote(ctx, "123", "Useless"); err == nil || err.Error() != "order_not_found" {
		t.Fatalf(`orders.AddNote(ctx, "123", "Useless") = '%s', want 'order_not_found'`, err.Error())
	}
}

func TestProductsReturnsProductsWhenSuccess(t *testing.T) {
	ctx := tests.Context()

	p, err := order.Products(ctx)
	if err != nil {
		t.Fatalf(`order.Products(ctx) = %v, %s, want []Products{}, nil`, p, err.Error())
	}

	if len(p) == 0 {
		t.Fatalf(`len(p) = %d, want > 0`, len(p))
	}

	op := p[0]

	if op.Quantity == 0 {
		t.Fatalf(`op.Quantity = %d, want > 0`, op.Quantity)
	}

	if op.Slug == "" {
		t.Fatalf(`op.Slug = %s, want not empty`, op.Slug)
	}
}

func TestSendConfirmationEmailReturnsEmailContentWhenSuccess(t *testing.T) {
	ctx := tests.Context()

	message.SetString(language.English, "order_created_email", "Hi %s,\nWoo hoo! Your order is on its way. Your order details can be found below.\n\n")
	message.SetString(language.English, "order_id_email", "Order ID: %s\n")
	message.SetString(language.English, "order_date_email", "Order date: %s\n")
	message.SetString(language.English, "order_total_email", "Order total: %.2f\n\n")
	message.SetString(language.English, "order_title", "Title")
	message.SetString(language.English, "order_quantity", "Quantity")
	message.SetString(language.English, "order_price", "Price")
	message.SetString(language.English, "order_total", "Total")
	message.SetString(language.English, "order_link", "Link")
	message.SetString(language.English, "order_footer_email", "\nSee you around,\nThe Customer Experience Team at gifthub shop")

	tpl, err := order.SendConfirmationEmail(ctx)
	if err != nil {
		t.Fatalf(`order.SendConfirmationEmail(ctx) = '%s', %v, want not empty, nil`, tpl, err)
	}

	expected := `Hi Prince Benton Bins,
Woo hoo! Your order is on its way. Your order details can be found below.

Order ID: test
Order date: Friday, November 11
Order total: 105.50

+-----------------------------------+----------+-------+-------+---------------------------------------------------------+
| TITLE                             | QUANTITY | PRICE | TOTAL | LINK                                                    |
+-----------------------------------+----------+-------+-------+---------------------------------------------------------+
| Un joli pull tricoté par ma maman |        1 | 100.5 | 100.5 | http://localhost/test-un-joli-pull-tricoté-par-ma-maman |
+-----------------------------------+----------+-------+-------+---------------------------------------------------------+

See you around,
The Customer Experience Team at gifthub shop`
	log.Println(tpl)
	log.Println(expected)
	if tpl != expected {
		t.Fatalf(`tpl != expected`)
	}
}

func TestTotalReturnsTheOrderTotalWhenSuccess(t *testing.T) {
	p := []products.Product{{Quantity: 1, Price: 11}, {Quantity: 2, Price: 24.5}}
	total := order.Total(p)
	if total != 65 {
		t.Fatalf(`order.Total(p)  = %f, want 65`, total)
	}
}
