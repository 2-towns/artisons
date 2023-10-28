package users

import "testing"

// TestFindBySessionID expects to succeed
func TestFindBySessionID(t *testing.T) {
	alive := true
	user := createUser(alive)
	u, err := findBySessionID(user.SID)

	if err != nil || u.SID == "" || u.Email == "" {
		t.Fatalf("findBySessionID(user.SID) = %v, %v, want User, nil", u, err)
	}
}

// TestFindBySessionIDWithoutSID expects to fail because of sid emptyness
func TestFindBySessionIDWithoutSID(t *testing.T) {
	u, err := findBySessionID("")

	if err == nil || err.Error() != "unauthorized" || u.Email != "" {
		t.Fatalf("findBySessionID('') = %v, %v, want User{}, 'unauthorized'", u, err)
	}
}

// TestFindBySessionIDExpired expects to fail because of session expired
func TestFindBySessionIDExpired(t *testing.T) {
	alive := false
	user := createUser(alive)
	u, err := findBySessionID(user.SID)

	if err == nil || err.Error() != "unauthorized" || u.Email != "" {
		t.Fatalf("findBySessionID(user.SID) = %v, %v, want User, nil", u, err)
	}
}
