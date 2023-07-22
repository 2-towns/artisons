package users

import (
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
)

var AuthUser = User{
	Email:    faker.Email(),
	Username: strings.ToLower(faker.Username()),
}

func init() {
	AuthUser.Persist("passw0rd")
	Login(AuthUser.Username, "passw0rd")
}

// TestLogin makes sure than the TestUser can login
func TestLogin(t *testing.T) {
	sid, err := Login(TestUser.Username, "passw0rd")
	if err != nil || sid == "" {
		t.Fatalf(`TestLogin(TestUser.Username, "passw0rd")) = %s, %v, not want sid == "" , error`, sid, err)
	}
}

// TestLoginFailsWithBadUsername fails when the username is invalid
func TestLoginFailsWithBadUsername(t *testing.T) {
	if sid, err := Login("I do not Exist", "passw0rd"); err == nil || err.Error() != "user_login_failed" {
		t.Fatalf("Login('I do not Exist', 'passw0rd') = %s, %v, want 'user_login_failed', error", sid, err)
	}
}

// TestLoginFailsWithNotExistingUsername fails when the username does not exist
func TestLoginFailsWithNotExistingUsername(t *testing.T) {
	if sid, err := Login("idonotexist", "passw0rd"); err == nil || err.Error() != "user_login_failed" {
		t.Fatalf("Login('idonotexist', 'passw0rd') = %s, %v, want 'user_login_failed', error", sid, err)
	}
}

// TestLoginFailsWithBadPassword fails when the password does not match
func TestLoginFailsWithBadPassword(t *testing.T) {
	if sid, err := Login(TestUser.Username, "bad"); err == nil || err.Error() != "user_login_failed" {
		t.Fatalf("Login(TestUser.Username, 'bad') = %s, %v, want 'user_login_failed', error", sid, err)
	}
}

// TestLogout makes sure than an user is logout
func TestLogout(t *testing.T) {
	if err := AuthUser.Logout(); err != nil {
		t.Fatalf("AuthUser.Logout(), %v, want '', error", err)
	}
}
