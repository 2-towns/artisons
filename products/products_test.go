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

func TestImagePathReturnsCorrectPathWhenSuccess(t *testing.T) {
	pid := "123"
	index := 1
	_, p := ImagePath(pid, index)
	expected := "../web/images/123/1"

	if p != expected {
		t.Fatalf(`TestImagePath("123", 1) = %s, want %s`, p, expected)
	}
}

func TestAvailableReturnsTrueWhenSuccess(t *testing.T) {
	ctx := context.Background()
	pid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "product:"+pid, "status", "online")
	c := tests.Context()

	if exists := Available(c, pid); !exists {
		t.Fatalf(`Available(pid) = %v, want true`, exists)
	}
}

func TestAvailableReturnsFalseWhenProductIsNotFound(t *testing.T) {
	c := tests.Context()

	if exists := Available(c, "toto"); exists {
		t.Fatalf(`Available(c, pid) = %v, want false`, exists)
	}
}

func TestAvailablesReturnsTrueWhenSuccess(t *testing.T) {
	ctx := context.Background()
	pid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "product:"+pid, "status", "online")
	c := tests.Context()

	if exists := Availables(c, []string{pid}); !exists {
		t.Fatalf(`Availables(c, pid) = %v, want true`, exists)
	}
}

func TestAvailablesReturnsFalseWhenProductsAreNotFound(t *testing.T) {
	c := tests.Context()

	if exists := Availables(c, []string{"toto"}); exists {
		t.Fatalf(`Availables(c, pid) = %v, want false`, exists)
	}
}

func TestValidateReturnsErrorWhenSkuIsEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Sku = ""

	if err := p.Validate(c); err == nil || err.Error() != "product_sku_invalid" {
		t.Fatalf(`p.Validate(c) = %v, want not "product_sku_invalid"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenSkuIsInvalid(t *testing.T) {
	c := tests.Context()

	p := product
	p.Sku = "!!!"

	if err := p.Validate(c); err == nil || err.Error() != "product_sku_invalid" {
		t.Fatalf(`p.Validate(c) = %v, want not "product_sku_invalid"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenTitleIsEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Title = ""

	if err := p.Validate(c); err == nil || err.Error() != "product_title_invalid" {
		t.Fatalf(`p.Validate(c) = %v, want not "product_title_invalid"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenDescriptionIsEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Description = ""

	if err := p.Validate(c); err == nil || err.Error() != "product_description_invalid" {
		t.Fatalf(`p.Validate(c) = %v, want not "product_description_invalid"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenCurrencyIsEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Currency = ""

	if err := p.Validate(c); err == nil || err.Error() != "product_currency_invalid" {
		t.Fatalf(`p.Validate(c) = %v, want not "product_currency_invalid"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenCurrencyIsNotSupported(t *testing.T) {
	c := tests.Context()

	p := product
	p.Currency = "ABC"

	if err := p.Validate(c); err == nil || err.Error() != "product_currency_invalid" {
		t.Fatalf(`p.Validate(c) = %v, want not "product_currency_invalid"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenStatusIsEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Status = ""

	if err := p.Validate(c); err == nil || err.Error() != "product_status_invalid" {
		t.Fatalf(`p.Validate(c) = %v, want not "product_status_invalid"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenStatusIsNotSupported(t *testing.T) {
	c := tests.Context()

	p := product
	p.Status = "ABC"

	if err := p.Validate(c); err == nil || err.Error() != "product_status_invalid" {
		t.Fatalf(`p.Validate(c) = %v, want not "product_status_invalid"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenLengthIsEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Length = 0

	if err := p.Validate(c); err == nil || err.Error() != "product_length_invalid" {
		t.Fatalf(`p.Validate(c) = %v, want not "product_length_invalid`, err.Error())
	}
}

func TestValidateReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	if err := product.Validate(c); err != nil {
		t.Fatalf(`p.Validate(c) = %v, want nil`, err.Error())
	}
}

func TestSaveReturnsNilWhenSuccess(t *testing.T) {
	c := tests.Context()

	if err := product.Save(c); err != nil {
		t.Fatalf(`p.Save(c) = %v, want nil`, err.Error())
	}
}

func TestSaveReturnsErrorWhenPidIsEmpty(t *testing.T) {
	c := tests.Context()
	p := Product{ID: ""}
	if err := p.Save(c); err == nil {
		t.Fatalf(`p.Save(c) = %v, want "product_pid_required"`, err.Error())
	}
}

func TestFindReturnsErrorWhenPidIsMissing(t *testing.T) {
	c := tests.Context()
	if _, err := Find(c, ""); err == nil || err.Error() != "product_id_required" {
		t.Fatalf(`Find(c,"") = %v, want "product_id_required"`, err.Error())
	}
}

func TestFindReturnsErrorWhenPidDoesNotExist(t *testing.T) {
	c := tests.Context()
	if _, err := Find(c, ""); err == nil || err.Error() != "product_id_required" {
		t.Fatalf(`Find(c,"hello") = %v, want "product_id_required"`, err.Error())
	}
}

func TestFindReturnsProductWhenSuccess(t *testing.T) {
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

	if p.Length == 0 {
		t.Fatalf(`p.Length = %v, want > 0`, p.Length)
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

func TestSearchReturnsProductsWhenTitleIsFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keywords: "pull"})
	if err != nil {
		t.Fatalf(`Find(c,"test") = %v, want "product_not_found"`, err.Error())
	}

	if len(p) == 0 {
		t.Fatalf(`len(p) = %d, want > 0`, len(p))
	}

	if p[0].ID != "test" {
		t.Fatalf(`p[0].ID = %s, want "test"`, p[0].ID)
	}
}

func TestSearchReturnsProductsWhenDescriptionIsFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keywords: "Lorem"})
	if err != nil {
		t.Fatalf(`Find(c,"test") = %v, want "product_not_found"`, err.Error())
	}

	if len(p) == 0 {
		t.Fatalf(`len(p) = %d, want > 0`, len(p))
	}

	if p[0].ID != "test" {
		t.Fatalf(`p[0].ID = %s, want "test"`, p[0].ID)
	}
}

func TestSearchReturnsProductsWhenSkuIsFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keywords: "skutest"})
	if err != nil {
		t.Fatalf(`Find(c,"test") = %v, want "product_not_found"`, err.Error())
	}

	if len(p) == 0 {
		t.Fatalf(`len(p) = %d, want > 0`, len(p))
	}

	if p[0].ID != "test" {
		t.Fatalf(`p[0].ID = %s, want "test"`, p[0].ID)
	}
}

func TestSearchReturnsEmptySliceWhenKeywordIsNotFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keywords: "crazy"})
	if err != nil {
		t.Fatalf(`Find(c,"test") = %v, want "product_not_found"`, err.Error())
	}

	if len(p) > 0 {
		t.Fatalf(`len(p) = %d, want 0`, len(p))
	}
}

func TestSearchReturnsProductsWhenPriceIsMoreThanPriceMin(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{PriceMin: 50})
	if err != nil {
		t.Fatalf(`Find(c,"test") = %v, want "product_not_found"`, err.Error())
	}

	if len(p) == 0 {
		t.Fatalf(`len(p) = %d, want > 0`, len(p))
	}

	if p[0].ID != "test" {
		t.Fatalf(`p[0].ID = %s, want "test"`, p[0].ID)
	}
}

func TestSearchReturnsEmptySliceWhenPriceMinIsOutOfRange(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{PriceMin: 1500000000})
	if err != nil {
		t.Fatalf(`Find(c,"test") = %v, want "product_not_found"`, err.Error())
	}

	if len(p) != 0 {
		t.Fatalf(`len(p) = %d, want 0`, len(p))
	}
}

func TestSearchReturnsProductsWhenPriceIsLessThanPriceMax(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{PriceMax: 150})
	if err != nil {
		t.Fatalf(`Find(c,"test") = %v, want "product_not_found"`, err.Error())
	}

	if len(p) == 0 {
		t.Fatalf(`len(p) = %d, want > 0`, len(p))
	}
}

func TestSearchReturnsEmptySliceWhenPriceMaxIsOutOfRange(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{PriceMax: 0.5})
	if err != nil {
		t.Fatalf(`Find(c,"test") = %v, want "product_not_found"`, err.Error())
	}

	if len(p) != 0 {
		t.Fatalf(`len(p) = %d, want 0`, len(p))
	}
}

func TestSearchReturnsProductsWhenTagsAreFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Tags: []string{"gift"}})
	if err != nil {
		t.Fatalf(`Find(c,"test") = %v, want "product_not_found"`, err.Error())
	}

	if len(p) == 0 {
		t.Fatalf(`len(p) = %d, want > 0`, len(p))
	}

	if p[0].ID != "test" {
		t.Fatalf(`p[0].ID = %s, want "test"`, p[0].ID)
	}
}

func TestSearchReturnsEmptySliceWhenTagsAreNotFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Tags: []string{"crazy"}})
	if err != nil {
		t.Fatalf(`Find(c,"test") = %v, want "product_not_found"`, err.Error())
	}

	if len(p) != 0 {
		t.Fatalf(`len(p) = %d, want 0`, len(p))
	}
}

func TestSearchReturnsProductsWhenMetaAreFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Meta: map[string]string{"color": "blue"}})
	if err != nil {
		t.Fatalf(`Find(c,"test") = %v, want "product_not_found"`, err.Error())
	}

	if len(p) == 0 {
		t.Fatalf(`len(p) = %d, want > 0`, len(p))
	}

	if p[0].ID != "test" {
		t.Fatalf(`p[0].ID = %s, want "test"`, p[0].ID)
	}
}

func TestSearchReturnsEmptySliceWhenMetaAreNotFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Meta: map[string]string{"color": "crazy"}})
	if err != nil {
		t.Fatalf(`Find(c,"test") = %v, want "product_not_found"`, err.Error())
	}

	if len(p) != 0 {
		t.Fatalf(`len(p) = %d, want 0`, len(p))
	}
}

func TestURLReturnsTheProductURLWhenSuccess(t *testing.T) {
	if product.URL() != "http://localhost/123-title" {
		t.Fatalf(`product.URL()  = %s, want 'http://localhost/123-title'`, product.URL())
	}
}
