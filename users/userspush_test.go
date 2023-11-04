package users

import (
	"gifthub/tests"
	"testing"
)

// TestAddWPToken expects to succeed
func TestAddWPToken(t *testing.T) {
	ctx := tests.Context()
	if err := user.AddWPToken(ctx, "{}"); err != nil {
		t.Fatalf("user.AddWPToken('{}') = %v, want nil", err)
	}
}

// TestAddWPTokenEmpty expects to fail because of token emptyness
func TestAddWPTokenEmpty(t *testing.T) {
	ctx := tests.Context()
	if err := user.AddWPToken(ctx, ""); err == nil || err.Error() != "user_wptoken_required" {
		t.Fatalf("user.AddWPToken('') = %v, want user_wptoken_required", err)
	}
}

// TestDeleteWPToken expects to succeed
func TestDeleteWPToken(t *testing.T) {
	ctx := tests.Context()
	if err := user.DeleteWPToken(ctx, user.SID); err != nil {
		t.Fatalf("user.TestDeleteWPToken(user.SID) = %v, want nil", err)
	}
}

// TestDeleteWPTokenEmpty expects to fail because of token emptyness
func TestDeleteWPTokenEmpty(t *testing.T) {
	ctx := tests.Context()
	if err := user.DeleteWPToken(ctx, ""); err == nil || err.Error() != "unauthorized" {
		t.Fatalf("u.DeleteWPToken('') = %v, want unauthorized", err)
	}
}
