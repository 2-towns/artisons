package users

import (
	"gifthub/tests"
	"testing"
)

func TestFindBySessionIDReturnsSessionWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	u, err := findBySessionID(ctx, "test")

	if err != nil || u.SID == "" || u.Email == "" {
		t.Fatalf(`findBySessionID("test) = %v, %v, want User, nil`, u, err)
	}
}

func TestFindBySessionIDReturnsErrorWhenSidIsEmpty(t *testing.T) {
	ctx := tests.Context()
	u, err := findBySessionID(ctx, "")

	if err == nil || err.Error() != "unauthorized" || u.Email != "" {
		t.Fatalf("findBySessionID('') = %v, %v, want User{}, 'unauthorized'", u, err)
	}
}

func TestFindBySessionIDReturnsErrorWhenSessionIsExpired(t *testing.T) {
	ctx := tests.Context()
	u, err := findBySessionID(ctx, "expired")

	if err == nil || err.Error() != "unauthorized" || u.Email != "" {
		t.Fatalf(`findBySessionID("expired") = %v, %v, want User, nil`, u, err)
	}
}
