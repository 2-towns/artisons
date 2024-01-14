package admin

import (
	"gifthub/tests"
	"mime/multipart"
	"testing"
)

var value = map[string][]string{
	"price":       {""},
	"quantity":    {""},
	"discount":    {""},
	"weight":      {""},
	"tags":        {""},
	"title":       {""},
	"description": {""},
	"sku":         {""},
	"status":      {""},
}

func TestProcessProductFormReturnsErrorWhenPriceIsEmpty(t *testing.T) {
	c := tests.Context()

	f := multipart.Form{Value: value}

	if _, err := processProductFrom(c, f, ""); err == nil || err.Error() != "input:price" {
		t.Fatalf(`processProductFrom(c, f, "") = _, %v, want _, 'input:price'`, err.Error())
	}
}

func TestProcessProductFormReturnsErrorWhenQuantityIsEmpty(t *testing.T) {
	c := tests.Context()

	v := make(map[string][]string)
	for k, val := range value {
		v[k] = val
	}

	v["price"] = []string{"12.4"}
	f := multipart.Form{Value: v}

	if _, err := processProductFrom(c, f, ""); err == nil || err.Error() != "input:quantity" {
		t.Fatalf(`processProductFrom(c, f, "") = _, %v, want _, 'input:quantity'`, err.Error())
	}
}

func TestProcessProductFormReturnsErrorWhenPriceIsInvalid(t *testing.T) {
	c := tests.Context()

	v := make(map[string][]string)
	for k, val := range value {
		v[k] = val
	}

	v["price"] = []string{"hello"}
	v["quantity"] = []string{"1"}
	f := multipart.Form{Value: v}

	if _, err := processProductFrom(c, f, ""); err == nil || err.Error() != "input:price" {
		t.Fatalf(`processProductFrom(c, f, "") = _, %v, want _, 'input:price'`, err.Error())
	}
}

func TestProcessProductFormReturnsErrorWhenQuantityIsInvalid(t *testing.T) {
	c := tests.Context()

	v := make(map[string][]string)
	for k, val := range value {
		v[k] = val
	}

	v["price"] = []string{"12.5"}
	v["quantity"] = []string{"hello"}
	f := multipart.Form{Value: v}

	if _, err := processProductFrom(c, f, ""); err == nil || err.Error() != "input:quantity" {
		t.Fatalf(`processProductFrom(c, f, "") = _, %v, want _, 'input:quantity'`, err.Error())
	}
}

func TestProcessProductFormReturnsErrorWhenDiscountIsInvalid(t *testing.T) {
	c := tests.Context()

	v := make(map[string][]string)
	for k, val := range value {
		v[k] = val
	}

	v["price"] = []string{"12.5"}
	v["quantity"] = []string{"1"}
	v["discount"] = []string{"hello"}
	f := multipart.Form{Value: v}

	if _, err := processProductFrom(c, f, ""); err == nil || err.Error() != "input:discount" {
		t.Fatalf(`processProductFrom(c, f, "") = _, %v, want _, 'input:discount'`, err.Error())
	}
}

func TestProcessProductFormReturnsErrorWhenWeightIsInvalid(t *testing.T) {
	c := tests.Context()

	v := make(map[string][]string)
	for k, val := range value {
		v[k] = val
	}

	v["price"] = []string{"12.5"}
	v["quantity"] = []string{"1"}
	v["weight"] = []string{"hello"}
	f := multipart.Form{Value: v}

	if _, err := processProductFrom(c, f, ""); err == nil || err.Error() != "input:weight" {
		t.Fatalf(`processProductFrom(c, f, "") = _, %v, want _, 'input:weight'`, err.Error())
	}
}

func TestProcessProductFormReturnsErrorWhenAFieldIsInvalid(t *testing.T) {
	c := tests.Context()

	v := make(map[string][]string)
	for k, val := range value {
		v[k] = val
	}

	v["price"] = []string{"12.5"}
	v["quantity"] = []string{"1"}
	f := multipart.Form{Value: v}

	if _, err := processProductFrom(c, f, ""); err == nil || err.Error() != "input:title" {
		t.Fatalf(`processProductFrom(c, f, "") = _, %v, want _, 'input:title'`, err.Error())
	}
}

func TestProcessProductFormReturnsErrorWhenNotPicture(t *testing.T) {
	c := tests.Context()

	v := make(map[string][]string)
	for k, val := range value {
		v[k] = val
	}

	v["price"] = []string{"12.5"}
	v["quantity"] = []string{"1"}
	v["title"] = []string{"title"}
	v["description"] = []string{"description"}
	v["sku"] = []string{"sku"}
	v["status"] = []string{"online"}

	f := multipart.Form{Value: v}

	if _, err := processProductFrom(c, f, ""); err == nil || err.Error() != "input:image_1" {
		t.Fatalf(`processProductFrom(c, f, "") = _, %v, want _, 'input:image_1'`, err.Error())
	}
}

func TestProcessProductFormReturnsProductWhenEditingWithNotPicture(t *testing.T) {
	c := tests.Context()

	v := make(map[string][]string)
	for k, val := range value {
		v[k] = val
	}

	v["price"] = []string{"12.5"}
	v["quantity"] = []string{"1"}
	v["title"] = []string{"title"}
	v["description"] = []string{"description"}
	v["sku"] = []string{"sku"}
	v["status"] = []string{"online"}

	f := multipart.Form{Value: v}

	if _, err := processProductFrom(c, f, "123"); err != nil {
		t.Fatalf(`processProductFrom(c, f, "123") = _, %v, want _, nil`, err)
	}
}
