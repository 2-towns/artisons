package users

import "testing"

// TestAddWPToken expects to succeed
func TestAddWPToken(t *testing.T) {
	alive := true
	u := createUser(alive)
	if err := u.AddWPToken("{}"); err != nil {
		t.Fatalf("u.AddWPToken('{}') = %v, want nil", err)
	}
}

// TestAddWPTokenEmpty expects to fail because of token emptyness
func TestAddWPTokenEmpty(t *testing.T) {
	alive := true
	u := createUser(alive)
	if err := u.AddWPToken(""); err == nil || err.Error() != "user_wptoken_required" {
		t.Fatalf("u.AddWPToken('') = %v, want user_wptoken_required", err)
	}
}

// TestDeleteWPToken expects to succeed
func TestDeleteWPToken(t *testing.T) {
	alive := true
	u := createUser(alive)
	if err := u.DeleteWPToken(u.SID); err != nil {
		t.Fatalf("u.TestDeleteWPToken(u.SID) = %v, want nil", err)
	}
}

// TestDeleteWPTokenEmpty expects to fail because of token emptyness
func TestDeleteWPTokenEmpty(t *testing.T) {
	alive := true
	u := createUser(alive)
	if err := u.DeleteWPToken(""); err == nil || err.Error() != "unauthorized" {
		t.Fatalf("u.DeleteWPToken('') = %v, want unauthorized", err)
	}
}
