package users

import (
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
)

// TestUserList get the user list from redis
func TestUserList(t *testing.T) {
	users, err := List(0)

	if err != nil {
		t.Fatal(`The user list should not failed.`)
	}

	if len(users) == 0 {
		t.Fatal(`The users should not be empty`)
	}
}

// TestUserPersist get the user list from redis
func TestUserPersist(t *testing.T) {
	u := User{
		Email:    faker.Email(),
		Username: strings.ToLower(faker.Username()),
	}

	err := u.Persist("passw0rd")

	if err != nil {
		t.Fatalf(`the persist failed because of %s`, err.Error())
	}
}

// TestUserPersistFailedWithUsernameUpper fails when username has uppercase
func TestUserPersistFailedWithUsernameUpper(t *testing.T) {
	u := User{
		Email:    faker.Email(),
		Username: strings.ToUpper(faker.Username()),
	}

	err := u.Persist("passw0rd")

	if err == nil {
		t.Fatalf(`the persist should failed because the username has uppercase`)
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
		t.Fatalf(`the persist should failed because the username is empty`)
	}
}

// TestUserPersistFailedWithEmailEmpty fails when email is empty
func TestUserPersistFailedWithEmailEmpty(t *testing.T) {
	u := User{
		Email:    "",
		Username: strings.ToLower(faker.Username()),
	}

	err := u.Persist("passw0rd")

	if err == nil {
		t.Fatalf(`the persist should failed because the username is empty`)
	}
}

// TestUserPersistFailedWithBadEmail fails when email is incorrect
func TestUserPersistFailedWithBadEmail(t *testing.T) {
	u := User{
		Email:    faker.Username(),
		Username: strings.ToLower(faker.Username()),
	}

	err := u.Persist("passw0rd")

	if err == nil {
		t.Fatalf(`the persist should failed because the email is incorrect`)
	}
}

// TestUserPersistFailedWithEmptyPassword fails when password is empty
func TestUserPersistFailedWithEmptyPassword(t *testing.T) {
	u := User{
		Email:    faker.Username(),
		Username: strings.ToLower(faker.Username()),
	}

	err := u.Persist("")

	if err == nil {
		t.Fatalf(`the persist should failed because the email is incorrect`)
	}
}

// TestUserPersistFailedWithExistingUsername fails when the username already exists
func TestUserPersistFailedWithExistingUsername(t *testing.T) {
	u := User{
		Email:    "toto",
		Username: strings.ToLower(faker.Username()),
	}

	err := u.Persist("passw0rd")

	if err == nil {
		t.Fatalf(`the persist should failed because the email is incorrect`)
	}
}
