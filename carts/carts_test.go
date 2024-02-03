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

	if _, err := Add(ctx, "PDT97", quantity); err != nil {
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

func TestGetIDReturnsUserIDWhenTheUserIsSignedIn(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.UserID, 123)

	cid, err := GetCID(ctx)
	if err != nil || cid != "123" {
		t.Fatalf(" getCID(ctx) = %s, %v, want '123', nil", cid, err)
	}
}

func TestGetIDReturnsCartIDWhenTheUserIsNotSignedIn(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.Cart, "1221XLS")

	cid, err := GetCID(ctx)
	if err != nil || cid != "1221XLS" {
		t.Fatalf(" getCID(ctx) = %s, %v, want '1221XLS', nil", cid, err)
	}
}

func TestDeleteReturnsNilWhenQuantityIsLowerThanTheCart(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.Cart, cart.ID)
	quantity := 1

	db.Redis.HSet(ctx, "cart:"+cart.ID, "abc", 2).Result()

	if err := Delete(ctx, "abc", quantity); err != nil {
		t.Fatalf(`Delete(ctx, "abc", quantity) = %v, want nil`, err)
	}

	qty, _ := db.Redis.HGet(ctx, "cart:"+cart.ID, "abc").Result()

	if qty != "1" {
		t.Fatalf(`qty = %s, want '1'`, qty)
	}
}

func TestDeleteReturnsNilWhenQuantityIsSameThanTheCart(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.Cart, cart.ID)
	quantity := 2

	db.Redis.HSet(ctx, "cart:"+cart.ID, "abc", 2).Result()

	if err := Delete(ctx, "abc", quantity); err != nil {
		t.Fatalf(`Delete(ctx, "abc", quantity) = %v, want nil`, err)
	}

	qty, _ := db.Redis.HGet(ctx, "cart:"+cart.ID, "abc").Result()

	if qty != "" {
		t.Fatalf(`qty = %s, want ''`, qty)
	}
}

func TestDeleteReturnsNilWhenQuantityIsMoreThanTheCart(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.Cart, cart.ID)
	quantity := 3

	db.Redis.HSet(ctx, "cart:"+cart.ID, "abc", 2).Result()

	if err := Delete(ctx, "abc", quantity); err != nil {
		t.Fatalf(`Delete(ctx, "abc", quantity) = %v, want nil`, err)
	}

	qty, _ := db.Redis.HGet(ctx, "cart:"+cart.ID, "abc").Result()

	if qty != "" {
		t.Fatalf(`qty = %s, want ''`, qty)
	}
}

func TestMergeReturnsErrorWhenUserIdDoesNotExist(t *testing.T) {
	ctx := tests.Context()
	exp := "you are not authorized to process this request"

	if err := Merge(ctx); err == nil || err.Error() != exp {
		t.Fatalf(`Merge(ctx) = %v, want '%s'`, err, exp)
	}
}

func TestMergeReturnsErrorWhenCartIdDoesNotExist(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.UserID, 1)
	ctx = context.WithValue(ctx, contexts.Cart, "CAR1")

	db.Redis.HSet(ctx, "cart:CAR1", "PDT1", "2")
	db.Redis.HSet(ctx, "cart:CAR1", "PDT2", "1")
	db.Redis.HSet(ctx, "cart:1", "PDT1", "1")
	db.Redis.HSet(ctx, "cart:1", "PDT2", "")

	if err := Merge(ctx); err != nil {
		t.Fatalf(`Merge(ctx) = %v, want nil`, err)
	}

	val, _ := db.Redis.HGetAll(ctx, "cart:1").Result()

	if val["PDT1"] != "3" {
		t.Fatalf(`val["PDT1"] = %s, want '3'`, val["PDT1"])
	}

	if val["PDT2"] != "1" {
		t.Fatalf(`val["PDT2"] = %s, want '1'`, val["PDT2"])
	}
}

func TestExistsReturnsFalseWhenTheCartNotInContext(t *testing.T) {
	ctx := tests.Context()

	if res := Exists(ctx, "crazy"); res {
		t.Fatalf(`Exists(ctx, "crazy") = true, want false`)
	}
}

func TestExistsReturnsFalseWhenTheCartDoesNotExist(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.Cart, "CAR1")

	if res := Exists(ctx, "CAR1"); res {
		t.Fatalf(`Exists(ctx, "CAR1") = true, want false`)
	}
}

func TestExistsReturnsTrueWhenTheCartExists(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.Cart, "CAR1")

	db.Redis.HSet(ctx, "cart:CAR1", "PDT1", "1")
	db.Redis.Expire(ctx, "cart:CAR1", conf.CartDuration)

	if res := Exists(ctx, "CAR1"); !res {
		t.Fatalf(`Exists(ctx, "CAR1") = false, want true`)
	}
}
