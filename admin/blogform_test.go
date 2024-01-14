package admin

import (
	"gifthub/tests"
	"mime/multipart"
	"testing"
)

var article = map[string][]string{
	"title":       {"My  great article"},
	"description": {"My great description"},
	"status":      {"online"},
	"image":       {"1.jpeg"},
	"lang":        {"abc"},
}

func TestProcessBlogFormReturnsErrorWhenTitleInvalid(t *testing.T) {
	c := tests.Context()

	a := make(map[string][]string)
	for k, val := range article {
		a[k] = val
	}

	a["title"] = []string{""}

	f := multipart.Form{Value: a}

	if _, err := processBlogFrom(c, f, ""); err == nil || err.Error() != "input:title" {
		t.Fatalf(`processBlogFrom(c, f, "") = _, %v, want _, 'input:title'`, err.Error())
	}
}

func TestProcessBlogtFormReturnsErrorWhenDescriptionIsInvalid(t *testing.T) {
	c := tests.Context()

	a := make(map[string][]string)
	for k, val := range article {
		a[k] = val
	}

	a["description"] = []string{""}

	f := multipart.Form{Value: a}

	if _, err := processBlogFrom(c, f, ""); err == nil || err.Error() != "input:description" {
		t.Fatalf(`processBlogFrom(c, f, "") = _, %v, want _, 'input:description'`, err.Error())
	}
}

func TestProcessBlogtFormReturnsErrorWhenStatusIsInvalid(t *testing.T) {
	c := tests.Context()

	a := make(map[string][]string)
	for k, val := range article {
		a[k] = val
	}

	a["status"] = []string{""}

	f := multipart.Form{Value: a}

	if _, err := processBlogFrom(c, f, ""); err == nil || err.Error() != "input:status" {
		t.Fatalf(`processBlogFrom(c, f, "") = _, %v, want _, 'input:status'`, err.Error())
	}
}

func TestProcessBlogtFormReturnsErrorWhenLangIsEmpty(t *testing.T) {
	c := tests.Context()

	a := make(map[string][]string)
	for k, val := range article {
		a[k] = val
	}

	a["lang"] = []string{""}

	f := multipart.Form{Value: a}

	if _, err := processBlogFrom(c, f, ""); err == nil || err.Error() != "input:lang" {
		t.Fatalf(`processBlogFrom(c, f, "") = _, %v, want _, 'input:lang'`, err.Error())
	}
}

func TestProcessBlogtFormReturnsErrorWhenLangIsInvalid(t *testing.T) {
	c := tests.Context()

	a := make(map[string][]string)
	for k, val := range article {
		a[k] = val
	}

	a["lang"] = []string{"!!!"}

	f := multipart.Form{Value: a}

	if _, err := processBlogFrom(c, f, ""); err == nil || err.Error() != "input:lang" {
		t.Fatalf(`processBlogFrom(c, f, "") = _, %v, want _, 'input:lang'`, err.Error())
	}
}

func TestProcessBlogtFormReturnsErrorWhenImageIsInvalid(t *testing.T) {
	c := tests.Context()

	a := make(map[string][]string)
	for k, val := range article {
		a[k] = val
	}

	a["image"] = []string{""}

	f := multipart.Form{Value: a}

	if _, err := processBlogFrom(c, f, ""); err == nil || err.Error() != "input:image" {
		t.Fatalf(`processBlogFrom(c, f, "") = _, %v, want _, 'input:image'`, err.Error())
	}
}
