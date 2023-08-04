package carts

import (
	"context"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/string/stringutil"
	"testing"

	"github.com/go-faker/faker/v4"
)

// TestCartAdd adds an product into user cart.
func TestCartAdd(t *testing.T) {
	ctx := context.Background()
	cid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "cart:"+cid, "cid", cid)
	db.Redis.Expire(ctx, "cart:"+cid, conf.CartDuration)

	pid, _ := stringutil.Random()
	title := faker.Sentence()
	db.Redis.HSet(ctx, "product:"+pid, "title", title)

	var quantity int64 = 1

	if err := Add(cid, pid, quantity); err != nil {
		t.Fatalf("Add(uid,pid,quantity), %v, want nil, error", err)
	}
}

// TestCartAddWithCIDMisvalue adding a product into a non existing
// cart returns an error.
func TestCartAddWithCIDMisvalue(t *testing.T) {
	ctx := context.Background()
	cid, _ := stringutil.Random()

	pid, _ := stringutil.Random()
	title := faker.Sentence()
	db.Redis.HSet(ctx, "product:"+pid, "title", title)

	var quantity int64 = 1

	if err := Add(cid, pid, quantity); err == nil {
		t.Fatalf("Add(uid,pid,quantity), %v, not want nil, error", err)
	}
}

// TestCartAddWithPIDMisvalue adding a non existing product into a
// cart returns an error.
func TestCartAddWithPIDMisvalue(t *testing.T) {
	ctx := context.Background()
	cid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "cart:"+cid, "cid", cid)
	db.Redis.Expire(ctx, "cart:"+cid, conf.CartDuration)

	pid, _ := stringutil.Random()
	var quantity int64 = 1

	if err := Add(cid, pid, quantity); err == nil {
		t.Fatalf("Add(uid,pid,quantity), %v, not want nil, error", err)
	}
}

// TestUdpateCID generates an cart ID.
func TestRefreshCID(t *testing.T) {
	cid, err := RefreshCID("")
	if cid == "" || err != nil {
		t.Fatalf("UpdateCID('') = %s, %v, not want ', error", cid, err)
	}
}

// TestRefreshCIDExisting refreshes the cart ID.
func TestRefreshCIDExisting(t *testing.T) {
	s, _ := stringutil.Random()
	cid, err := RefreshCID(s)
	if cid == "" || err != nil {
		t.Fatalf("UpdateCID(s) = %s, %v, want string, error", cid, err)
	}
}

// TestGetCart get a cart from its ID.
func TestGetCart(t *testing.T) {
	ctx := context.Background()
	cid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "cart:"+cid, "cid", cid)
	db.Redis.Expire(ctx, "cart:"+cid, conf.CartDuration)

	pid, _ := stringutil.Random()
	title := faker.Sentence()
	db.Redis.HSet(ctx, "product:"+pid, "title", title)

	c, err := Get(cid)
	// TODO: Add the test len(c.Products) == 0
	if c.ID == "" || err != nil {
		t.Fatalf("Get(cid) = %v, %v, want Cart, nil", c, err)
	}
}

// TestGetCartWithCIDNotExisting returns an error because the
// CID does not exist.
func TestGetCartWithCIDNotExisting(t *testing.T) {
	c, err := Get("toto")
	if c.ID != "" || len(c.Products) != 0 || err == nil || err.Error() != "cart_not_found" {
		t.Fatalf("Get(cid) = %v, %v, want Cart{}, 'cart_not_found'", c, err)
	}
}

// TestCartUpdateDelivery updates the cart delivery.
func TestCartUpdateDelivery(t *testing.T) {
	ctx := context.Background()
	cid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "cart:"+cid, "cid", cid)
	c := Cart{ID: cid}

	if err := c.UpdateDelivery("collect"); err != nil {
		t.Fatalf("c.UpdateDelivery('collect') = %v, want nil", err)
	}
}

// TestCartUpdateDeliveryWithMisvalue returns an error because
// the delivery value is wrong.
func TestCartUpdateDeliveryWithMisvalue(t *testing.T) {
	ctx := context.Background()
	cid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "cart:"+cid, "cid", cid)
	c := Cart{ID: cid}

	if err := c.UpdateDelivery("toto"); err == nil || err.Error() != "unauthorized" {
		t.Fatalf("c.UpdateDelivery('toto') = %v, want unauthorized", err)
	}
}

// TestCartUpdatePayment updates the cart payment.
func TestCartUpdatePayment(t *testing.T) {
	ctx := context.Background()
	cid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "cart:"+cid, "cid", cid)
	c := Cart{ID: cid}

	if err := c.UpdatePayment("card"); err != nil {
		t.Fatalf("c.UpdatePayment('cart') = %v, want nil", err)
	}
}

// TestCartUpdatePaymentWithMisvalue returns an error because
// the delivery value is wrong.
func TestCartUpdatePaymentWithMisvalue(t *testing.T) {
	ctx := context.Background()
	cid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "cart:"+cid, "cid", cid)
	c := Cart{ID: cid}

	if err := c.UpdatePayment("toto"); err == nil || err.Error() != "unauthorized" {
		t.Fatalf("c.UpdatePayment('toto') = %v, want unauthorized", err)
	}
}
