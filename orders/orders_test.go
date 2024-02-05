package orders

import (
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/products"
	"artisons/string/stringutil"
	"artisons/tests"
	"artisons/users"
	"context"
	"fmt"
	"testing"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var order Order = Order{
	ID:       tests.OrderID,
	UID:      tests.UserID1,
	Delivery: "collect",
	Payment:  "cash",
	Total:    tests.OrderTotal,
	Quantities: map[string]int{
		tests.ProductID1: 1,
	},
	DeliveryCost: 5,
	Address: users.Address{
		Firstname:     tests.OrderFirstName,
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
		t.Fatalf(`valid = %v want true`, valid)
	}
}

func TestIsValidDeliveryFalseWhenInvalid(t *testing.T) {
	ctx := tests.Context()
	if valid := IsValidDelivery(ctx, tests.DoesNotExist); valid {
		t.Fatalf(`IsValidDelivery(ctx, tests.DoesNotExist) = %v want false`, valid)
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

func TestValidateReturnsErrorWhenDeliveryIsInvalid(t *testing.T) {
	o := order
	o.Delivery = "toto"
	ctx := tests.Context()

	if err := o.Validate(ctx); err == nil || err.Error() != "you are not authorized to process this request" {
		t.Fatalf(`o.Validate(ctx) = %v, want string, 'unauthorized'`, err)
	}
}

func TestValidateReturnsErrorWhenPaymentIsInvalid(t *testing.T) {
	o := order
	o.Payment = tests.DoesNotExist
	ctx := tests.Context()

	if err := o.Validate(ctx); err == nil || err.Error() != "you are not authorized to process this request" {
		t.Fatalf(`o.Save(ctx) =  %v, want string, 'unauthorized'`, err)
	}
}

func TestValidateReturnsErrorWhenProductsIsEmpty(t *testing.T) {
	o := order
	o.Quantities = map[string]int{}
	ctx := tests.Context()

	if err := o.Validate(ctx); err == nil || err.Error() != "the cart is empty" {
		t.Fatalf(`o.Validate(ctx) = %v, want string, 'the cart is empty'`, err)
	}
}

func TestSaveReturnsErrorWhenProductsAreUnavailable(t *testing.T) {
	o := order
	o.Quantities = map[string]int{tests.DoesNotExist: 1}
	ctx := tests.Context()

	if oid, err := o.Save(ctx); oid != "" || err == nil || err.Error() != "the cart is empty" {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want "", 'the cart is empty'`, oid, err)
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
	if err := UpdateStatus(ctx, o.ID, tests.DoesNotExist); err == nil || err.Error() != "input:status" {
		t.Fatalf(`UpdateStatus(ctx, o.ID, tests.DoesNotExist) = '%s', want 'input:status'`, err.Error())
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

	if err := AddNote(ctx, order.ID, "Ta commande, tu peux te la garder !"); err != nil {
		t.Fatalf(`AddNote(ctx, order.ID, "Ta commande, tu peux te la garder !") = '%v', want nil`, err)
	}
}

func TestAddNoteReturnsErrorWhenNoteIsEmpty(t *testing.T) {
	ctx := tests.Context()

	if err := AddNote(ctx, order.ID, ""); err == nil || err.Error() != "input:note" {
		t.Fatalf(`orders.AddNote(ctx, order.ID, "") = '%s', want 'input:note'`, err.Error())
	}
}

func TestAddNoteReturnsErrorWhenOrderDoesNotExist(t *testing.T) {
	ctx := tests.Context()

	if err := AddNote(ctx, tests.DoesNotExist, "Useless"); err == nil || err.Error() != "oops the data is not found" {
		t.Fatalf(`orders.AddNote(ctx, tests.DoesNotExist, "Useless") = '%s', want 'oops the data is not found'`, err.Error())
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
	//
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
		t.Fatalf(`order.SendConfirmationEmail(ctx) = '%s', %v, want not empty, nil`, tpl, err)
	}

	expected := `Hi Arnaud,
Woo hoo! Your order is on its way. Your order details can be found below.

Order ID: ` + order.ID + `
Order date: Friday, November 11
Order total: ` + fmt.Sprintf("%.2f", order.Total) + `

+-----------------------------+----------+-------+-------+--------------------------------------------------------+
| TITLE                       | QUANTITY | PRICE | TOTAL | LINK                                                   |
+-----------------------------+----------+-------+-------+--------------------------------------------------------+
| T-shirt Tester c'est douter |        1 | 100.5 | 100.5 | http://localhost/PDT1-t-shirt-tester-c-est-douter.html |
+-----------------------------+----------+-------+-------+--------------------------------------------------------+

See you around,
The Customer Experience Team at artisons shop`

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
	p, err := Search(c, Query{Keywords: "created"}, 0, conf.DashboardMostItems)
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
	p, err := Search(c, Query{Keywords: "card"}, 0, conf.DashboardMostItems)
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
	p, err := Search(c, Query{Keywords: "home"}, 0, conf.ItemsPerPage)
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

func TestSearchReturnsOrdersWhenUIDIsFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{UID: tests.UserID1}, 0, conf.ItemsPerPage)
	if err != nil {
		t.Fatalf(`Search(c, Query{UID:  tests.UserID1}, conf.ItemsPerPage)) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if p.Orders[0].ID == "" {
		t.Fatalf(`p.Orders[0].ID = %s, want not empty`, p.Orders[0].ID)
	}
}

func TestSearchReturnsNoOrdersWhenDeliveryIsInvalid(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keywords: tests.DoesNotExist}, 0, conf.ItemsPerPage)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keyword: "crazy"}, 0, conf.ItemsPerPage) = %v, want nil`, err.Error())
	}

	if p.Total != 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}
}

func TestSearchReturnUpdatedAtSortedOrdersWhenEndIsBack(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{}, 0, 999)
	if err != nil {
		t.Fatalf(`Search(c, Query{}, 0, 999) = %v, want nil`, err.Error())
	}

	if p.Total <= 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	a := 0
	b := 0

	for idx, val := range p.Orders {
		if val.ID == tests.OrderID {
			a = idx
		}

		if val.ID == tests.OrderID2 {
			b = idx
		}
	}

	if a > b {
		t.Fatal(`a < b, want a > b`)
	}
}

func TestSearchReturnUpdatedAtSortedOrdersWhenEndIsFront(t *testing.T) {
	c := tests.Context()
	c = context.WithValue(c, contexts.Domain, "front")
	p, err := Search(c, Query{}, 0, 999)
	if err != nil {
		t.Fatalf(`Search(c, Query{}, 0, 999) = %v, want nil`, err.Error())
	}

	if p.Total <= 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	a := 0
	b := 0

	for idx, val := range p.Orders {
		if val.ID == tests.OrderID {
			a = idx
		}

		if val.ID == tests.OrderID2 {
			b = idx
		}
	}

	if a < b {
		t.Fatal(`a < b, want a > b`)
	}
}
