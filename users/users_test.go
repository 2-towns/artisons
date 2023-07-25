package users

import (
	"gifthub/string/stringutil"
	"testing"

	"github.com/go-faker/faker/v4"
)

// TestUserList get the user list from redis
func TestUserList(t *testing.T) {
	users, err := List(0)

	if err != nil || len(users) == 0 || users[0].ID == 0 {
		t.Fatalf("List(0) = %v, %v, not want [{ID: 0}], error", users, err)
	}
}

// TestMagicCode generate a magic code for an email
func TestMagicCode(t *testing.T) {
	magic, err := MagicCode(faker.Email())
	if magic == "" || err != nil {
		t.Fatalf("TestMagicCode(faker.Email()) = %s, %v, not want '', error", magic, err)
	}
}

// TestMagicCodeTwice generate a magic code for an email
func TestMagicCodeTwice(t *testing.T) {
	magic, _ := MagicCode(faker.Email())
	magic, err := MagicCode(faker.Email())
	if magic == "" || err != nil {
		t.Fatalf("TestMagicCode(faker.Email()) = %s, %v, not want '', error", magic, err)
	}
}

// TestMagicCodeWithEmailEmpty fails when email is empty
func TestMagicCodeWithEmailEmpty(t *testing.T) {
	magic, err := MagicCode("")
	if magic != "" || err == nil {
		t.Fatalf("TestMagicCode('') = %s, %v, want '', error", magic, err)
	}
}

// TestMagicCodeFailedWithBadEmail fails when email is incorrect
func TestMagicCodeFailedWithBadEmail(t *testing.T) {
	magic, err := MagicCode("toto")
	if magic != "" || err == nil {
		t.Fatalf("TestMagicCode('toto') = %s, %v, want '', error", magic, err)
	}
}

// TestDeleteUser deletes an existing user
func TestDeleteUser(t *testing.T) {
	email := faker.Email()

	id, _ := saveUser(email, "toto")
	sid, _ := saveSID(id)
	u, _ := findBySessionID(sid)

	err := u.Delete()
	if err != nil {
		t.Fatalf("Delete(), %v, want nil, error", err)
	}
}

// TestDeleteUserNotExisting does nothing
func TestDeleteUserNotExisting(t *testing.T) {
	err := User{}.Delete()
	if err != nil {
		t.Fatalf("Delete(), %v, want nil, error", err)
	}
}

// Testlogin authenticates an user
func TestLogin(t *testing.T) {
	magic, _ := stringutil.Random()
	saveUser(faker.Email(), magic)
	sid, err := Login(magic)

	if sid == "" || err != nil {
		t.Fatalf("TestLogin(magic) = %s, %v, not want '', error", sid, err)
	}
}

// TestLoginWithEmptyMagic try to authenticate an user with empty magic
func TestLoginWithEmptyMagic(t *testing.T) {
	sid, err := Login("")
	if sid != "" || err == nil {
		t.Fatalf("TestLogin('') = %s, %v, want '', error", sid, err)
	}
}

// TestLoginWithNotExistingMagic try to authenticate an user with not existing magic
func TestLoginWithNotExistingMagic(t *testing.T) {
	sid, err := Login("titi")
	if sid != "" || err == nil {
		t.Fatalf("TestLogin('') = %s, %v, want '', error", sid, err)
	}
}

// TestLogout logout an user
func TestLogout(t *testing.T) {
	email := faker.Email()

	id, _ := saveUser(email, "toto")
	sid, _ := saveSID(id)

	err := Logout(id, sid)
	if err != nil {
		t.Fatalf("Logout(id, sid), %v, want nil, error", err)
	}
}

// TestLogoutWithZeroID returns an error
func TestLogoutWithZeroID(t *testing.T) {
	err := Logout(0, "123")
	if err == nil {
		t.Fatalf("Logout(0, '123'), %v, not want nil, error", err)
	}
}

// TestLogoutWithEmptySID returns an error
func TestLogoutWithEmptySID(t *testing.T) {
	err := Logout(124, "")
	if err == nil {
		t.Fatalf("Logout(124, ''), %v, not want nil, error", err)
	}
}

// TestLogoutWithNotExistingData does nothing
func TestLogoutWithNotExistingData(t *testing.T) {
	err := Logout(124, "122")
	if err != nil {
		t.Fatalf("Logout(124, '122'), %v, not want nil, error", err)
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
