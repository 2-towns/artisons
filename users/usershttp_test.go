package users

import "testing"

// TestFindBySessionID expects to succeed
func TestFindBySessionID(t *testing.T) {
	alive := true
	_, sid := createUser(alive)
	u, err := findBySessionID(sid)

	if err != nil || u.Email == "" {
		t.Fatalf("findBySessionID(sid) = %v, %v, want User, nil", u, err)
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
	_, sid := createUser(alive)
	u, err := findBySessionID(sid)

	if err == nil || err.Error() != "unauthorized" || u.Email != "" {
		t.Fatalf("findBySessionID(sid) = %v, %v, want User, nil", u, err)
	}
}
