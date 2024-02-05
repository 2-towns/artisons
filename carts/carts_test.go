package carts

import (
	"artisons/http/contexts"
	"artisons/string/stringutil"
	"artisons/tests"
	"context"
	"fmt"
	"testing"
)

var cart Cart = Cart{ID: tests.CartID}

func TestAddReturnsNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.Cart, cart.ID)
	quantity := 1

	if _, err := Add(ctx, tests.CartProductID, quantity); err != nil {
		t.Fatalf(`Add(ctx, tests.CartProductID, quantity), %v, want nil, error`, err)
	}
}

func TestAddReturnsSuccessWhenCidDoesNotExist(t *testing.T) {
	ctx := tests.Context()
	quantity := 1

	if cid, err := Add(ctx, tests.CartProductID, quantity); cid == "" || err != nil {
		t.Fatalf(`Add(ctx, tests.CartProductID, quantity), %v, want nil, error`, err)
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
	cid, err := RefreshCID(ctx, tests.CartID)
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
	ctx = context.WithValue(ctx, contexts.Cart, tests.CartID)

	if err := UpdateDelivery(ctx, "collect"); err != nil {
		t.Fatalf(`cart.UpdateDelivery(ctx, cid, "collect") = %v, want nil`, err)
	}
}

func TestUpdateDeliveryWhenDeliveryIsInvalid(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.Cart, tests.CartID)

	if err := UpdateDelivery(ctx, "toto"); err == nil || err.Error() != "you are not authorized to process this request" {
		t.Fatalf(`cart.UpdateDelivery(ctx, cid, "toto") = %v, want unauthorized`, err)
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
	ctx = context.WithValue(ctx, contexts.UserID, tests.UserID1)

	cid, err := GetCID(ctx)
	if err != nil || cid != fmt.Sprintf("%d", tests.UserID1) {
		t.Fatalf(" getCID(ctx) = %s, %v, want '%s', nil", cid, err, tests.CartID)
	}
}

func TestGetIDReturnsCartIDWhenTheUserIsNotSignedIn(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.Cart, tests.CartID)

	cid, err := GetCID(ctx)
	if err != nil || cid != tests.CartID {
		t.Fatalf(" getCID(ctx) = %s, %v, want '%s', nil", cid, err, tests.CartID)
	}
}

func TestDeleteReturnsNilWhenQuantityIsLowerThanTheCart(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.Cart, cart.ID)
	quantity := 1

	tests.AddToCart(ctx, cart.ID, tests.ProductID1, "2")

	if err := Delete(ctx, tests.ProductID1, quantity); err != nil {
		t.Fatalf(`Delete(ctx, tests.ProductID1, quantity) = %v, want nil`, err)
	}

	qty := tests.Quantity(ctx, cart.ID, tests.ProductID1)

	if qty != "1" {
		t.Fatalf(`qty = %s, want '1'`, qty)
	}
}

func TestDeleteReturnsNilWhenQuantityIsSameThanTheCart(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.Cart, cart.ID)
	quantity := 2

	tests.AddToCart(ctx, cart.ID, tests.ProductID1, "2")

	if err := Delete(ctx, tests.ProductID1, quantity); err != nil {
		t.Fatalf(`Delete(ctx, tests.ProductID1, quantity) = %v, want nil`, err)
	}

	qty := tests.Quantity(ctx, cart.ID, tests.ProductID1)

	if qty != "" {
		t.Fatalf(`qty = %s, want ''`, qty)
	}
}

func TestDeleteReturnsNilWhenQuantityIsMoreThanTheCart(t *testing.T) {
	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.Cart, cart.ID)
	quantity := 3

	tests.AddToCart(ctx, cart.ID, tests.ProductID1, "2")

	if err := Delete(ctx, tests.ProductID1, quantity); err != nil {
		t.Fatalf(`Delete(ctx, tests.ProductID1, quantity) = %v, want nil`, err)
	}

	qty := tests.Quantity(ctx, cart.ID, tests.ProductID1)

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
	ctx = context.WithValue(ctx, contexts.UserID, tests.UserID1)
	ctx = context.WithValue(ctx, contexts.Cart, cart.ID)
	ucart := fmt.Sprintf("%d", tests.UserID1)

	tests.AddToCart(ctx, cart.ID, tests.ProductID1, "2")
	tests.AddToCart(ctx, cart.ID, tests.ProductID2, "1")
	tests.AddToCart(ctx, ucart, tests.ProductID1, "1")
	tests.AddToCart(ctx, ucart, tests.ProductID2, "")

	if err := Merge(ctx); err != nil {
		t.Fatalf(`Merge(ctx) = %v, want nil`, err)
	}

	qty1 := tests.Quantity(ctx, ucart, tests.ProductID1)
	if qty1 != "3" {
		t.Fatalf(`qty1 = %s, want '3'`, qty1)
	}

	qty2 := tests.Quantity(ctx, ucart, tests.ProductID2)
	if qty2 != "1" {
		t.Fatalf(`qty2 = %s, want '1'`, qty2)
	}
}

func TestExistsReturnsFalseWhenTheCartNotInContext(t *testing.T) {
	ctx := tests.Context()

	if res := Exists(ctx, tests.CartID); res {
		t.Fatalf(`Exists(ctx, tests.DoesNotExist) = true, want false`)
	}
}

func TestExistsReturnsTrueWhenTheCartExists(t *testing.T) {
	ctx := tests.Context()

	tests.AddToCart(ctx, tests.CartID, tests.ProductID1, "1")

	if res := Exists(ctx, tests.CartID); !res {
		t.Fatalf(`Exists(ctx, tests.CartID,) = false, want true`)
	}
}
