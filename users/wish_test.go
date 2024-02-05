package users

import (
	"artisons/tests"
	"fmt"
	"slices"
	"testing"
)

var wuid = fmt.Sprintf("%d", tests.UserID1)

func TestWishAddReturnsNilSuccess(t *testing.T) {
	ctx := tests.Context()

	if err := user.Wish(ctx, wuid); err != nil {
		t.Fatalf(`user.Wish(ctx, wuid) = %v, want nil`, err)
	}
}

func TestUnWishReturnsNilSuccess(t *testing.T) {
	ctx := tests.Context()
	if err := user.UnWish(ctx, wuid); err != nil {
		t.Fatalf(`user.UnWish(ctx, wuid) = %v, want nil`, err)
	}
}

func TestWishesReturnWishes(t *testing.T) {
	ctx := tests.Context()
	wishes, err := user.Wishes(ctx)

	if err != nil {
		t.Fatalf(`user.Wishes(ctx) = %v, %v, want nil`, wishes, err)
	}

	if len(wishes) == 0 {
		t.Fatalf(`len(wishes) = %d, want 0`, len(wishes))
	}

	if !slices.Contains(wishes, tests.ProductID1) {
		t.Fatalf(`slices.Contains(wishes,tests.ProductID1) = %v, want true`, slices.Contains(wishes, tests.ProductID1))
	}

	if !slices.Contains(wishes, tests.ProductID2) {
		t.Fatalf(`slices.Contains(wishes,tests.ProductID2) = %v, want true`, slices.Contains(wishes, tests.ProductID2))
	}
}

func TestHasWishReturnsFalseWhenItDoesNotExist(t *testing.T) {
	ctx := tests.Context()
	if has := user.HasWish(ctx, tests.DoesNotExist); has {
		t.Fatalf(`user.HasWish(ctx, tests.DoesNotExist) = %v, want false`, has)
	}
}

func TestHasWishReturnsTrueWhenItExists(t *testing.T) {
	ctx := tests.Context()
	if has := user.HasWish(ctx, tests.ProductID1); !has {
		t.Fatalf(`user.HasWish(ctx, tests.ProductID1) = %v, want true`, has)
	}
}
