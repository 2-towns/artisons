package products

import (
	"context"
	"gifthub/db"
	"gifthub/locales"
	"gifthub/string/stringutil"
	"gifthub/tests"
	"os"
	"testing"
)

var product = Product{
	ID:          "123",
	Title:       "Title",
	Description: "Description",
	Price:       32.5,
	Slug:        "title",
	MID:         "12345",
	Sku:         "123456",
	Currency:    "EUR",
	Quantity:    1,
	Length:      1,
	Status:      "online",
	Weight:      1.5,
}

func TestMain(m *testing.M) {
	// Write code here to run before tests
	locales.LoadEn()

	// Run tests
	exitVal := m.Run()

	// Write code here to run after tests

	// Exit with exit value from tests
	os.Exit(exitVal)
}

// TestImagePath expects to succeed
func TestImagePath(t *testing.T) {
	pid := "123"
	index := 1
	_, p := ImagePath(pid, index)
	expected := "../web/images/123/1"

	if p != expected {
		t.Fatalf(`TestImagePath("123", 1) = %s, want %s`, p, expected)
	}
}

// TestProductAvailable expects to succeed when the product exists
func TestProductAvailable(t *testing.T) {
	ctx := context.Background()
	pid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "product:"+pid, "status", "online")
	c := tests.Context()

	if exists := Available(c, pid); !exists {
		t.Fatalf(`Available(pid) = %v, want true`, exists)
	}
}

// TestProductAvailableNotFound expects to fail because of product non existence
func TestProductAvailableNotFound(t *testing.T) {
	c := tests.Context()

	if exists := Available(c, "toto"); exists {
		t.Fatalf(`Available(c, pid) = %v, want false`, exists)
	}
}

// TestProductsAvailables expects to succeed
func TestProductsAvailables(t *testing.T) {
	ctx := context.Background()
	pid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "product:"+pid, "status", "online")
	c := tests.Context()

	if exists := Availables(c, []string{pid}); !exists {
		t.Fatalf(`Availables(c, pid) = %v, want true`, exists)
	}
}

// TestProductsAvailablesNotFound expects to fail because of products non existence
func TestProductsAvailablesNotFound(t *testing.T) {
	c := tests.Context()

	if exists := Availables(c, []string{"toto"}); exists {
		t.Fatalf(`Availables(c, pid) = %v, want false`, exists)
	}
}

func TestValidateSkuEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Sku = ""

	if err := p.Validate(c); err == nil || err.Error() != "The field sku is required." {
		t.Fatalf(`p.Validate(c) = %v, want not "The field sku is required."`, err.Error())
	}
}

func TestValidateSkuMalFormated(t *testing.T) {
	c := tests.Context()

	p := product
	p.Sku = "!!!"

	if err := p.Validate(c); err == nil || err.Error() != "The field sku is not correct." {
		t.Fatalf(`p.Validate(c) = %v, want not "The field sku is not correct."`, err.Error())
	}
}

func TestValidateTitleEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Title = ""

	if err := p.Validate(c); err == nil || err.Error() != "The field title is required." {
		t.Fatalf(`p.Validate(c) = %v, want not "The field title is required."`, err.Error())
	}
}

func TestValidateDescriptionEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Description = ""

	if err := p.Validate(c); err == nil || err.Error() != "The field description is required." {
		t.Fatalf(`p.Validate(c) = %v, want not "The field description is required."`, err.Error())
	}
}

func TestValidateCurrencyEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Currency = ""

	if err := p.Validate(c); err == nil || err.Error() != "The field currency is not correct." {
		t.Fatalf(`p.Validate(c) = %v, want not "The field currency is not correct."`, err.Error())
	}
}

func TestValidateCurrencyNotSupported(t *testing.T) {
	c := tests.Context()

	p := product
	p.Currency = "ABC"

	if err := p.Validate(c); err == nil || err.Error() != "The field currency is not correct." {
		t.Fatalf(`p.Validate(c) = %v, want not "The field currency is not correct."`, err.Error())
	}
}

func TestValidateStatusEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Status = ""

	if err := p.Validate(c); err == nil || err.Error() != "The field status is not correct." {
		t.Fatalf(`p.Validate(c) = %v, want not "The field status is not correct."`, err.Error())
	}
}

func TestValidateStatusNotSupported(t *testing.T) {
	c := tests.Context()

	p := product
	p.Status = "ABC"

	if err := p.Validate(c); err == nil || err.Error() != "The field status is not correct." {
		t.Fatalf(`p.Validate(c) = %v, want not "The field status is not correct."`, err.Error())
	}
}

func TestValidateLengthEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Length = 0

	if err := p.Validate(c); err == nil || err.Error() != "The field length is required." {
		t.Fatalf(`p.Validate(c) = %v, want not "The field length required."`, err.Error())
	}
}

func TestValidate(t *testing.T) {
	c := tests.Context()

	if err := product.Validate(c); err != nil {
		t.Fatalf(`p.Validate(c) = %v, want nil`, err.Error())
	}
}

func TestSave(t *testing.T) {
	c := tests.Context()

	if err := product.Save(c); err != nil {
		t.Fatalf(`p.Save(c) = %v, want nil`, err.Error())
	}
}

func TestSaveWithoutPID(t *testing.T) {
	c := tests.Context()
	p := Product{ID: ""}
	if err := p.Save(c); err == nil {
		t.Fatalf(`p.Save(c) = %v, want "product_pid_required"`, err.Error())
	}
}

func TestFindWithoutPID(t *testing.T) {
	c := tests.Context()
	if _, err := Find(c, ""); err == nil || err.Error() != "product_not_found" {
		t.Fatalf(`Find(c,"") = %v, want "product_not_found"`, err.Error())
	}
}

func TestFindWithPIDNotExisting(t *testing.T) {
	c := tests.Context()
	if _, err := Find(c, ""); err == nil || err.Error() != "product_not_found" {
		t.Fatalf(`Find(c,"hello") = %v, want "product_not_found"`, err.Error())
	}
}

func TestFind(t *testing.T) {
	c := tests.Context()
	p, err := Find(c, "test")
	if err != nil {
		t.Fatalf(`Find(c,"test") = %v, want "product_not_found"`, err.Error())
	}

	if p.Sku == "" {
		t.Fatalf(`p.Sku = %v, want string`, p.Sku)
	}

	if p.Currency != "EUR" {
		t.Fatalf(`p.Currency = %v, want string`, p.Currency)
	}

	if p.Description == "" {
		t.Fatalf(`p.Description  = %v, want string`, p.Description)
	}

	if p.MID == "" {
		t.Fatalf(`p.MID = %v, want string`, p.MID)
	}

	if p.ID != "test" {
		t.Fatalf(`p.PID = %v, want "test"`, p.ID)
	}

	if p.Slug == "" {
		t.Fatalf(`p.Slug = %v, want string`, p.Slug)
	}

	if p.Status != "online" {
		t.Fatalf(`p.Status = %v, want string`, p.Status)
	}

	if p.Title == "" {
		t.Fatalf(`p.Title = %v, want string`, p.Title)
	}

	if p.Length <= 0 {
		t.Fatalf(`p.Length = %v, want >0`, p.Length)
	}

	if p.Length != len(p.Images) {
		t.Fatalf(`p.Images = %d, want "%d"`, len(p.Images), p.Length)
	}

	if p.Links[0] != "http://google.fr" {
		t.Fatalf(`p.Links[0] = %v, want "http://google.fr"`, p.Links[0])
	}

	if p.Tags[0] != "gift" {
		t.Fatalf(`p.Tags[0] = %v, want "gift"`, p.Tags[0])
	}

	if p.Meta["color"] != "blue" {
		t.Fatalf(`p.Meta["color"] = %v, want "blue"`, p.Meta["color"])
	}

}
