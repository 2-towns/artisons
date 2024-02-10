package stringutil

import (
	"testing"
)

func TestRandom(t *testing.T) {
	r, err := Random()

	if err != nil || r == "" {
		t.Fatalf(`random, err = %s, %v, want string, nil`, r, err)
	}
}
