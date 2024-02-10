package carts

import (
	"artisons/db"
	"artisons/http/contexts"
	"artisons/tests"
	"context"
	"path"
	"runtime"
	"testing"
)

var cur string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	cur = path.Dir(filename) + "/"
}

func TestAdd(t *testing.T) {
	ctx := tests.Context()
	c := Cart{ID: ""}

	tests.ImportData(ctx, cur+"testdata/cart.redis")

	t.Run("Product not found", func(t *testing.T) {
		qty := 1

		if cid, err := c.Add(ctx, "idontexist", qty); err == nil || err.Error() != "oops the data is not found" {
			t.Fatalf(`cid = %s, err = %v, want "", nil, oops the data is not found`, cid, err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		qty := 1

		if cid, err := c.Add(ctx, "PDT1", qty); cid == "" || err != nil {
			t.Fatalf(`cid = %s, err = %v, want not empty, nil`, cid, err)
		}
	})
}

func TestGet(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/cart.redis")

	t.Run("Empty", func(t *testing.T) {
		c := Cart{ID: ""}
		cart, err := c.Get(ctx)

		if err != nil {
			t.Fatalf(`err = %v, want nil`, err)
		}

		if cart.ID != "" {
			t.Fatalf(`id = %s, want ""`, cart.ID)
		}

		if len(cart.Products) != 0 {
			t.Fatalf(`len(products) = %d, want 0`, len(cart.Products))
		}
	})

	t.Run("Success", func(t *testing.T) {
		c := Cart{ID: "CAR1"}
		cart, err := c.Get(ctx)

		if err != nil {
			t.Fatalf(`err = %v, want nil`, err)
		}

		if cart.ID == "" {
			t.Fatalf(`id = %s, want not empty`, cart.ID)
		}

		if len(cart.Products) == 0 {
			t.Fatalf(`len(products) = %d, want > 0`, len(cart.Products))
		}
	})
}

func TestDelete(t *testing.T) {
	ctx := tests.Context()
	c := Cart{ID: "CAR1"}

	tests.ImportData(ctx, cur+"testdata/delete.redis")

	t.Run("Does not exist", func(t *testing.T) {
		qty := 1

		if err := c.Delete(ctx, "idontexist", qty); err == nil || err.Error() != "oops the data is not found" {
			t.Fatalf(`err = %v, want "", nil, oops the data is not found`, err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		qty := 1

		if err := c.Delete(ctx, "PDT1", qty); err != nil {
			t.Fatalf(`err = %v, want not empty, nil`, err)
		}

		q, _ := db.Redis.HGet(ctx, "cart:CAR1", "PDT1").Result()
		if q != "1" {
			t.Fatalf(`qty1 = %s, want '1'`, q)
		}
	})
}

func TestRefresh(t *testing.T) {
	ctx := tests.Context()
	cid, err := RefreshCID(ctx, "CAR1")
	if cid == "" || err != nil {
		t.Fatalf(`cid = %s, err = %v, "", nil`, cid, err)
	}
}

func TestUpdateDelivery(t *testing.T) {
	ctx := tests.Context()
	c := Cart{ID: "CAR1"}

	t.Run("Success", func(t *testing.T) {
		if err := c.UpdateDelivery(ctx, "collect"); err != nil {
			t.Fatalf(`err = %v, want nil`, err)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		if err := c.UpdateDelivery(ctx, "iaminvalid"); err == nil || err.Error() != "you are not authorized to process this request" {
			t.Fatalf(`err = %v, want "you are not authorized to process this request"`, err)
		}
	})
}

func TestUpdatePayment(t *testing.T) {
	ctx := tests.Context()
	c := Cart{ID: "CAR1"}

	t.Run("Success", func(t *testing.T) {
		if err := c.UpdatePayment(ctx, "card"); err != nil {
			t.Fatalf(`err = %v, want nil`, err)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		if err := c.UpdatePayment(ctx, "iaminvalid"); err == nil || err.Error() != "you are not authorized to process this request" {
			t.Fatalf(`err = %v, want "you are not authorized to process this request"`, err)
		}
	})
}

func TestExists(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/cart.redis")

	t.Run("Empty", func(t *testing.T) {
		if res := Exists(ctx, ""); res {
			t.Fatalf(`exists = true, want false`)
		}
	})

	t.Run("Does not exist", func(t *testing.T) {
		if res := Exists(ctx, "idontexist"); res {
			t.Fatalf(`exists = true, want false`)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		if res := Exists(ctx, "CAR1"); !res {
			t.Fatalf(`exists = false, want true`)
		}
	})
}

func TestMerge(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/merge.redis")

	t.Run("Anonymous", func(t *testing.T) {
		if err := Merge(ctx, "CAR1"); err == nil || err.Error() != "you are not authorized to process this request" {
			t.Fatalf(`err = %v, want you are not authorized to process this request`, err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		ctx := tests.Context()
		ctx = context.WithValue(ctx, contexts.UserID, 1)

		if err := Merge(ctx, "CAR1"); err != nil {
			t.Fatalf(`Merge(ctx) = %v, want nil`, err)
		}

		qty1, _ := db.Redis.HGet(ctx, "cart:1", "PDT1").Result()
		if qty1 != "3" {
			t.Fatalf(`qty1 = %s, want '3'`, qty1)
		}

		qty2, _ := db.Redis.HGet(ctx, "cart:1", "PDT2").Result()
		if qty2 != "1" {
			t.Fatalf(`qty2 = %s, want '1'`, qty2)
		}
	})
}
