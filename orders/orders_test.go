package orders

import (
	"gifthub/string/stringutil"
	"gifthub/tests"
	"testing"
)

var order Order = Order{
	ID:       "test",
	UID:      1,
	Delivery: "collect",
	Payment:  "cash",
	Products: map[string]int64{
		"test": 1,
	},
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
	o.Products = map[string]int64{}
	ctx := tests.Context()

	if oid, err := o.Save(ctx); oid != "" || err == nil || err.Error() != "cart_empty" {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want string, 'cart_empty'`, oid, err)
	}
}

func TestSaveReturnsErrorWhenProductsAreUnavailable(t *testing.T) {
	o := order
	o.Products = map[string]int64{"toto12": 1}
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
