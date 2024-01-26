package users

import (
	"fmt"
	"gifthub/db"
	"gifthub/tests"
	"slices"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func init() {
	ctx := tests.Context()

	now := time.Now()

	db.Redis.ZAdd(ctx, fmt.Sprintf("wish:%d", user.ID), redis.Z{
		Member: "PDT1",
		Score:  float64(now.Unix()),
	})

	db.Redis.ZAdd(ctx, fmt.Sprintf("wish:%d", user.ID), redis.Z{
		Member: "PDT2",
		Score:  float64(now.Unix()),
	})

	db.Redis.ZAdd(ctx, fmt.Sprintf("wish:%d", user.ID), redis.Z{
		Member: "1",
		Score:  float64(now.Unix()),
	})
}

func TestWishAddReturnsNilSuccess(t *testing.T) {
	ctx := tests.Context()
	if err := user.Wish(ctx, "1"); err != nil {
		t.Fatalf(`user.Wish(ctx, "1") = %v, want nil`, err)
	}
}

func TestUnWishReturnsNilSuccess(t *testing.T) {
	ctx := tests.Context()
	if err := user.UnWish(ctx, "1"); err != nil {
		t.Fatalf(`user.UnWish(ctx, "1") = %v, want nil`, err)
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

	if !slices.Contains(wishes, "PDT1") {
		t.Fatalf(`slices.Contains(wishes, "PDT1") = %v, want true`, slices.Contains(wishes, "PDT1"))
	}

	if !slices.Contains(wishes, "PDT2") {
		t.Fatalf(`slices.Contains(wishes, "PDT2") = %v, want true`, slices.Contains(wishes, "PDT2"))
	}
}

func TestHasWishReturnsFalseWhenItDoesNotExist(t *testing.T) {
	ctx := tests.Context()
	if has := user.HasWish(ctx, "123"); has {
		t.Fatalf(`user.HasWish(ctx, "123") = %v, want false`, has)
	}
}

func TestHasWishReturnsTrueWhenItExists(t *testing.T) {
	ctx := tests.Context()
	if has := user.HasWish(ctx, "PDT1"); !has {
		t.Fatalf(`user.HasWish(ctx, "PDT1") = %v, want true`, has)
	}
}
