package stringutil

import "testing"

// TestRandom checks the generation of a random string
func TestRandom(t *testing.T) {
	r, err := Random()

	if err != nil || r == "" {
		t.Fatalf(`Random() = %s, %v, not want "", error`, r, err)
	}
}

// TestSlugify checks the slug result
func TestSlugify(t *testing.T) {
	s := Slugify("VERy nice title 12")

	if s != "very-nice-title-12" {
		t.Fatalf(`Slugify("VERy nice title 12") = %s, want "very-nice-title-12", error`, s)
	}
}