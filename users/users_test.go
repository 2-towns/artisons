package users

import (
	"context"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/string/stringutil"
	"gifthub/tests"
	"math/rand"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
)

func createUser(alive bool) User {
	email := faker.Email()
	sid, _ := stringutil.Random()
	id := rand.Int63n(10000)
	ctx := context.Background()
	now := time.Now()

	db.Redis.HSet(ctx, fmt.Sprintf("user:%d", id), "email", email)
	db.Redis.HSet(ctx, fmt.Sprintf("user:%d", id), "id", id)
	db.Redis.HSet(ctx, fmt.Sprintf("user:%d", id), "created_at", now.Format(time.RFC3339))
	db.Redis.HSet(ctx, fmt.Sprintf("user:%d", id), "updated_at", now.Format(time.RFC3339))

	if alive {
		db.Redis.Set(ctx, "auth:"+sid, id, conf.SessionDuration)
		db.Redis.HSet(ctx, fmt.Sprintf("user:%d", id), "auth:"+sid, email)
	}

	return User{
		SID: sid,
		Devices: map[string]string{
			"auth:" + sid: "Mozilla/5.0 Gecko/20100101 Firefox/115.0",
		},
		ID: id,
	}
}

func createLinkedMagic() string {
	magic, _ := stringutil.Random()
	id := rand.Intn(10000)
	ctx := context.Background()
	db.Redis.Set(ctx, "magic:"+magic, id, 0)

	return magic
}

func createAddress() Address {
	ra := faker.GetRealAddress()
	return Address{
		Lastname:      faker.Name(),
		Firstname:     faker.Name(),
		Address:       ra.Address,
		City:          ra.City,
		Complementary: ra.Address,
		Zipcode:       ra.PostalCode,
		Phone:         faker.Phonenumber(),
	}

}

// TestUserList expects to succeed
func TestUserList(t *testing.T) {
	ctx := tests.Context()
	users, err := List(ctx, 0)
	if err != nil || len(users) == 0 || users[0].ID == 0 {
		t.Fatalf("List(ctx, 0) = '%v', %v, want User, nil", users, err)
	}
}

// TestMagicCode expects to succeed
func TestMagicCode(t *testing.T) {
	ctx := tests.Context()
	magic, err := MagicCode(ctx, faker.Email())
	if magic == "" || err != nil {
		t.Fatalf("TestMagicCode(ctx, faker.Email()) = '%s', %v, want string, nil", magic, err)
	}
}

// TestMagicCodeTwice expects to succeed even if it's used more than one time
func TestMagicCodeTwice(t *testing.T) {
	ctx := tests.Context()
	MagicCode(ctx, faker.Email())
	magic, err := MagicCode(ctx, faker.Email())
	if magic == "" || err != nil {
		t.Fatalf("TestMagicCode(ctx, faker.Email()) = '%s', %v, want string, nil", magic, err)
	}
}

// TestMagicCodeWithoutEmail expects to fail because of email emptyness
func TestMagicCodeWithoutEmail(t *testing.T) {
	ctx := tests.Context()
	magic, err := MagicCode(ctx, "")
	if magic != "" || err == nil || err.Error() != "user_email_invalid" {
		t.Fatalf("TestMagicCode(ctx, '') = '%s', %v, want '', 'user_email_invalid'", magic, err)
	}
}

// TestMagicCodeFailedWithBadEmail expects to fail because of email misvalue
func TestMagicCodeFailedWithBadEmail(t *testing.T) {
	ctx := tests.Context()
	magic, err := MagicCode(ctx, "toto")
	if magic != "" || err == nil || err.Error() != "user_email_invalid" {
		t.Fatalf("TestMagicCode(ctx, 'toto') = '%s', %v, want '', 'user_email_invalid'", magic, err)
	}
}

// TestDeleteUser expects to succeed
func TestDeleteUser(t *testing.T) {
	alive := true
	u := createUser(alive)
	ctx := tests.Context()
	err := u.Delete(ctx)
	if err != nil {
		t.Fatalf("Delete(ctx) = %v, want nil", err)
	}
}

// TestDeleteUserNotExisting expects to succeed
func TestDeleteUserNotExisting(t *testing.T) {
	ctx := tests.Context()
	err := User{}.Delete(ctx)
	if err != nil {
		t.Fatalf("Delete(ctx) = %v, want nil", err)
	}
}

// Testlogin expects to succeed
func TestLogin(t *testing.T) {
	magic := createLinkedMagic()
	ctx := tests.Context()
	sid, err := Login(ctx, magic, "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid == "" || err != nil {
		t.Fatalf("Login(ctx, magic, 'Mozilla/5.0 Gecko/20100101 Firefox/115.0') = '%s', %v, want string, nil", sid, err)
	}
}

// TestLoginWithoutDevice expects to fail because of device emptyness
func TestLoginWithoutDevice(t *testing.T) {
	magic := createLinkedMagic()
	ctx := tests.Context()
	sid, err := Login(ctx, magic, "")
	if sid != "" || err == nil || err.Error() != "user_device_required" {
		t.Fatalf("Login(ctx, magic,'') = '%s', %v, want '', 'user_device_required'", sid, err)
	}
}

// TestLoginWithoutMagic expects to fail because of magic emptyness
func TestLoginWithoutMagic(t *testing.T) {
	ctx := tests.Context()
	sid, err := Login(ctx, "", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil || err.Error() != "user_magic_code_required" {
		t.Fatalf("Login(ctx, '','Mozilla/5.0 Gecko/20100101 Firefox/115.0') = '%s', %v, want '', 'user_magic_code_required'", sid, err)
	}
}

// TestLoginWithNotExistingMagic expects to fail because of magic non existence
func TestLoginWithNotExistingMagic(t *testing.T) {
	ctx := tests.Context()
	sid, err := Login(ctx, "titi", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil {
		t.Fatalf("Login(ctx, 'titi') = '%s', %v, want '', 'user_magic_code_required'", sid, err)
	}
}

// TestLogout expects to succeed
func TestLogout(t *testing.T) {
	alive := true
	u := createUser(alive)
	ctx := tests.Context()
	err := Logout(ctx, u.SID)
	if err != nil {
		t.Fatalf("Logout(ctx, u.SID) = %v, want nil", err)
	}
}

// TestLogoutWithoutSID expects to fail because of id emptyness
func TestLogoutWithoutSID(t *testing.T) {
	ctx := tests.Context()
	err := Logout(ctx, "")
	if err == nil || err.Error() != "unauthorized" {
		t.Fatalf("Logout(ctx, '') = %v, want 'user_logout_invalid'", err)
	}
}

// TestLogoutWithExpiredSession expects to fail because of session expiration
func TestLogoutWithExpiredSession(t *testing.T) {
	alive := false
	u := createUser(alive)
	ctx := tests.Context()
	err := Logout(ctx, u.SID)
	if err == nil || err.Error() != "unauthorized" {
		t.Fatalf("Logout(ctx, u.SID) = %v, want 'user_logout_invalid'", err)
	}
}

// TestLogoutWithNotExistingData expects to fails because of session misvalue
func TestLogoutWithNotExistingData(t *testing.T) {
	ctx := tests.Context()
	err := Logout(ctx, "122")
	if err == nil || err.Error() != "unauthorized" {
		t.Fatalf("Logout(ctx, 124, '122') = %v, want nil", err)
	}
}

// TestSessions expects to succeed
func TestSessions(t *testing.T) {
	alive := true
	u := createUser(alive)

	ctx := tests.Context()
	sessions, err := u.Sessions(ctx)
	if len(sessions) == 0 || err != nil {
		t.Fatalf("u.Session(ctx) = %v, %v, want []Session, nil", sessions, err)
	}

	session := sessions[0]
	if session.ID == "" || session.Device == "" || session.TTL == 0 {
		t.Fatalf("sessions[0] = %v, want Session", session)
	}
}

// TestSessionsExpired expects to succeed with empty array when sessions are expired
func TestSessionsExpired(t *testing.T) {
	alive := false
	u := createUser(alive)

	ctx := tests.Context()
	sessions, err := u.Sessions(ctx)
	if len(sessions) != 0 || err != nil {
		t.Fatalf("u.Session(ctx) = %v, %v, want []Session, nil", sessions, err)
	}
}

// TestSessionsEmpty expects to succeed with empty array when sessions are empty
func TestSessionsEmpty(t *testing.T) {
	alive := true
	u := createUser(alive)
	u.Devices = map[string]string{}

	ctx := tests.Context()
	sessions, err := u.Sessions(ctx)
	if len(sessions) != 0 || err != nil {
		t.Fatalf("u.Session(ctx) = %v, %v, want [], nil", sessions, err)
	}
}

// TestSaveAddress expects to succeed
func TestSaveAddress(t *testing.T) {
	a := createAddress()
	alive := true
	u := createUser(alive)

	ctx := tests.Context()
	err := u.SaveAddress(ctx, a)
	if err != nil {
		t.Fatalf("SaveAddress(ctx, a) = %v, want nil", err)
	}
}

// TestSaveAddressWithoutComplementary expects to succeed with empty complementary
func TestSaveAddressWithoutComplementary(t *testing.T) {
	a := createAddress()
	a.Complementary = ""
	alive := true
	u := createUser(alive)

	ctx := tests.Context()
	err := u.SaveAddress(ctx, a)
	if err != nil {
		t.Fatalf("SaveAddress(ctx, a) = %v, want nil", err)
	}
}

// TestSaveAddressUIDEmpty expects to fail because of id emptyness
func TestSaveAddressUIDEmpty(t *testing.T) {
	a := createAddress()
	alive := true
	u := createUser(alive)
	u.ID = 0

	ctx := tests.Context()
	err := u.SaveAddress(ctx, a)
	if err == nil || err.Error() != "something_went_wrong" {
		t.Fatalf("SaveAddress(ctx, a) = %v, want 'something_went_wrong'", err)
	}
}

// TestSaveAddressWithoutFirstname expects to fail because of firstname emptyness
func TestSaveAddressWithoutFirstname(t *testing.T) {
	a := createAddress()
	a.Firstname = ""
	alive := true
	u := createUser(alive)

	ctx := tests.Context()
	err := u.SaveAddress(ctx, a)
	if err == nil || err.Error() != "user_firstname_required" {
		t.Fatalf("SaveAddress(ctx, a) = %v, want 'user_firstname_required'", err)
	}
}

// TestSaveAddressWithoutLastname expects to fail because of lastname emptyness
func TestSaveAddressWithoutLastname(t *testing.T) {
	a := createAddress()
	a.Lastname = ""
	alive := true
	u := createUser(alive)

	ctx := tests.Context()
	err := u.SaveAddress(ctx, a)
	if err == nil || err.Error() != "user_lastname_required" {
		t.Fatalf("SaveAddress(ctx, a) = %v, want 'user_lastname_required'", err)
	}
}

// TestSaveAddressWithoutAddress expects to fail because of address emptyness
func TestSaveAddressWithoutAddress(t *testing.T) {
	a := createAddress()
	a.Address = ""
	alive := true
	u := createUser(alive)

	ctx := tests.Context()
	err := u.SaveAddress(ctx, a)
	if err == nil || err.Error() != "user_address_required" {
		t.Fatalf("SaveAddress(ctx, a) = %v, want 'user_address_required'", err)
	}
}

// TestSaveAddressWithoutCity expects to fail because of city emptyness
func TestSaveAddressWithoutCity(t *testing.T) {
	a := createAddress()
	a.City = ""
	alive := true
	u := createUser(alive)

	ctx := tests.Context()
	err := u.SaveAddress(ctx, a)
	if err == nil || err.Error() != "user_city_required" {
		t.Fatalf("SaveAddress(ctx, a) = %v, want 'user_city_required'", err)
	}
}

// TestSaveAddressWithoutZipcode expects to fail because of zipcode emptyness
func TestSaveAddressWithoutZipcode(t *testing.T) {
	a := createAddress()
	a.Zipcode = ""
	alive := true
	u := createUser(alive)

	ctx := tests.Context()
	err := u.SaveAddress(ctx, a)
	if err == nil || err.Error() != "user_zipcode_required" {
		t.Fatalf("SaveAddress(ctx, a) = %v, want 'user_zipcode_required'", err)
	}
}

// TestSaveAddressWithoutPhone expects to fail because of phone emptyness
func TestSaveAddressWithoutPhone(t *testing.T) {
	a := createAddress()
	a.Phone = ""
	alive := true
	u := createUser(alive)

	ctx := tests.Context()
	err := u.SaveAddress(ctx, a)
	if err == nil || err.Error() != "user_phone_required" {
		t.Fatalf("SaveAddress(ctx, a) = %v, want 'user_phone_required'", err)
	}
}

// TestGetUser expects to succeed
func TestGetUser(t *testing.T) {
	alive := true
	u := createUser(alive)

	ctx := tests.Context()
	user, err := Get(ctx, u.ID)
	if err != nil || user.ID == 0 {
		t.Fatalf("users.Get(ctx, u.ID) = %v, %v, want User, nil", user, err)
	}
}

// TestGetUserEmpty expects to fail because of id emptyness
func TestGetUserEmpty(t *testing.T) {
	ctx := tests.Context()
	user, err := Get(ctx, 0)
	if err == nil || err.Error() != "user_not_found" || user.ID != 0 {
		t.Fatalf("users.Get(ctx, 0) = %v, %v, wan t User{}, 'user_not_found'", user, err)
	}
}

// TestGetUserNotExisting expects to fail because of user non existence
func TestGetUserNotExisting(t *testing.T) {
	ctx := tests.Context()
	user, err := Get(ctx, 123)
	if err == nil || err.Error() != "user_not_found" || user.ID != 0 {
		t.Fatalf("users.Get(ctx, 0) = %v, %v, want User{}, 'user_not_found'", user, err)
	}
}
