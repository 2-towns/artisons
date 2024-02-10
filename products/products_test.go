package products

import (
	"artisons/conf"
	"artisons/tests"
	"errors"
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var cur string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	cur = path.Dir(filename) + "/"
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

func TestImagePath(t *testing.T) {
	pid := "/products/123.jpeg"
	p := ImagePath(pid)
	expected := conf.WorkingSpace + "web/images/products/123.jpeg"

	if !strings.HasSuffix(p, expected) {
		t.Fatalf(`path = %s, want %s`, p, expected)
	}
}

func TestAvailable(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/products.redis")

	var tests = []struct {
		name      string
		pid       string
		available bool
	}{
		{"pid=PDT1", "PDT1", true},
		{"pid=idontexist", "idontexist", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if available := Available(ctx, tt.pid); available != tt.available {
				t.Fatalf(`available = %v, want %v`, available, tt.available)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	ctx := tests.Context()

	var tests = []struct {
		name  string
		field string
		value string
		err   error
	}{
		{"sku=!!!", "Sku", "!!!", errors.New("input:sku")},
		{"title=", "Title", "", errors.New("input:title")},
		{"description=", "Description", "", errors.New("input:description")},
		{"status=", "Status", "", errors.New("input:status")},
		{"status=abc", "Status", "abc", errors.New("input:status")},
		{"slug=", "Slug", "", errors.New("input:slug")},
		{"success", "", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := product

			if tt.field != "" {
				reflect.ValueOf(&p).Elem().FieldByName(tt.field).SetString(tt.value)
			}

			if err := p.Validate(ctx); fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.err) {
				t.Fatalf(`err = %v, want %s`, err, tt.err)
			}
		})
	}
}

func TestSave(t *testing.T) {
	c := tests.Context()

	if _, err := product.Save(c); err != nil {
		t.Fatalf(`err = %v, want nil`, err.Error())
	}
}

func TestFind(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/products.redis")

	var tests = []struct {
		name string
		pid  string
		err  error
	}{
		{"pid=''", "", errors.New("oops the data is not found")},
		{"pid='PDT1'", "PDT1", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Find(ctx, tt.pid); fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.err) {
				t.Fatalf(`err = %v, want %v`, err, tt.err)
			}
		})
	}
}

func TestFindAll(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/products.redis")

	var tests = []struct {
		name  string
		pids  []string
		count int
	}{
		{"pid=''", []string{""}, 0},
		{"pid='PDT1'", []string{"PDT1"}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if p, err := FindAll(ctx, tt.pids); err != nil || len(p) != tt.count {
				t.Fatalf(`count = %d, err = %v, want %d, nil`, len(p), err, tt.count)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	ctx := tests.Context()

	tests.Del(ctx, "product")
	tests.ImportData(ctx, cur+"testdata/products.redis")

	var tests = []struct {
		name  string
		q     Query
		count int
	}{
		{"keywords=T-Shirt", Query{Keywords: "t-shirt"}, 1},
		{"keywords=c'est", Query{Keywords: "c'est"}, 1},
		{"keywords=hello douter", Query{Keywords: "hello douter"}, 1},
		{"keywords=unisexe", Query{Keywords: "unisexe"}, 1},
		{"keywords=SKU1", Query{Keywords: "SKU1"}, 1},
		{"keywords=idontexist", Query{Keywords: "idontexist"}, 0},
		{"min=10", Query{PriceMin: 200}, 1},
		{"min=500", Query{PriceMin: 500}, 0},
		{"min=150", Query{PriceMax: 150}, 1},
		{"min=50", Query{PriceMax: 50}, 0},
		{"slug=t-shirt-tester-c-est-douter", Query{Slug: "t-shirt-tester-c-est-douter"}, 1},
		{"tags=clothes", Query{Tags: []string{"clothes"}}, 1},
		{"color=blue", Query{Meta: map[string][]string{"color": {"blue"}}}, 1},
		{"color=blue cyan", Query{Meta: map[string][]string{"color": {"blue cyan"}}}, 1},
		{"color=idontexist", Query{Meta: map[string][]string{"color": {"idontexist"}}}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := Search(ctx, tt.q, 0, 10)

			if err != nil {
				t.Fatalf(`err = %v, want nil`, err.Error())
			}

			if p.Total != tt.count {
				t.Fatalf(`total = %d, want %d`, p.Total, tt.count)
			}

			if len(p.Products) != tt.count {
				t.Fatalf(`len(products) = %d, want %d`, len(p.Products), tt.count)
			}
		})
	}
}

func TestURL(t *testing.T) {
	if product.URL() != "http://localhost/123-title.html" {
		t.Fatalf(`url  = %s, want 'http://localhost/123-title.html'`, product.URL())
	}
}

func TestDelete(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/products.redis")

	if err := Delete(ctx, ""); err == nil || err.Error() != "input:id" {
		t.Fatalf(`err = %s want "input:id"`, err.Error())
	}
}

func TestList(t *testing.T) {
	c := tests.Context()
	pds, err := List(c, []string{"PDT1"})

	if err != nil {
		t.Fatalf(`err = %v want not nil`, err.Error())
	}

	if len(pds) == 0 {
		t.Fatalf(`len(pds) = %d want > 0`, len(pds))
	}

	if pds[0].ID == "" {
		t.Fatalf(`id = %s want not empty`, pds[0].ID)
	}
}
