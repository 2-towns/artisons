package products

import (
	"artisons/conf"
	"artisons/db"
	"artisons/locales"
	"artisons/string/stringutil"
	"artisons/tests"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
)

func init() {
	ctx := tests.Context()
	now := time.Now().Unix()

	db.Redis.HSet(ctx, "product:PDT1",
		"id", "PDT1",
		"sku", "SKU1",
		"title", db.Escape("T-shirt Tester c'est douter"),
		"description", db.Escape("T-shirt développeur unisexe Tester c'est douter"),
		"slug", stringutil.Slugify(db.Escape("T-shirt Tester c'est douter")),
		"currency", "EUR",
		"price", 100.5,
		"quantity", rand.Intn(10),
		"status", "online",
		"weight", rand.Float32(),
		"mid", faker.Phonenumber(),
		"meta", "color_blue;color_blue cyan",
		"tags", "clothes",
		"image_1", "PDT1.jpeg",
		"image_2", "PDT1.jpeg",
		"links", "",
		"created_at", now,
		"updated_at", now,
		"type", "product",
	)

	db.Redis.HSet(ctx, "pids", db.Escape("T-shirt développeur unisexe Tester c'est douter"), "PDT1")
}

var product = Product{
	ID:          "123",
	Title:       "Title",
	Description: "Description",
	Price:       32.5,
	Slug:        "title",
	MID:         "12345",
	Sku:         "123456",
	Quantity:    1,
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
	pid := "/products/123.jpeg"
	p := ImagePath(pid)
	expected := conf.WorkingSpace + "web/images/products/123.jpeg"

	if !strings.HasSuffix(p, expected) {
		t.Fatalf(`strings.HasSuffix(p, expected) = %s, want %s`, p, expected)
	}
}

func TestAvailableReturnsTrueWhenSuccess(t *testing.T) {
	ctx := tests.Context()
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
	ctx := tests.Context()
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

func TestValidateReturnsErrorWhenSkuIsInvalid(t *testing.T) {
	c := tests.Context()

	p := product
	p.Sku = "!!!"

	if err := p.Validate(c); err == nil || err.Error() != "input:sku" {
		t.Fatalf(`p.Validate(c) = %v, want not "input:sku"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenTitleIsEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Title = ""

	if err := p.Validate(c); err == nil || err.Error() != "input:title" {
		t.Fatalf(`p.Validate(c) = %v, want not "input:title"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenDescriptionIsEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Description = ""

	if err := p.Validate(c); err == nil || err.Error() != "input:description" {
		t.Fatalf(`p.Validate(c) = %v, want not "input:description"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenStatusIsEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Status = ""

	if err := p.Validate(c); err == nil || err.Error() != "input:status" {
		t.Fatalf(`p.Validate(c) = %v, want not "input:status"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenStatusIsNotSupported(t *testing.T) {
	c := tests.Context()

	p := product
	p.Status = "ABC"

	if err := p.Validate(c); err == nil || err.Error() != "input:status" {
		t.Fatalf(`p.Validate(c) = %v, want not "input:status"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenSlugIsEmpty(t *testing.T) {
	c := tests.Context()

	p := product
	p.Slug = ""

	if err := p.Validate(c); err == nil || err.Error() != "input:slug" {
		t.Fatalf(`p.Validate(c) = %v, want not "input:slug"`, err.Error())
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

	if _, err := product.Save(c); err != nil {
		t.Fatalf(`p.Save(c) = %v, want nil`, err.Error())
	}
}

func TestFindReturnsErrorWhenPidIsMissing(t *testing.T) {
	c := tests.Context()
	if _, err := Find(c, ""); err == nil || err.Error() != "input:id" {
		t.Fatalf(`Find(c,"") = %v, want "input:id"`, err.Error())
	}
}

func TestFindReturnsErrorWhenPidDoesNotExist(t *testing.T) {
	c := tests.Context()
	if _, err := Find(c, "doesnotexist"); err == nil || err.Error() != "oops the data is not found" {
		t.Fatalf(`Find(c, "doesnotexist") = %v, want "oops the data is not found"`, err.Error())
	}
}

func TestFindAllReturnsProductsWhenPidsExist(t *testing.T) {
	c := tests.Context()
	p, err := FindAll(c, []string{"PDT1"})
	if err != nil {
		t.Fatalf(`FindAll(c, []string{"PDT1"}) = %v, %v, want products, nil`, p, err.Error())
	}

	if len(p) == 0 {
		t.Fatal(`len(p) = 0, want 1`)
	}

	if p[0].ID == "" {
		t.Fatal(`p[0].ID = "", want "PDT1"`)
	}
}

func TestFindAllReturnsEmptyArrayWhenPidsDoNotExist(t *testing.T) {
	c := tests.Context()
	p, err := FindAll(c, []string{"crazy"})
	if err != nil {
		t.Fatalf(`FindAll(c, []string{"crazy"}) = %v, %v, want products, nil`, p, err.Error())
	}

	if len(p) > 0 {
		t.Fatalf(`len(p) = %d, want 0`, len(p))
	}
}

func TestFindReturnsProductWhenSuccess(t *testing.T) {
	c := tests.Context()
	p, err := Find(c, "PDT1")
	if err != nil {
		t.Fatalf(`Find(c, "PDT1") = %v, want nil`, err.Error())
	}

	if p.Sku == "" {
		t.Fatalf(`p.Sku = %v, want string`, p.Sku)
	}

	if p.Description == "" {
		t.Fatalf(`p.Description  = %v, want string`, p.Description)
	}

	if p.ID != "PDT1" {
		t.Fatalf(`p.PID = %v, want "PDT1"`, p.ID)
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

	if p.Image1 == "" {
		t.Fatalf(`p.Image1 = %v, want string`, p.Image1)
	}

	if p.Tags[0] != "clothes" {
		t.Fatalf(`p.Tags[0] = %s, want "clothes"`, p.Tags[0])
	}
}

func TestSearchReturnsProductsWhenTitleIsFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keywords: "T-Shirt"}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: "T-Shirt"}) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if len(p.Products) == 0 {
		t.Fatalf(`len(p.Products) = %d, want > 0`, len(p.Products))
	}

	if p.Products[0].ID != "PDT1" {
		t.Fatalf(`p[0].ID = %s, want "PDT1"`, p.Products[0].ID)
	}
}

func TestSearchReturnsNoErrorWhenUsingQuote(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keywords: "c'est"}, 0, 10)

	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: "c'est"}) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if len(p.Products) == 0 {
		t.Fatalf(`len(p.Products) = %d, want > 0`, len(p.Products))
	}

	if p.Products[0].ID != "PDT1" {
		t.Fatalf(`p[0].ID = %s, want "PDT1"`, p.Products[0].ID)
	}
}

func TestSearchReturnsNoErrorWhenUsingDash(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keywords: "t-shirt"}, 0, 10)

	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: "t-shirt"}) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if len(p.Products) == 0 {
		t.Fatalf(`len(p.Products) = %d, want > 0`, len(p.Products))
	}

	if p.Products[0].ID != "PDT1" {
		t.Fatalf(`p[0].ID = %s, want "PDT1"`, p.Products[0].ID)
	}
}

func TestSearchReturnsNoErrorWhenUsingSpace(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keywords: "hello douter"}, 0, 10)

	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: "hello douter"}) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if len(p.Products) == 0 {
		t.Fatalf(`len(p.Products) = %d, want > 0`, len(p.Products))
	}

	if p.Products[0].ID != "PDT1" {
		t.Fatalf(`p[0].ID = %s, want "PDT1"`, p.Products[0].ID)
	}
}

func TestSearchReturnsProductsWhenDescriptionIsFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keywords: "unisexe"}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: "unisexe"}) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if len(p.Products) == 0 {
		t.Fatalf(`len(p.Products) = %d, want > 0`, len(p.Products))
	}

	if p.Products[0].ID != "PDT1" {
		t.Fatalf(`p[0].ID = %s, want "PDT1"`, p.Products[0].ID)
	}
}

func TestSearchReturnsProductsWhenSkuIsFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keywords: "SKU1"}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: "SKU1"}) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if len(p.Products) == 0 {
		t.Fatalf(`len(p.Products) = %d, want > 0`, len(p.Products))
	}

	if p.Products[0].ID != "PDT1" {
		t.Fatalf(`p[0].ID = %s, want "PDT1"`, p.Products[0].ID)
	}
}

func TestSearchReturnsEmptySliceWhenKeywordIsNotFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Keywords: "crazy"}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: "crazy"}) = %v, want nil`, err.Error())
	}

	if p.Total != 0 {
		t.Fatalf(`p.Total = %d, want = 0`, p.Total)
	}

	if len(p.Products) != 0 {
		t.Fatalf(`len(p.Products) = %d, want = 0`, len(p.Products))
	}

	if len(p.Products) > 0 {
		t.Fatalf(`len(p.Products) = %d, want 0`, len(p.Products))
	}
}

func TestSearchReturnsProductsWhenPriceIsMoreThanPriceMin(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{PriceMin: 50}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{PriceMin: 50}) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if len(p.Products) == 0 {
		t.Fatalf(`len(p.Products) = %d, want > 0`, len(p.Products))
	}

	if p.Products[0].ID == "" {
		t.Fatalf(`p.Products[0].ID = %s, want not empty`, p.Products[0].ID)
	}
}

func TestSearchReturnsEmptySliceWhenPriceMinIsOutOfRange(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{PriceMin: 1500000000}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{PriceMin: 1500000000}) = %v, want nil`, err.Error())
	}

	if p.Total != 0 {
		t.Fatalf(`p.Total = %d, want = 0`, p.Total)
	}

	if len(p.Products) != 0 {
		t.Fatalf(`len(p.Products) = %d, want = 0`, len(p.Products))
	}
}

func TestSearchReturnsProductsWhenPriceIsLessThanPriceMax(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{PriceMax: 150}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{PriceMax: 150}) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if len(p.Products) == 0 {
		t.Fatalf(`len(p.Products) = %d, want > 0`, len(p.Products))
	}
}

func TestSearchReturnsEmptySliceWhenPriceMaxIsOutOfRange(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{PriceMax: 0.5}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{PriceMax: 0.5}) = %v, want nil`, err.Error())
	}

	if p.Total != 0 {
		t.Fatalf(`p.Total = %d, want = 0`, p.Total)
	}

	if len(p.Products) != 0 {
		t.Fatalf(`len(p.Products) = %d, want = 0`, len(p.Products))
	}
}

func TestSearchReturnsProductWhenSlugIsFound(t *testing.T) {
	c := tests.Context()
	slug := stringutil.Slugify("T-shirt Tester c'est douter")
	p, err := Search(c, Query{Slug: slug}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Slug: slug}) = %v, want nil`, err.Error())
	}

	if p.Total != 0 {
		t.Fatalf(`p.Total = %d, want = 0`, p.Total)
	}

	if len(p.Products) != 0 {
		t.Fatalf(`len(p.Products) = %d, want = 0`, len(p.Products))
	}
}

func TestSearchReturnsProductsWhenTagsAreFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Tags: []string{"clothes"}}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Tags: []string{"clothes"}}) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if len(p.Products) == 0 {
		t.Fatalf(`len(p.Products) = %d, want > 0`, len(p.Products))
	}

	if p.Products[0].ID == "" {
		t.Fatalf(`p.Products[0].ID = %s, want not empty`, p.Products[0].ID)
	}
}

func TestSearchReturnsEmptySliceWhenTagsAreNotFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Tags: []string{"crazy"}}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Tags: []string{"crazy"}}) = %v, want nil`, err.Error())
	}

	if p.Total != 0 {
		t.Fatalf(`p.Total = %d, want = 0`, p.Total)
	}

	if len(p.Products) != 0 {
		t.Fatalf(`len(p.Products) = %d, want = 0`, len(p.Products))
	}
}

func TestSearchReturnsProductsWhenMetaAreFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Meta: map[string][]string{"color": {"blue"}}}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Meta: map[string]string{"color": "blue"}}) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if len(p.Products) == 0 {
		t.Fatalf(`len(p.Products) = %d, want > 0`, len(p.Products))
	}

	if p.Products[0].ID == "" {
		t.Fatalf(`p.Products[0].ID = %s, want ""`, p.Products[0].ID)
	}
}

func TestSearchReturnsProductsWhenMetaWithSpaceAreFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Meta: map[string][]string{"color": {"blue cyan"}}}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Meta: map[string]string{"color": "blue cyan"}}) = %v, want nil`, err.Error())
	}

	if p.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, p.Total)
	}

	if len(p.Products) == 0 {
		t.Fatalf(`len(p.Products) = %d, want > 0`, len(p.Products))
	}

	if p.Products[0].ID == "" {
		t.Fatalf(`p.Products[0].ID = %s, want ""`, p.Products[0].ID)
	}
}

func TestSearchReturnsEmptySliceWhenMetaAreNotFound(t *testing.T) {
	c := tests.Context()
	p, err := Search(c, Query{Meta: map[string][]string{"color": {"crazy"}}}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Meta: map[string]string{"color": "crazy"}}) = %v, want nil`, err.Error())
	}

	if p.Total != 0 {
		t.Fatalf(`p.Total = %d, want = 0`, p.Total)
	}

	if len(p.Products) != 0 {
		t.Fatalf(`len(p.Products) = %d, want = 0`, len(p.Products))
	}
}

func TestURLReturnsTheProductURLWhenSuccess(t *testing.T) {
	if product.URL() != "http://localhost/123-title.html" {
		t.Fatalf(`product.URL()  = %s, want 'http://localhost/123-title.html'`, product.URL())
	}
}

func TestDeleteReturnsErrorWhenIdIsEmpty(t *testing.T) {
	c := tests.Context()
	if err := Delete(c, ""); err == nil || err.Error() != "input:id" {
		t.Fatalf(`Delete(c, "") = %s want "input:id"`, err.Error())
	}
}

func TestListReturnsProductsWhenOk(t *testing.T) {
	c := tests.Context()
	pds, err := List(c, []string{"PDT1"})

	if err != nil {
		t.Fatalf(`List(c, []string{"PDT1"}) = %v, %v want not empty, nil`, pds, err.Error())
	}

	if len(pds) == 0 {
		t.Fatalf(`len(pds) = %d want > 0`, len(pds))
	}

	if pds[0].ID == "" {
		t.Fatalf(`pds[0].ID = %s want not empty`, pds[0].ID)
	}
}
