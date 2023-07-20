package users

import (
	"context"
	"fmt"
	"gifthub/db"
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

	if users[0].ID != 1 {
		t.Fatalf(`The first ID is wrong %d`, users[0].ID)
	}
}

// TestUserPersist get the user list from redis
func TestUserPersist(t *testing.T) {
	u := User{
		Email:    faker.Email(),
		Username: strings.ToLower(faker.Username()),
	}

	id, err := u.Persist("passw0rd")
	if id == 0 {
		t.Fatalf(`the user id should not be 0`)
	}
	if err != nil {
		t.Fatalf(`the persist failed because of %s`, err.Error())
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
	if id != 0 {
		t.Fatalf(`the user id should be 0`)
	}

	if err == nil {
		t.Fatalf(`the persist should fail because the username is empty`)
	}

	if err.Error() != "user_email_invalid" {
		t.Fatalf(`the error message is incorrect`)
	}
}

// TestUserPersistFailedWithBadEmail fails when email is incorrect
func TestUserPersistFailedWithBadEmail(t *testing.T) {
	u := User{
		Email:    faker.Username(),
		Username: strings.ToLower(faker.Username()),
	}

	id, err := u.Persist("passw0rd")
	if id != 0 {
		t.Fatalf(`the user id should be 0`)
	}
	if err == nil {
		t.Fatalf(`the persist should fail because the email is incorrect`)
	}
	if err.Error() != "user_email_invalid" {
		t.Fatalf(`the error message is incorrect`)
	}
}

// TestUserPersistFailedWithEmptyPassword fails when password is empty
func TestUserPersistFailedWithEmptyPassword(t *testing.T) {
	u := User{
		Email:    faker.Username(),
		Username: strings.ToLower(faker.Username()),
	}

	id, err := u.Persist("")
	if id != 0 {
		t.Fatalf(`the user id should be 0`)
	}
	if err == nil {
		t.Fatalf(`the persist should fail because the email is incorrect`)
	}
	if err.Error() != "user_password_required" {
		t.Fatalf(`the error message is incorrect`)
	}
}

// TestUserPersistFailedWithExistingUsername fails when the username already exists
func TestUserPersistFailedWithExistingUsername(t *testing.T) {
	username := strings.ToLower(faker.Username())

	u := User{
		Email:    faker.Email(),
		Username: username,
	}

	_, err := u.Persist("passw0rd")
	if err != nil {
		t.Fatalf(`the first persist should work but got %s`, err.Error())
	}

	id, err := u.Persist("passw0rd")
	if id != 0 {
		t.Fatalf(`the user id should be 0`)
	}
	if err == nil {
		t.Fatalf(`the persist should fail because the username is already added`)
	}
	if err.Error() != "user_username_exists" {
		t.Fatalf(`the error message is incorrect`)
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
	if id != 0 {
		t.Fatalf(`the user id should be 0`)
	}
	if err == nil {
		t.Fatalf(`the persist should fail because the username is incorrect`)
	}
	if err.Error() != "user_username_invalid" {
		t.Fatalf(`the error message is incorrect`)
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
		t.Fatalf(`the delete should not fail`)
	}

	ctx := context.Background()

	id, err := db.Redis.HGet(ctx, fmt.Sprintf("user:%d", u.ID), "id").Result()
	if id != "" {
		t.Fatalf(`the id should be empty`)
	}
	if err == nil {
		t.Fatalf(`the id hget should fail`)
	}

	id, err = db.Redis.HGet(ctx, "user", u.Username).Result()
	if id != "" {
		t.Fatalf(`the id should be empty`)
	}
	if err == nil {
		t.Fatalf(`the username hget should fail`)
	}

	result, _ := db.Redis.SIsMember(ctx, "users", u.ID).Result()
	if result == true {
		t.Fatalf(`the member should not be a part of users`)
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
		t.Fatalf(`the delete should not fail`)
	}
}
