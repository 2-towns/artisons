package users

import (
	"context"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/string/stringutil"
	"math/rand"
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
	sid, _ := stringutil.Random()
	id := rand.Intn(10000)
	ctx := context.Background()

	db.Redis.Set(ctx, "auth:"+sid, id, 0)
	db.Redis.HSet(ctx, fmt.Sprintf("user:%d", id), "email", email)

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
	id := rand.Intn(10000)
	ctx := context.Background()
	db.Redis.Set(ctx, "magic:"+magic, id, 0)

	sid, err := Login(magic, "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid == "" || err != nil {
		t.Fatalf("TestLogin(magic) = %s, %v, not want '', error", sid, err)
	}
}

// TestLoginWithoutDevice fails when device is empty
func TestLoginWithoutDevice(t *testing.T) {
	magic, _ := stringutil.Random()
	id := rand.Intn(10000)
	ctx := context.Background()
	db.Redis.Set(ctx, "magic:"+magic, id, 0)

	sid, err := Login(magic, "")
	if sid != "" || err == nil {
		t.Fatalf("TestLogin(magic) = %s, %v, want '', error", sid, err)
	}
}

// TestLoginWithEmptyMagic try to authenticate an user with empty magic
func TestLoginWithEmptyMagic(t *testing.T) {
	sid, err := Login("", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil {
		t.Fatalf("TestLogin('') = %s, %v, want '', error", sid, err)
	}
}

// TestLoginWithNotExistingMagic try to authenticate an user with not existing magic
func TestLoginWithNotExistingMagic(t *testing.T) {
	sid, err := Login("titi", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil {
		t.Fatalf("TestLogin('') = %s, %v, want '', error", sid, err)
	}
}

// TestLogout logout an user
func TestLogout(t *testing.T) {
	ctx := context.Background()
	sid, _ := stringutil.Random()
	id := rand.Int63n(10000)
	db.Redis.Set(ctx, "auth:"+sid, id, 0)
	db.Redis.HSet(ctx, "session:"+sid, "id", id)

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

// TestSessions returns the user sessions
func TestSessions(t *testing.T) {
	sid, _ := stringutil.Random()
	id := rand.Int63n(10000)
	u := User{
		ID: id,
		Devices: map[string]string{
			"auth:" + sid: "Mozilla/5.0 Gecko/20100101 Firefox/115.0",
		},
	}

	ctx := context.Background()
	db.Redis.Set(ctx, "auth:"+sid, id, conf.SessionDuration)

	sessions, err := u.Sessions()
	if len(sessions) == 0 || err != nil {
		t.Fatalf("u.Session() = %v, %v, not want empty, error", sessions, err)
	}

	session := sessions[0]
	if session.ID == "" || session.Device == "" || session.TTL == 0 {
		t.Fatalf("u.Session() = %v, %v, not want empty, error", session, err)
	}
}

// TestSessionsExpired returns an empty session array because
// the session is expired
func TestSessionsExpired(t *testing.T) {
	sid, _ := stringutil.Random()
	id := rand.Int63n(10000)
	u := User{
		ID: id,
		Devices: map[string]string{
			"auth:" + sid: "Mozilla/5.0 Gecko/20100101 Firefox/115.0",
		},
	}

	sessions, err := u.Sessions()
	if len(sessions) != 0 || err != nil {
		t.Fatalf("u.Session() = %v, %v, want empty, error", sessions, err)
	}
}

// TestSessionsEmpty returns an empty session
func TestSessionsEmpty(t *testing.T) {
	id := rand.Int63n(10000)
	u := User{
		ID:      id,
		Devices: map[string]string{},
	}

	sessions, err := u.Sessions()
	if len(sessions) != 0 || err != nil {
		t.Fatalf("u.Session() = %v, %v, want empty, error", sessions, err)
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
