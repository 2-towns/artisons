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
	db.Redis.HSet(ctx, "product:"+pid, "status", "online")
	products := map[string]int64{
		pid: 1,
	}

	o := Order{
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
		t.Fatalf(`o.Save() = %s, %v, want string, nil`, oid, err)
	}
}

// TestOrderSaveDeliveryMisvalue expects to fail because of wrong delivery value
func TestOrderSaveDeliveryMisvalue(t *testing.T) {
	o := createOrder()
	o.Delivery = "toto"
	oid, err := o.Save()
	if oid != "" || err == nil || err.Error() != "unauthorized" {
		t.Fatalf(`o.Save() = %s, %v, want string, 'unauthorized'`, oid, err)
	}
}

// TestOrderSavePaymentMisvalue expects to fail because of wrong payment value
func TestOrderSavePaymentMisvalue(t *testing.T) {
	o := createOrder()
	o.Payment = "toto"
	oid, err := o.Save()
	if oid != "" || err == nil || err.Error() != "unauthorized" {
		t.Fatalf(`o.Save() = %s, %v, want string, 'unauthorized'`, oid, err)
	}
}

// TestOrderSaveProductsEmpty expects to fail because of products emptyness
func TestOrderSaveProductsEmpty(t *testing.T) {
	o := createOrder()
	o.Products = map[string]int64{}
	oid, err := o.Save()
	if oid != "" || err == nil || err.Error() != "cart_empty" {
		t.Fatalf(`o.Save() = %s, %v, want string, 'cart_empty'`, oid, err)
	}
}

// TestOrderSaveProductsUnavailable expects to fail because of products availability
func TestOrderSaveProductsUnavailable(t *testing.T) {
	o := createOrder()
	o.Products = map[string]int64{"toto12": 1}
	oid, err := o.Save()
	if oid != "" || err == nil || err.Error() != "cart_empty" {
		t.Fatalf(`o.Save() = %s, %v, want "", 'cart_empty'`, oid, err)
	}
}
