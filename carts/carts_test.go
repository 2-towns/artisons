package carts

import (
	"context"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/string/stringutil"
	"testing"
)

func createCart() Cart {
	ctx := context.Background()
	cid, _ := stringutil.Random()

	db.Redis.HSet(ctx, "cart:"+cid, "cid", cid)
	db.Redis.Expire(ctx, "cart:"+cid, conf.CartDuration)

	c := Cart{ID: cid}

	return c
}

func createProduct() string {
	ctx := context.Background()
	pid, _ := stringutil.Random()

	db.Redis.HSet(ctx, "product:"+pid, "status", "online")

	return pid
}

func addProductToCart(cid, pid string) {
	ctx := context.Background()
	db.Redis.HSet(ctx, "cart:"+cid, "product:"+pid, 1)

}

// TestCartAdd expects to succeed
func TestCartAdd(t *testing.T) {
	c := createCart()
	pid := createProduct()
	var quantity int64 = 1

	if err := Add(c.ID, pid, quantity); err != nil {
		t.Fatalf("Add(uid,pid,quantity), %v, want nil, error", err)
	}
}

// TestCartAddWithCIDMisvalue expects to fail because of cid misvalue
func TestCartAddWithCIDMisvalue(t *testing.T) {
	cid, _ := stringutil.Random()
	pid := createProduct()
	var quantity int64 = 1

	if err := Add(cid, pid, quantity); err == nil {
		t.Fatalf("Add(uid,pid,quantity), %v, not want nil, error", err)
	}
}

// TestCartAddWithPIDMisvalue expects to fail because of pid misvalue
func TestCartAddWithPIDMisvalue(t *testing.T) {
	c := createCart()
	pid, _ := stringutil.Random()
	var quantity int64 = 1

	if err := Add(c.ID, pid, quantity); err == nil {
		t.Fatalf("Add(uid,pid,quantity), %v, not want nil, error", err)
	}
}

// TestUdpateCID expects to succeed
func TestRefreshCID(t *testing.T) {
	cid, err := RefreshCID("")
	if cid == "" || err != nil {
		t.Fatalf("UpdateCID('') = %s, %v, not want ', error", cid, err)
	}
}

// TestRefreshCIDExisting expects to succeed with existing cartID
func TestRefreshCIDExisting(t *testing.T) {
	s, _ := stringutil.Random()
	cid, err := RefreshCID(s)
	if cid == "" || err != nil {
		t.Fatalf("UpdateCID(s) = %s, %v, want string, error", cid, err)
	}
}

// TestGetCart expects to succeed
func TestGetCart(t *testing.T) {
	c := createCart()
	pid := createProduct()
	addProductToCart(c.ID, pid)

	c, err := Get(c.ID)
	// TODO: Add the test len(c.Products) == 0
	if c.ID == "" || err != nil {
		t.Fatalf("Get(cid) = %v, %v, want Cart, nil", c, err)
	}
}

// TestGetCartWithCIDNotExisting expects to fail because of cid non existence
func TestGetCartWithCIDNotExisting(t *testing.T) {
	c, err := Get("toto")
	if c.ID != "" || len(c.Products) != 0 || err == nil || err.Error() != "cart_not_found" {
		t.Fatalf("Get(cid) = %v, %v, want Cart{}, 'cart_not_found'", c, err)
	}
}

// TestCartUpdateDelivery expects to succeed
func TestCartUpdateDelivery(t *testing.T) {
	c := createCart()
	if err := c.UpdateDelivery("collect"); err != nil {
		t.Fatalf("c.UpdateDelivery('collect') = %v, want nil", err)
	}
}

// TestCartUpdateDeliveryWithMisvalue expects to fail because of delivery misvalue
func TestCartUpdateDeliveryWithMisvalue(t *testing.T) {
	c := createCart()
	if err := c.UpdateDelivery("toto"); err == nil || err.Error() != "unauthorized" {
		t.Fatalf("c.UpdateDelivery('toto') = %v, want unauthorized", err)
	}
}

// TestCartUpdatePayment expects to succeed
func TestCartUpdatePayment(t *testing.T) {
	c := createCart()
	if err := c.UpdatePayment("card"); err != nil {
		t.Fatalf("c.UpdatePayment('card') = %v, want nil", err)
	}
}

// TestCartUpdatePaymentWithMisvalue  expects to fail because of payment misvalue
func TestCartUpdatePaymentWithMisvalue(t *testing.T) {
	c := createCart()
	if err := c.UpdatePayment("toto"); err == nil || err.Error() != "unauthorized" {
		t.Fatalf("c.UpdatePayment('toto') = %v, want 'unauthorized'", err)
	}
}
