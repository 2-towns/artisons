package stringutil

import (
	"testing"
)

// TestRandom expects to succeed
func TestRandom(t *testing.T) {
	r, err := Random()

	if err != nil || r == "" {
		t.Fatalf(`Random() = %s, %v, want string, nil`, r, err)
	}
}
