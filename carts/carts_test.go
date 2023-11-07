package carts

import (
	"gifthub/string/stringutil"
	"gifthub/tests"
	"testing"
)

var cart Cart = Cart{ID: "test"}

func TestAddReturnsNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	quantity := 1

	if err := Add(ctx, "test", "test", quantity); err != nil {
		t.Fatalf(`Add(ctx, "test", "test", quantity), %v, want nil, error`, err)
	}
}

func TestAddReturnsErrorWhenCidIsInvalid(t *testing.T) {
	ctx := tests.Context()
	quantity := 1

	if err := Add(ctx, "world", "test", quantity); err == nil {
		t.Fatalf(`Add(ctx, "test", "test", quantity), %v, not want nil, error`, err)
	}
}

func TestAddReturnsErrorWhenPidIsInvalid(t *testing.T) {
	pid, _ := stringutil.Random()
	ctx := tests.Context()
	quantity := 1

	if err := Add(ctx, "test", pid, quantity); err == nil {
		t.Fatalf(`Add(ctx, "test", pid, quantity), %v, not want nil, error`, err)
	}
}

func TestRefreshCIDReturnsCidWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	cid, err := RefreshCID(ctx, "")
	if cid == "" || err != nil {
		t.Fatalf("UpdateCID(ctx, '') = %s, %v, not want ', error", cid, err)
	}
}

func TestRefreshCIDReturnsCidWhenCidExisting(t *testing.T) {
	ctx := tests.Context()
	s, _ := stringutil.Random()
	cid, err := RefreshCID(ctx, s)
	if cid == "" || err != nil {
		t.Fatalf("UpdateCID(ctx, s) = %s, %v, want string, error", cid, err)
	}
}

func TestGetReturnsCartWhenSuccess(t *testing.T) {
	ctx := tests.Context()

	c, err := Get(ctx, "test")
	if c.ID == "" || err != nil {
		t.Fatalf(`Get(ctx, "test") = %v, %v, want Cart, nil`, c, err)
	}
}

func TestGetReturnsErrorWhenCidIsNotExisting(t *testing.T) {
	ctx := tests.Context()
	c, err := Get(ctx, "toto")
	if c.ID != "" || len(c.Products) != 0 || err == nil || err.Error() != "cart_not_found" {
		t.Fatalf(`Get(ctx, "toto") = %v, %v, want Cart{}, 'cart_not_found'`, c, err)
	}
}

func TestUpdateDeliveryReturnsNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	if err := cart.UpdateDelivery(ctx, "collect"); err != nil {
		t.Fatalf("cart.UpdateDelivery(ctx,'collect') = %v, want nil", err)
	}
}

func TestUpdateDeliveryWhenDeliveryIsInvalid(t *testing.T) {
	ctx := tests.Context()
	if err := cart.UpdateDelivery(ctx, "toto"); err == nil || err.Error() != "unauthorized" {
		t.Fatalf("cart.UpdateDelivery(ctx,'toto') = %v, want unauthorized", err)
	}
}

func TestUpdatePaymentReturnsNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	if err := cart.UpdatePayment(ctx, "card"); err != nil {
		t.Fatalf("cart.UpdatePayment(ctx, 'card') = %v, want nil", err)
	}
}

func TestUpdatePaymentReturnsErrorWhenPaymentIsInvalid(t *testing.T) {
	ctx := tests.Context()
	if err := cart.UpdatePayment(ctx, "toto"); err == nil || err.Error() != "unauthorized" {
		t.Fatalf("cart.UpdatePayment(ctx, 'toto') = %v, want 'unauthorized'", err)
	}
}
