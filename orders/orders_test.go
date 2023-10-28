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
	oid, err := o.Save()
	if err != nil || oid == "" {
		t.Fatalf(`o.Save() = '%s', %v, want string, nil`, oid, err)
	}
}

// TestOrderSaveDeliveryMisvalue expects to fail because of wrong delivery value
func TestOrderSaveDeliveryMisvalue(t *testing.T) {
	o := createOrder()
	o.Delivery = "toto"
	oid, err := o.Save()
	if oid != "" || err == nil || err.Error() != "unauthorized" {
		t.Fatalf(`o.Save() = '%s', %v, want string, 'unauthorized'`, oid, err)
	}
}

// TestOrderSavePaymentMisvalue expects to fail because of wrong payment value
func TestOrderSavePaymentMisvalue(t *testing.T) {
	o := createOrder()
	o.Payment = "toto"
	oid, err := o.Save()
	if oid != "" || err == nil || err.Error() != "unauthorized" {
		t.Fatalf(`o.Save() = '%s', %v, want string, 'unauthorized'`, oid, err)
	}
}

// TestOrderSaveProductsEmpty expects to fail because of products emptyness
func TestOrderSaveProductsEmpty(t *testing.T) {
	o := createOrder()
	o.Products = map[string]int64{}
	oid, err := o.Save()
	if oid != "" || err == nil || err.Error() != "cart_empty" {
		t.Fatalf(`o.Save() = '%s', %v, want string, 'cart_empty'`, oid, err)
	}
}

// TestOrderSaveProductsUnavailable expects to fail because of products availability
func TestOrderSaveProductsUnavailable(t *testing.T) {
	o := createOrder()
	o.Products = map[string]int64{"toto12": 1}
	oid, err := o.Save()
	if oid != "" || err == nil || err.Error() != "cart_empty" {
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
