package stringutil

import "testing"

// TestRandom expects to succeed
func TestRandom(t *testing.T) {
	r, err := Random()

	if err != nil || r == "" {
		t.Fatalf(`Random() = %s, %v, want string, nil`, r, err)
	}
}

// TestSlugify expects to succeed
func TestSlugify(t *testing.T) {
	s := Slugify("VERy nice title 12")

	if s != "very-nice-title-12" {
		t.Fatalf(`Slugify("VERy nice title 12") = %s, want "very-nice-title-12"`, s)
	}
}
