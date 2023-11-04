package users

import (
	"gifthub/tests"
	"testing"
)

// TestFindBySessionID expects to succeed
func TestFindBySessionID(t *testing.T) {
	ctx := tests.Context()
	u, err := findBySessionID(ctx, "test")

	if err != nil || u.SID == "" || u.Email == "" {
		t.Fatalf(`findBySessionID("test) = %v, %v, want User, nil`, u, err)
	}
}

// TestFindBySessionIDWithoutSID expects to fail because of sid emptyness
func TestFindBySessionIDWithoutSID(t *testing.T) {
	ctx := tests.Context()
	u, err := findBySessionID(ctx, "")

	if err == nil || err.Error() != "unauthorized" || u.Email != "" {
		t.Fatalf("findBySessionID('') = %v, %v, want User{}, 'unauthorized'", u, err)
	}
}

// TestFindBySessionIDExpired expects to fail because of session expired
func TestFindBySessionIDExpired(t *testing.T) {
	ctx := tests.Context()
	u, err := findBySessionID(ctx, "expired")

	if err == nil || err.Error() != "unauthorized" || u.Email != "" {
		t.Fatalf(`findBySessionID("expired") = %v, %v, want User, nil`, u, err)
	}
}
