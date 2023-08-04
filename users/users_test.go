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
		t.Fatalf("List(0) = '%v', %v, want User, nil", users, err)
	}
}

// TestMagicCode generate a magic code for an email
func TestMagicCode(t *testing.T) {
	magic, err := MagicCode(faker.Email())
	if magic == "" || err != nil {
		t.Fatalf("TestMagicCode(faker.Email()) = '%s', %v, want string, nil", magic, err)
	}
}

// TestMagicCodeTwice generate a magic code for an email
func TestMagicCodeTwice(t *testing.T) {
	magic, _ := MagicCode(faker.Email())
	magic, err := MagicCode(faker.Email())
	if magic == "" || err != nil {
		t.Fatalf("TestMagicCode(faker.Email()) = '%s', %v, want string, nil", magic, err)
	}
}

// TestMagicCodeWithEmailEmpty fails when email is empty
func TestMagicCodeWithEmailEmpty(t *testing.T) {
	magic, err := MagicCode("")
	if magic != "" || err == nil || err.Error() != "user_email_invalid" {
		t.Fatalf("TestMagicCode('') = '%s', %v, want '', 'user_email_invalid'", magic, err)
	}
}

// TestMagicCodeFailedWithBadEmail fails when email is incorrect
func TestMagicCodeFailedWithBadEmail(t *testing.T) {
	magic, err := MagicCode("toto")
	if magic != "" || err == nil || err.Error() != "user_email_invalid" {
		t.Fatalf("TestMagicCode('toto') = '%s', %v, want '', 'user_email_invalid'", magic, err)
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
		t.Fatalf("Delete() = %v, want nil", err)
	}
}

// TestDeleteUserNotExisting does nothing
func TestDeleteUserNotExisting(t *testing.T) {
	err := User{}.Delete()
	if err != nil {
		t.Fatalf("Delete() = %v, want nil", err)
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
		t.Fatalf("TestLogin(magic) = '%s', %v, want string, nil", sid, err)
	}
}

// TestLoginWithoutDevice fails when device is empty
func TestLoginWithoutDevice(t *testing.T) {
	magic, _ := stringutil.Random()
	id := rand.Intn(10000)
	ctx := context.Background()
	db.Redis.Set(ctx, "magic:"+magic, id, 0)

	sid, err := Login(magic, "")
	if sid != "" || err == nil || err.Error() != "user_device_required" {
		t.Fatalf("TestLogin(magic) = '%s', %v, want '', 'user_device_required'", sid, err)
	}
}

// TestLoginWithEmptyMagic try to authenticate an user with empty magic
func TestLoginWithEmptyMagic(t *testing.T) {
	sid, err := Login("", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil || err.Error() != "user_magic_code_required" {
		t.Fatalf("TestLogin('') = '%s', %v, want '', 'user_magic_code_required'", sid, err)
	}
}

// TestLoginWithNotExistingMagic try to authenticate an user with not existing magic
func TestLoginWithNotExistingMagic(t *testing.T) {
	sid, err := Login("titi", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil {
		t.Fatalf("TestLogin('') = '%s', %v, want '', 'user_magic_code_required'", sid, err)
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
		t.Fatalf("Logout(id, sid) = %v, want nil", err)
	}
}

// TestLogoutWithZeroID returns an error
func TestLogoutWithZeroID(t *testing.T) {
	err := Logout(0, "123")
	if err == nil || err.Error() != "user_logout_invalid" {
		t.Fatalf("Logout(0, '123') = %v, want 'user_logout_invalid'", err)
	}
}

// TestLogoutWithEmptySID returns an error
func TestLogoutWithEmptySID(t *testing.T) {
	err := Logout(124, "")
	if err == nil || err.Error() != "user_logout_invalid" {
		t.Fatalf("Logout(124, '') = %v, want 'user_logout_invalid'", err)
	}
}

// TestLogoutWithNotExistingData does nothing
func TestLogoutWithNotExistingData(t *testing.T) {
	err := Logout(124, "122")
	if err != nil {
		t.Fatalf("Logout(124, '122') = %v, want nil", err)
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
		t.Fatalf("u.Session() = %v, %v, want []Session, nil", sessions, err)
	}

	session := sessions[0]
	if session.ID == "" || session.Device == "" || session.TTL == 0 {
		t.Fatalf("sessions[0] = %v, want Session", session)
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
		t.Fatalf("u.Session() = %v, %v, want []Session, nil", sessions, err)
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
		t.Fatalf("u.Session() = %v, %v, want [], nil", sessions, err)
	}
}

// TestSaveAddress saves a new address into Redis
func TestSaveAddress(t *testing.T) {
	ra := faker.GetRealAddress()
	a := Address{
		Lastname:      faker.Name(),
		Firstname:     faker.Name(),
		Address:       ra.Address,
		City:          ra.City,
		Complementary: "",
		Zipcode:       ra.PostalCode,
		Phone:         faker.Phonenumber(),
	}

	id := rand.Int63n(10000)
	u := User{
		ID: id,
	}
	err := u.SaveAddress(a)
	if err != nil {
		t.Fatalf("SaveAddress(a) = %v, want nil", err)
	}
}

// TestSaveAddressWithoutFirstname fails
func TestSaveAddressWithoutFirstname(t *testing.T) {
	ra := faker.GetRealAddress()
	a := Address{
		Lastname:      faker.Name(),
		Address:       ra.Address,
		City:          ra.City,
		Complementary: "",
		Zipcode:       ra.PostalCode,
		Phone:         faker.Phonenumber(),
	}

	id := rand.Int63n(10000)
	u := User{
		ID: id,
	}
	err := u.SaveAddress(a)
	if err == nil || err.Error() != "user_firstname_required" {
		t.Fatalf("SaveAddress(a) = %v, want 'user_firstname_required'", err)
	}
}

// TestSaveAddressWithoutLastname fails
func TestSaveAddressWithoutLastname(t *testing.T) {
	ra := faker.GetRealAddress()
	a := Address{
		Firstname:     faker.Name(),
		Address:       ra.Address,
		City:          ra.City,
		Complementary: "",
		Zipcode:       ra.PostalCode,
		Phone:         faker.Phonenumber(),
	}

	id := rand.Int63n(10000)
	u := User{
		ID: id,
	}
	err := u.SaveAddress(a)
	if err == nil || err.Error() != "user_lastname_required" {
		t.Fatalf("SaveAddress(a) = %v, want 'user_lastname_required'", err)
	}
}

// TestSaveAddressWithoutAddress fails
func TestSaveAddressWithoutAddress(t *testing.T) {
	ra := faker.GetRealAddress()
	a := Address{
		Firstname:     faker.Name(),
		Lastname:      faker.Name(),
		City:          ra.City,
		Complementary: "",
		Zipcode:       ra.PostalCode,
		Phone:         faker.Phonenumber(),
	}

	id := rand.Int63n(10000)
	u := User{
		ID: id,
	}
	err := u.SaveAddress(a)
	if err == nil || err.Error() != "user_address_required" {
		t.Fatalf("SaveAddress(a) = %v, want 'user_address_required'", err)
	}
}

// TestSaveAddressWithoutCity fails
func TestSaveAddressWithoutCity(t *testing.T) {
	ra := faker.GetRealAddress()
	a := Address{
		Lastname:      faker.Name(),
		Firstname:     faker.Name(),
		Address:       ra.Address,
		Complementary: "",
		Zipcode:       ra.PostalCode,
		Phone:         faker.Phonenumber(),
	}

	id := rand.Int63n(10000)
	u := User{
		ID: id,
	}
	err := u.SaveAddress(a)
	if err == nil || err.Error() != "user_city_required" {
		t.Fatalf("SaveAddress(a) = %v, want 'user_city_required'", err)
	}
}

// TestSaveAddressWithoutZipcode fails
func TestSaveAddressWithoutZipcode(t *testing.T) {
	ra := faker.GetRealAddress()
	a := Address{
		Lastname:      faker.Name(),
		Firstname:     faker.Name(),
		Address:       ra.Address,
		Complementary: "",
		City:          ra.City,
		Phone:         faker.Phonenumber(),
	}

	id := rand.Int63n(10000)
	u := User{
		ID: id,
	}
	err := u.SaveAddress(a)
	if err == nil || err.Error() != "user_zipcode_required" {
		t.Fatalf("SaveAddress(a) = %v, want 'user_zipcode_required'", err)
	}
}

// TestSaveAddressWithoutPhone fails
func TestSaveAddressWithoutPhone(t *testing.T) {
	ra := faker.GetRealAddress()
	a := Address{
		Lastname:      faker.Name(),
		Firstname:     faker.Name(),
		Address:       ra.Address,
		Complementary: "",
		City:          ra.City,
		Zipcode:       ra.PostalCode,
		Phone:         "",
	}

	id := rand.Int63n(10000)
	u := User{
		ID: id,
	}
	err := u.SaveAddress(a)
	if err == nil || err.Error() != "user_phone_required" {
		t.Fatalf("SaveAddress(a) = %v, want 'user_phone_required'", err)
	}
}
