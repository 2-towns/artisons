package products

import (
	"context"
	"gifthub/db"
	"gifthub/string/stringutil"
	"gifthub/tests"
	"testing"
)

// TestImagePath expects to succeed
func TestImagePath(t *testing.T) {
	pid := "123"
	index := 1
	_, p := ImagePath(pid, index)
	expected := "../web/images/123/1"

	if p != expected {
		t.Fatalf(`TestImagePath("123", 1) = %s, want %s`, p, expected)
	}
}

// TestProductAvailable expects to succeed when the product exists
func TestProductAvailable(t *testing.T) {
	ctx := context.Background()
	pid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "product:"+pid, "status", "online")
	c := tests.Context()

	if exists := Available(c, pid); !exists {
		t.Fatalf(`Available(pid) = %v, want true`, exists)
	}
}

// TestProductAvailableNotFound expects to fail because of product non existence
func TestProductAvailableNotFound(t *testing.T) {
	c := tests.Context()

	if exists := Available(c, "toto"); exists {
		t.Fatalf(`Available(c, pid) = %v, want false`, exists)
	}
}

// TestProductsAvailables expects to succeed
func TestProductsAvailables(t *testing.T) {
	ctx := context.Background()
	pid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "product:"+pid, "status", "online")
	c := tests.Context()

	if exists := Availables(c, []string{pid}); !exists {
		t.Fatalf(`Availables(c, pid) = %v, want true`, exists)
	}
}

// TestProductsAvailablesNotFound expects to fail because of products non existence
func TestProductsAvailablesNotFound(t *testing.T) {
	c := tests.Context()

	if exists := Availables(c, []string{"toto"}); exists {
		t.Fatalf(`Availables(c, pid) = %v, want false`, exists)
	}
}
