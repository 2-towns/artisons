package users

import (
	"gifthub/tests"
	"testing"
)

// TestFindBySessionID expects to succeed
func TestFindBySessionID(t *testing.T) {
	alive := true
	user := createUser(alive)
	ctx := tests.Context()
	u, err := findBySessionID(ctx, user.SID)

	if err != nil || u.SID == "" || u.Email == "" {
		t.Fatalf("findBySessionID(user.SID) = %v, %v, want User, nil", u, err)
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
	alive := false
	user := createUser(alive)
	ctx := tests.Context()
	u, err := findBySessionID(ctx, user.SID)

	if err == nil || err.Error() != "unauthorized" || u.Email != "" {
		t.Fatalf("findBySessionID(user.SID) = %v, %v, want User, nil", u, err)
	}
}
