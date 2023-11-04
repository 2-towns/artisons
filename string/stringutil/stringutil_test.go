package stringutil

import (
	"fmt"
	"testing"
)

// TestRandom expects to succeed
func TestRandom(t *testing.T) {
	r, err := Random()

	if err != nil || r == "" {
		t.Fatalf(`Random() = %s, %v, want string, nil`, r, err)
	}
}

func ExampleSlugify(t *testing.T) {
	fmt.Println("VERy nice title 12")
	// Output: very-nice-title-12
}
