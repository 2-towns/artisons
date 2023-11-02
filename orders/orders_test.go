package orders

import (
	"context"
	"gifthub/db"
	"gifthub/string/stringutil"
	"gifthub/tests"
	"math/rand"
	"testing"
)

func createOrder() Order {
	ctx := context.Background()
	uid := rand.Int63n(10000)
	pid, _ := stringutil.Random()
	oid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "product:"+pid, "status", "online")
	db.Redis.HSet(ctx, "order:"+oid, "uid", uid)
	db.Redis.HSet(ctx, "order:"+oid, "id", oid)
	products := map[string]int64{
		pid: 1,
	}

	o := Order{
		ID:       oid,
		UID:      uid,
		Delivery: "collect",
		Payment:  "cash",
		Products: products,
	}

	return o
}

// TestIsValidDelivery expects to succeed
func TestIsValidDelivery(t *testing.T) {
	ctx := tests.Context()
	if valid := IsValidDelivery(ctx, "collect"); !valid {
		t.Fatalf(`IsValidDelivery(ctx, "collect") = %v want true`, valid)
	}
}

// TestIsValidDeliveryMisvalue expects to fail because of delivery invalid
func TestIsValidDeliveryMisvalue(t *testing.T) {
	ctx := tests.Context()
	if valid := IsValidDelivery(ctx, "toto"); valid {
		t.Fatalf(`IsValidDelivery(ctx, "toto") = %v want false`, valid)
	}
}

// TestIsValidPayment expects to succeed
func TestIsValidPayment(t *testing.T) {
	ctx := tests.Context()
	if valid := IsValidPayment(ctx, "cash"); !valid {
		t.Fatalf(`IsValidPayment(ctx, "cash") = %v want true`, valid)
	}
}

// TestIsValidPaymentMisvalue expects to fail because of payment invalid
func TestIsValidPaymentMisvalue(t *testing.T) {
	ctx := tests.Context()
	if valid := IsValidPayment(ctx, "toto"); valid {
		t.Fatalf(`IsValidDelivery(ctx, "toto") = %v want false`, valid)
	}
}

// TestOrderSave expects to succeed
func TestOrderSave(t *testing.T) {
	o := createOrder()
	ctx := tests.Context()

	if oid, err := o.Save(ctx); err != nil || oid == "" {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want string, nil`, oid, err)
	}
}

// TestOrderSaveDeliveryMisvalue expects to fail because of wrong delivery value
func TestOrderSaveDeliveryMisvalue(t *testing.T) {
	o := createOrder()
	o.Delivery = "toto"
	ctx := tests.Context()

	if oid, err := o.Save(ctx); oid != "" || err == nil || err.Error() != "unauthorized" {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want string, 'unauthorized'`, oid, err)
	}
}

// TestOrderSavePaymentMisvalue expects to fail because of wrong payment value
func TestOrderSavePaymentMisvalue(t *testing.T) {
	o := createOrder()
	o.Payment = "toto"
	ctx := tests.Context()

	if oid, err := o.Save(ctx); oid != "" || err == nil || err.Error() != "unauthorized" {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want string, 'unauthorized'`, oid, err)
	}
}

// TestOrderSaveProductsEmpty expects to fail because of products emptyness
func TestOrderSaveProductsEmpty(t *testing.T) {
	o := createOrder()
	o.Products = map[string]int64{}
	ctx := tests.Context()

	if oid, err := o.Save(ctx); oid != "" || err == nil || err.Error() != "cart_empty" {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want string, 'cart_empty'`, oid, err)
	}
}

// TestOrderSaveProductsUnavailable expects to fail because of products availability
func TestOrderSaveProductsUnavailable(t *testing.T) {
	o := createOrder()
	o.Products = map[string]int64{"toto12": 1}
	ctx := tests.Context()

	if oid, err := o.Save(ctx); oid != "" || err == nil || err.Error() != "cart_empty" {
		t.Fatalf(`o.Save(ctx) = '%s', %v, want "", 'cart_empty'`, oid, err)
	}
}

// TestOrderUpdateStatus expects to succeed
func TestOrderUpdateStatus(t *testing.T) {
	ctx := tests.Context()
	o := createOrder()
	if err := UpdateStatus(ctx, o.ID, "processing"); err != nil {
		t.Fatalf(`UpdateStatus(ctx, o.ID, "processing") = %v, nil`, err)
	}
}

// TestOrderUpdateStatusWrong expects to fail because of invalid status
func TestOrderUpdateStatusWrong(t *testing.T) {
	ctx := tests.Context()
	o := createOrder()
	if err := UpdateStatus(ctx, o.ID, "toto"); err == nil || err.Error() != "order_bad_status" {
		t.Fatalf(`UpdateStatus(ctx, o.ID, "toto") = '%s', want 'order_bad_status'`, err.Error())
	}
}

// TestOrderUpdateStatusNotExisting expects to fail because of not existing order
func TestOrderUpdateStatusNotExisting(t *testing.T) {
	oid, _ := stringutil.Random()
	ctx := tests.Context()

	if err := UpdateStatus(ctx, oid, "processing"); err == nil || err.Error() != "order_not_found" {
		t.Fatalf(`UpdateStatus(ctx, oid, "processing") = '%s', want 'order_not_found'`, err.Error())
	}
}

// TestOrderFind expects to succeed
func TestOrderFind(t *testing.T) {
	o := createOrder()
	ctx := tests.Context()

	if oo, err := Find(ctx, o.ID); err != nil || oo.ID == "" {
		t.Fatalf(`Find(ctx, o.ID) = %v, %v, nil`, oo, err)
	}
}

// TestOrderUpdateStatusNotExisting expects to fail because of not existing order
func TestOrderFindNotExisting(t *testing.T) {
	oid, _ := stringutil.Random()
	ctx := tests.Context()

	if oo, err := Find(ctx, oid); err == nil || err.Error() != "order_not_found" {
		t.Fatalf(`Find(ctx, oid) = %v, %s, want Order{},'order_not_found'`, oo, err.Error())
	}
}

// TestAddNote expects to succeed
func TestAddNote(t *testing.T) {
	o := createOrder()
	ctx := tests.Context()

	if err := AddNote(ctx, o.ID, "Ta commande, tu peux te la garder !"); err != nil {
		t.Fatalf(`AddNote(ctx,o.ID, "Ta commande, tu peux te la garder !") = '%v', want nil`, err)
	}
}

// TestAddNote expects to fail because of note emptyness
func TestAddNoteWithEmptyString(t *testing.T) {
	o := createOrder()
	ctx := tests.Context()

	if err := AddNote(ctx, o.ID, ""); err == nil || err.Error() != "order_note_required" {
		t.Fatalf(`orders.AddNote(ctx,o.ID, "") = '%s', want 'order_note_required'`, err.Error())
	}
}

// TestAddNote expects to fail because of order not found
func TestAddNoteWithOrderNotExisting(t *testing.T) {
	ctx := tests.Context()

	if err := AddNote(ctx, "123", "Useless"); err == nil || err.Error() != "order_not_found" {
		t.Fatalf(`orders.AddNote(ctx, o.ID, "Useless") = '%s', want 'order_not_found'`, err.Error())
	}
}
