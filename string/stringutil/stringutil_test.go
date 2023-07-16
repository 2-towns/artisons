package stringutil

import "testing"

// TestRandom checks the generation of a random string
func TestRandom(t *testing.T) {
	r, err := Random()

	if err != nil {
		t.Fatalf(`the random should not failed: %s `, err.Error())
	}

	if r == "" {
		t.Fatal(`the random should not be empty`)
	}
}

// TestSlugify checks the slug result
func TestSlugify(t *testing.T) {
	s := Slugify("VERy nice title 12")

	if s != "very-nice-title-12" {
		t.Fatal(`the slug is not correct`)
	}
}
