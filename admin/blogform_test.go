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
}

func TestProcessBlogFormReturnsErrorWhenTitleInvalid(t *testing.T) {
	c := tests.Context()

	a := make(map[string][]string)
	for k, val := range article {
		a[k] = val
	}

	a["title"] = []string{""}

	f := multipart.Form{Value: a}

	if _, err := processBlogFrom(c, f, ""); err == nil || err.Error() != "input_title_invalid" {
		t.Fatalf(`processBlogFrom(c, f, "") = _, %v, want _, 'input_title_invalid'`, err.Error())
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

	if _, err := processBlogFrom(c, f, ""); err == nil || err.Error() != "input_description_invalid" {
		t.Fatalf(`processBlogFrom(c, f, "") = _, %v, want _, 'input_description_invalid'`, err.Error())
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

	if _, err := processBlogFrom(c, f, ""); err == nil || err.Error() != "input_status_invalid" {
		t.Fatalf(`processBlogFrom(c, f, "") = _, %v, want _, 'input_status_invalid'`, err.Error())
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

	if _, err := processBlogFrom(c, f, ""); err == nil || err.Error() != "input_image_required" {
		t.Fatalf(`processBlogFrom(c, f, "") = _, %v, want _, 'input_image_required'`, err.Error())
	}
}
