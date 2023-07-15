package users

import (
	"testing"
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
