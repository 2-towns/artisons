package carts

import (
	"artisons/addresses"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/products"
	"artisons/tests"
	"artisons/users"
	"context"
	"errors"
	"fmt"
	"path"
	"reflect"
	"runtime"
	"testing"

	"github.com/go-faker/faker/v4"
)

var cur string

var ra faker.RealAddress = faker.GetRealAddress()

var address addresses.Address = addresses.Address{
	Lastname:      faker.Name(),
	Firstname:     faker.Name(),
	Street:        ra.Address,
	City:          ra.City,
	Complementary: ra.Address,
	Zipcode:       ra.PostalCode,
	Phone:         faker.Phonenumber(),
}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	cur = path.Dir(filename) + "/"
}

func TestAdd(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/cart.redis")

	t.Run("Product not found", func(t *testing.T) {
		qty := 1

		if err := Add(ctx, 0, "idontexist", qty); err == nil || err.Error() != "oops the data is not found" {
			t.Fatalf(`err = %v, want "", nil, oops the data is not found`, err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		qty := 1

		if err := Add(ctx, 0, "PDT1", qty); err != nil {
			t.Fatalf(`err = %v, want not empty, nil`, err)
		}
	})
}

func TestGet(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/cart.redis")

	t.Run("Empty", func(t *testing.T) {
		cart, err := Get(ctx, 0)

		if err != nil {
			t.Fatalf(`err = %v, want nil`, err)
		}

		if cart.ID != 0 {
			t.Fatalf(`id = %d, want ""`, cart.ID)
		}

		if len(cart.Products) != 0 {
			t.Fatalf(`len(products) = %d, want 0`, len(cart.Products))
		}
	})

	t.Run("Success", func(t *testing.T) {
		cart, err := Get(ctx, 123)

		if err != nil {
			t.Fatalf(`err = %v, want nil`, err)
		}

		if cart.ID == 0 {
			t.Fatalf(`id = %d, want not empty`, cart.ID)
		}

		if len(cart.Products) == 0 {
			t.Fatalf(`len(products) = %d, want > 0`, len(cart.Products))
		}
	})
}

func TestDelete(t *testing.T) {
	ctx := tests.Context()

	t.Run("Does not exist", func(t *testing.T) {
		qty := 1

		if err := Delete(ctx, 123, "idontexist", qty); err == nil || err.Error() != "oops the data is not found" {
			t.Fatalf(`err = %v, want "", nil, oops the data is not found`, err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		tests.ImportData(ctx, cur+"testdata/cart.redis")

		qty := 1

		if err := Delete(ctx, 123, "PDT1", qty); err != nil {
			t.Fatalf(`err = %v, want not empty, nil`, err)
		}

		q, _ := db.Redis.HGet(ctx, "cart:123", "PDT1").Result()
		if q != "1" {
			t.Fatalf(`qty1 = %s, want '1'`, q)
		}
	})
}

func TestRefresh(t *testing.T) {
	ctx := tests.Context()
	err := RefreshCID(ctx, 123)
	if err != nil {
		t.Fatalf(`err = %v, "", nil`, err)
	}
}

func TestSaveAddress(t *testing.T) {
	ctx := tests.Context()
	c := Cart{ID: 123}
	err := c.SaveAddress(ctx, address)
	if err != nil {
		t.Fatalf(`err = %v, "", nil`, err)
	}
}

func TestUpdateDelivery(t *testing.T) {
	ctx := tests.Context()
	c := Cart{ID: 123}

	tests.ImportData(ctx, cur+"testdata/cart.redis")

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
	c := Cart{ID: 123}

	t.Run("Success", func(t *testing.T) {
		tests.ImportData(ctx, cur+"testdata/cart.redis")

		if err := c.UpdatePayment(ctx, "cash"); err != nil {
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
		if res := Exists(ctx, 0); res {
			t.Fatalf(`exists = true, want false`)
		}
	})

	t.Run("Does not exist", func(t *testing.T) {
		if res := Exists(ctx, 9999); res {
			t.Fatalf(`exists = true, want false`)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		if res := Exists(ctx, 123); !res {
			t.Fatalf(`exists = false, want true`)
		}
	})
}

func TestMerge(t *testing.T) {
	ctx := tests.Context()

	t.Run("Anonymous", func(t *testing.T) {
		tests.ImportData(ctx, cur+"testdata/cart.redis")

		if err := Merge(ctx, 123); err == nil || err.Error() != "you are not authorized to process this request" {
			t.Fatalf(`err = %v, want you are not authorized to process this request`, err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		tests.ImportData(ctx, cur+"testdata/cart.redis")

		ctx := tests.Context()
		ctx = context.WithValue(ctx, contexts.User, users.User{ID: 1})

		if err := Merge(ctx, 123); err != nil {
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

func TestValidate(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/cart.redis")

	var address addresses.Address = addresses.Address{
		Lastname:      faker.Name(),
		Firstname:     faker.Name(),
		Street:        ra.Address,
		City:          ra.City,
		Complementary: ra.Address,
		Zipcode:       ra.PostalCode,
		Phone:         faker.Phonenumber(),
	}

	cart := Cart{
		Delivery: "colissimo",
		Total:    100,
		Payment:  "cash",
		Products: []products.Product{{ID: "PDT1"}},
		Address:  address,
	}

	var tests = []struct {
		name  string
		field string
		value interface{}
		err   error
	}{
		{"delivery=idontexist", "Delivery", "idontexist", errors.New("you are not authorized to process this request")},
		{"payment=idontexist", "Payment", "idontexist", errors.New("you are not authorized to process this request")},
		{"products=", "Products", []products.Product{}, errors.New("the cart is empty")},
		{"products={notavailable}", "Products", []products.Product{{ID: "notavailable"}}, errors.New("some products are not available anymore")},
		{"products={PDT1}", "Products", []products.Product{{ID: "PDT1"}}, errors.New("the minimum amount is not reached")},
		{"success", "Products", []products.Product{{ID: "PDT1", Quantity: 10, Price: 100}}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cart

			if tt.field == "Products" {
				c.Products = tt.value.([]products.Product)
			} else if tt.field != "" {
				reflect.ValueOf(&c).Elem().FieldByName(tt.field).SetString(tt.value.(string))
			}

			if err := c.Validate(ctx); fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.err) {
				t.Fatalf(`err = %v, want %s`, err, tt.err)
			}
		})
	}
}

func TestCalculateTotal(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/cart.redis")

	var tests = []struct {
		name  string
		cart  Cart
		total float64
	}{
		{"total=60,delivery=collect", Cart{
			Delivery: "collect",
			Products: []products.Product{{Quantity: 1, Price: 11}, {Quantity: 2, Price: 24.5}},
		}, 60},
		{"total=16.99,delivery=colissimo", Cart{
			Delivery: "colissimo",
			Products: []products.Product{{Quantity: 1, Price: 11}},
		}, 16.99},
		{"total=30,delivery=colissimo", Cart{
			Delivery: "colissimo",
			Products: []products.Product{{Quantity: 1, Price: 30}},
		}, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.cart
			total, err := c.CalculateTotal(ctx)

			if fmt.Sprintf("%2.f", tt.total) != fmt.Sprintf("%2.f", total) || err != nil {
				t.Fatalf(`err = %v, total = %f, want %f`, err, total, tt.total)
			}
		})
	}
}

func TestNewCartID(t *testing.T) {
	ctx := tests.Context()
	if cid, err := NewCartID(ctx); cid == 0 || err != nil {
		t.Fatalf(`cid =%d, err = %v`, cid, err)
	}
}
