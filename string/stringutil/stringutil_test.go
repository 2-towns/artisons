package stringutil

import (
	"testing"
)

func TestRandomReturnsNumberWhenSuccess(t *testing.T) {
	r, err := Random()

	if err != nil || r == "" {
		t.Fatalf(`Random() = %s, %v, want string, nil`, r, err)
	}
}
