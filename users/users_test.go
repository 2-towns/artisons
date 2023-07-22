package users

import (
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
)

var TestUser = User{
	Email:    faker.Email(),
	Username: strings.ToLower(faker.Username()),
}

func init() {
	TestUser.Persist("passw0rd")
}

// TestUserList get the user list from redis
func TestUserList(t *testing.T) {
	users, err := List(0)

	if err != nil || len(users) == 0 || users[0].ID == 0 {
		t.Fatalf("List(0) = %v, %v, not want [{ID: 0}], error", users, err.Error())
	}
}

// TestUserPersist get the user list from redis
func TestUserPersist(t *testing.T) {
	u := User{
		Email:    faker.Email(),
		Username: strings.ToLower(faker.Username()),
	}

	id, err := u.Persist("passw0rd")
	if id == 0 || err != nil {
		t.Fatalf("Persist = %d, %v, not want 0, error", id, err.Error())
	}
}

/*
// TestUserPersistFailedWithUsernameUpper fails when username has uppercase
func TestUserPersistFailedWithUsernameUpper(t *testing.T) {
	u := User{
		Email:    faker.Email(),
		Username: strings.ToUpper(faker.Username()),
	}

	err := u.Persist("passw0rd")

	if err == nil {
		t.Fatalf(`the persist should fail because the username has uppercase`)
	}
}

// TestUserPersistFailedWithUsernameEmpty fails when username is empty
func TestUserPersistFailedWithUsernameEmpty(t *testing.T) {
	u := User{
		Email:    faker.Email(),
		Username: "",
	}

	err := u.Persist("passw0rd")

	if err == nil {
		t.Fatalf(`the persist should fail because the username is empty`)
	}
}
*/
// TestUserPersistFailedWithEmailEmpty fails when email is empty
func TestUserPersistFailedWithEmailEmpty(t *testing.T) {
	u := User{
		Email:    "",
		Username: strings.ToLower(faker.Username()),
	}

	id, err := u.Persist("passw0rd")
	if id != 0 || err == nil || err.Error() != "user_email_invalid" {
		t.Fatalf("Persist = %d, %v, want 0, error", id, err.Error())
	}
}

// TestUserPersistFailedWithBadEmail fails when email is incorrect
func TestUserPersistFailedWithBadEmail(t *testing.T) {
	u := User{
		Email:    faker.Username(),
		Username: strings.ToLower(faker.Username()),
	}

	id, err := u.Persist("passw0rd")
	if id != 0 || err == nil || err.Error() != "user_email_invalid" {
		t.Fatalf("Persist = %d, %v, want 0, error", id, err.Error())
	}
}

// TestUserPersistFailedWithEmptyPassword fails when password is empty
func TestUserPersistFailedWithEmptyPassword(t *testing.T) {
	u := User{
		Email:    faker.Username(),
		Username: strings.ToLower(faker.Username()),
	}

	id, err := u.Persist("")
	if id != 0 || err == nil || err.Error() != "user_password_required" {
		t.Fatalf("Persist = %d, %v, want 0, error", id, err.Error())
	}
}

// TestUserPersistFailedWithExistingUsername fails when the username already exists
func TestUserPersistFailedWithExistingUsername(t *testing.T) {
	id, err := TestUser.Persist("passw0rd")
	if id != 0 || err == nil || err.Error() != "user_username_exists" {
		t.Fatalf("Persist = %d, %v, want 0, error", id, err.Error())
	}
}

// TestUserPersistFailedWithBadUsername fails when the username is incorrect
func TestUserPersistFailedWithBadUsername(t *testing.T) {
	username := faker.Username()

	u := User{
		Email:    faker.Email(),
		Username: username,
	}

	id, err := u.Persist("passw0rd")
	if id != 0 || err == nil || err.Error() != "user_username_invalid" {
		t.Fatalf("Persist = %d, %v, want 0, error", id, err.Error())
	}
}

// TestDeleteUser deletes an existing user
func TestDeleteUser(t *testing.T) {
	u := User{
		Email:    faker.Email(),
		Username: strings.ToLower(faker.Username()),
	}

	u.Persist("passw0rd")

	err := u.Delete()
	if err != nil {
		t.Fatalf("Delete(), %v, want nil, error", err.Error())
	}
}

// TestDeleteUserNotExisting does nothing
func TestDeleteUserNotExisting(t *testing.T) {
	u := User{
		Email:    faker.Email(),
		Username: strings.ToLower(faker.Username()),
	}

	err := u.Delete()
	if err != nil {
		t.Fatalf("Delete(), %v, want nil, error", err.Error())
	}
}
