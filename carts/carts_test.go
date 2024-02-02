package carts

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/string/stringutil"
	"artisons/tests"
	"context"
	"testing"
)

func init() {
	ctx := tests.Context()

	db.Redis.HSet(ctx, "cart:99", "cid", "CAR99")
	db.Redis.Expire(ctx, "cart:99", conf.CartDuration)

	db.Redis.HSet(ctx, "product:PDT97",
		"id", "PDT97",
		"status", "online",
	)
}

var cart Cart = Cart{ID: "99"}

func TestAddReturnsNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.Cart, cart.ID)
	quantity := 1

	if cid, err := Add(ctx, "PDT97", quantity); cid == "" || err != nil {
		t.Fatalf(`Add(ctx, "PDT97", quantity), %v, want nil, error`, err)
	}
}

func TestAddReturnsSuccessWhenCidDoesNotExist(t *testing.T) {
	ctx := tests.Context()
	quantity := 1

	if cid, err := Add(ctx, "PDT97", quantity); cid == "" || err != nil {
		t.Fatalf(`Add(ctx, "PDT97", quantity), %v, want nil, error`, err)
	}
}

func TestAddReturnsErrorWhenPidIsInvalid(t *testing.T) {
	pid, _ := stringutil.Random()
	ctx := tests.Context()
	quantity := 1

	if cid, err := Add(ctx, pid, quantity); cid != "" || err == nil {
		t.Fatalf(`Add(ctx, pid, quantity), %v, not want nil, error`, err)
	}
}

func TestRefreshCIDReturnsCidWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	cid, err := RefreshCID(ctx, "CAR99")
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
	ctx = context.WithValue(ctx, contexts.Cart, cart.ID)

	c, err := Get(ctx)
	if c.ID == "" || err != nil {
		t.Fatalf(`Get(ctx) = %v, %v, want Cart, nil`, c, err)
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
	if err := cart.UpdateDelivery(ctx, "toto"); err == nil || err.Error() != "you are not authorized to process this request" {
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
	if err := cart.UpdatePayment(ctx, "toto"); err == nil || err.Error() != "you are not authorized to process this request" {
		t.Fatalf("cart.UpdatePayment(ctx, 'toto') = %v, want 'unauthorized'", err)
	}
}
