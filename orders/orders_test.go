package orders

import (
	"context"
	"gifthub/db"
	"gifthub/string/stringutil"
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
	if valid := IsValidDelivery("collect"); !valid {
		t.Fatalf(`IsValidDelivery("collect") = %v want true`, valid)
	}
}

// TestIsValidDeliveryMisvalue expects to fail because of delivery invalid
func TestIsValidDeliveryMisvalue(t *testing.T) {
	if valid := IsValidDelivery("toto"); valid {
		t.Fatalf(`IsValidDelivery("toto") = %v want false`, valid)
	}
}

// TestIsValidPayment expects to succeed
func TestIsValidPayment(t *testing.T) {
	if valid := IsValidPayment("cash"); !valid {
		t.Fatalf(`IsValidPayment("cash") = %v want true`, valid)
	}
}

// TestIsValidPaymentMisvalue expects to fail because of payment invalid
func TestIsValidPaymentMisvalue(t *testing.T) {
	if valid := IsValidPayment("toto"); valid {
		t.Fatalf(`IsValidDelivery("toto") = %v want false`, valid)
	}
}

// TestOrderSave expects to succeed
func TestOrderSave(t *testing.T) {
	o := createOrder()

	if oid, err := o.Save(); err != nil || oid == "" {
		t.Fatalf(`o.Save() = '%s', %v, want string, nil`, oid, err)
	}
}

// TestOrderSaveDeliveryMisvalue expects to fail because of wrong delivery value
func TestOrderSaveDeliveryMisvalue(t *testing.T) {
	o := createOrder()
	o.Delivery = "toto"

	if oid, err := o.Save(); oid != "" || err == nil || err.Error() != "unauthorized" {
		t.Fatalf(`o.Save() = '%s', %v, want string, 'unauthorized'`, oid, err)
	}
}

// TestOrderSavePaymentMisvalue expects to fail because of wrong payment value
func TestOrderSavePaymentMisvalue(t *testing.T) {
	o := createOrder()
	o.Payment = "toto"

	if oid, err := o.Save(); oid != "" || err == nil || err.Error() != "unauthorized" {
		t.Fatalf(`o.Save() = '%s', %v, want string, 'unauthorized'`, oid, err)
	}
}

// TestOrderSaveProductsEmpty expects to fail because of products emptyness
func TestOrderSaveProductsEmpty(t *testing.T) {
	o := createOrder()
	o.Products = map[string]int64{}

	if oid, err := o.Save(); oid != "" || err == nil || err.Error() != "cart_empty" {
		t.Fatalf(`o.Save() = '%s', %v, want string, 'cart_empty'`, oid, err)
	}
}

// TestOrderSaveProductsUnavailable expects to fail because of products availability
func TestOrderSaveProductsUnavailable(t *testing.T) {
	o := createOrder()
	o.Products = map[string]int64{"toto12": 1}

	if oid, err := o.Save(); oid != "" || err == nil || err.Error() != "cart_empty" {
		t.Fatalf(`o.Save() = '%s', %v, want "", 'cart_empty'`, oid, err)
	}
}

// TestOrderUpdateStatus expects to succeed
func TestOrderUpdateStatus(t *testing.T) {
	o := createOrder()
	if err := UpdateStatus(o.ID, "processing"); err != nil {
		t.Fatalf(`UpdateStatus(o.ID, "processing") = %v, %v, nil`, err)
	}
}

// TestOrderUpdateStatusWrong expects to fail because of invalid status
func TestOrderUpdateStatusWrong(t *testing.T) {
	o := createOrder()
	if err := UpdateStatus(o.ID, "toto"); err == nil || err.Error() != "unauthorized" || oo.ID != "" {
		t.Fatalf(`UpdateStatus(o.ID, "toto") = %v, %s, want 'unauthorized'`, err)
	}
}

// TestOrderUpdateStatusNotExisting expects to fail because of not existing order
func TestOrderUpdateStatusNotExisting(t *testing.T) {
	oid, _ := stringutil.Random()

	if err := UpdateStatus(oid, "processing"); err == nil || err.Error() != "order_not_found" {
		t.Fatalf(`UpdateStatus(oid, "processing") = %v, %s, want 'order_not_found'`, err)
	}
}

// TestOrderFind expects to succeed
func TestOrderFind(t *testing.T) {
	o := createOrder()

	if oo, err := Find(o.ID); err != nil || oo.ID == "" {
		t.Fatalf(`Find(o.ID) = %v, %v, nil`, oo, err)
	}
}

// TestOrderUpdateStatusNotExisting expects to fail because of not existing order
func TestOrderFindNotExisting(t *testing.T) {
	oid, _ := stringutil.Random()

	if oo, err := Find(oid); err == nil || err.Error() != "order_not_found" {
		t.Fatalf(`Find(oid) = %v, %s, want Order{},'order_not_found'`, oo, err.Error())
	}
}

// TestAddNote expects to succeed
func TestAddNote(t *testing.T) {
	o := createOrder()

	if err := AddNote(o.ID, "Ta commande, tu peux te la garder !"); err != nil {
		t.Fatalf(`AddNote(o.ID, "Ta commande, tu peux te la garder !") = '%v', want nil`, err)
	}
}

// TestAddNote expects to fail because of note emptyness
func TestAddNoteWithEmptyString(t *testing.T) {
	o := createOrder()

	if err := AddNote(o.ID, ""); err == nil || err.Error() != "order_note_required" {
		t.Fatalf(`orders.AddNote(o.ID, "") = '%s', want 'order_note_required'`, err.Error())
	}
}

// TestAddNote expects to fail because of order not found
func TestAddNoteWithOrderNotExisting(t *testing.T) {
	if err := AddNote("123", "Useless"); err == nil || err.Error() != "order_not_found" {
		t.Fatalf(`orders.AddNote(o.ID, "Useless") = '%s', want 'order_not_found'`, err.Error())
	}
}
