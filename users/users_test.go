package users

import (
	"gifthub/conf"
	"gifthub/db"
	"gifthub/tests"
	"testing"

	"github.com/go-faker/faker/v4"
)

var user User = User{
	ID:  1,
	SID: "test",
}

var ra faker.RealAddress = faker.GetRealAddress()
var address Address = Address{
	Lastname:      faker.Name(),
	Firstname:     faker.Name(),
	Street:        ra.Address,
	City:          ra.City,
	Complementary: ra.Address,
	Zipcode:       ra.PostalCode,
	Phone:         faker.Phonenumber(),
}

func TestListReturnsUsersWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	users, err := List(ctx, 0)
	if err != nil || len(users) == 0 || users[0].ID == 0 {
		t.Fatalf("List(ctx, 0) = '%v', %v, want User, nil", users, err)
	}
}

func TestOtpCodeReturnsCodeWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	otp, err := Otp(ctx, faker.Email())
	if otp == "" || err != nil {
		t.Fatalf("OtpCode(ctx, faker.Email()) = '%s', %v, want string, nil", otp, err)
	}
}

func TestOtpCodeReturnsCodeWhenUsedTwice(t *testing.T) {
	ctx := tests.Context()
	Otp(ctx, faker.Email())
	otp, err := Otp(ctx, faker.Email())
	if otp == "" || err != nil {
		t.Fatalf("OtpCode(ctx, faker.Email()) = '%s', %v, want string, nil", otp, err)
	}
}

func TestOtpCodeReturnsErrorWhenEmailIsEmpty(t *testing.T) {
	ctx := tests.Context()
	otp, err := Otp(ctx, "")
	if otp != "" || err == nil || err.Error() != "input_email_invalid" {
		t.Fatalf("OtpCode(ctx, '') = '%s', %v, want '', 'input_email_invalid'", otp, err)
	}
}

func TestOtpCodeReturnsErrorWhenEmailIsInvalid(t *testing.T) {
	ctx := tests.Context()
	otp, err := Otp(ctx, "toto")
	if otp != "" || err == nil || err.Error() != "input_email_invalid" {
		t.Fatalf("OtpCode(ctx, 'toto') = '%s', '%v', want '', 'input_email_invalid'", otp, err)
	}
}

func TestDeleteReturnsNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	err := User{ID: 2}.Delete(ctx)
	if err != nil {
		t.Fatalf("User{ID: 2}.Delete(ctx) = %v, want nil", err)
	}
}

func TestDeleteReturnsErrorWhenUserDoesNotExist(t *testing.T) {
	ctx := tests.Context()
	err := User{}.Delete(ctx)
	if err != nil {
		t.Fatalf("User{}.Delete(ctx) = %v, want nil", err)
	}
}

func TestLoginReturnsSidWhenSuccess(t *testing.T) {
	ctx := tests.Context()

	otp := "hello-world"

	db.Redis.Set(ctx, "hellow@world.com:otp", otp, conf.SessionDuration)
	db.Redis.Set(ctx, "otp:glue", "hellow@world.com", conf.SessionDuration)

	sid, err := Login(ctx, otp, "glue", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid == "" || err != nil {
		t.Fatalf(`Login(ctx, "hello-world", "glue", 'Mozilla/5.0 Gecko/20100101 Firefox/115.0') = '%s', %v, want string, nil`, sid, err)
	}
}

func TestLoginReturnsErrorWhenGlueIsEmpty(t *testing.T) {
	ctx := tests.Context()

	otp := "hello-world"

	sid, err := Login(ctx, otp, "", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil || err.Error() != "error_http_unauthorized" {
		t.Fatalf(`Login(ctx, "hello-world", "", 'Mozilla/5.0 Gecko/20100101 Firefox/115.0') = '%s', %v, want string, nil`, sid, err)
	}
}

func TestLoginReturnsErrorWhenOtpDoesNotExistForGlue(t *testing.T) {
	ctx := tests.Context()

	otp := "hello-world"

	db.Redis.Set(ctx, "otp:glue", "hellow@world.com", conf.SessionDuration)

	sid, err := Login(ctx, otp, "glue", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil || err.Error() != "error_http_unauthorized" {
		t.Fatalf(`Login(ctx, "hello-world", "glue", 'Mozilla/5.0 Gecko/20100101 Firefox/115.0') = '%s', %v, want string, nil`, sid, err)
	}
}

func TestLoginReturnsErrorWhenOtpDoesNotMatch(t *testing.T) {
	ctx := tests.Context()

	otp := "hello-world"

	db.Redis.Set(ctx, "hellow@world.com:otp", otp, conf.SessionDuration)
	db.Redis.Set(ctx, "otp:glue", "hellow@world.com", conf.SessionDuration)

	sid, err := Login(ctx, "hello", "glue", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil || err.Error() != "error_otp_mismatch" {
		t.Fatalf(`Login(ctx, "hello", "glue", 'Mozilla/5.0 Gecko/20100101 Firefox/115.0') = '%s', %v, want string, nil`, sid, err)
	}
}

func TestLoginReturnsErrorWhenOtpIsBlocked(t *testing.T) {
	ctx := tests.Context()

	otp := "hello-world"

	db.Redis.Set(ctx, "hellow@world.com:otp", otp, conf.SessionDuration)
	db.Redis.Set(ctx, "hellow@world.com:otp:attempts", 2, conf.SessionDuration)
	db.Redis.Set(ctx, "otp:glue", "hellow@world.com", conf.SessionDuration)

	sid, err := Login(ctx, "hello", "glue", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil || err.Error() != "error_otp_locked" {
		t.Fatalf(`Login(ctx, "hello", "glue", 'Mozilla/5.0 Gecko/20100101 Firefox/115.0') = '%s', %v, want string, nil`, sid, err)
	}
}

func TestLoginReturnsErrorWhenEmailOtpIsNotFound(t *testing.T) {
	ctx := tests.Context()

	otp := "hello-world"

	db.Redis.Set(ctx, "hellow@world.com:otp", otp, conf.SessionDuration)

	sid, err := Login(ctx, otp, "glue", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil || err.Error() != "error_http_unauthorized" {
		t.Fatalf(`Login(ctx, "hello-world", "glue", 'Mozilla/5.0 Gecko/20100101 Firefox/115.0') = '%s', %v, want string, nil`, sid, err)
	}
}

func TestLoginReturnsErrorWhenDeviceIsMissing(t *testing.T) {
	ctx := tests.Context()
	sid, err := Login(ctx, "otp", "glue", "")
	if sid != "" || err == nil || err.Error() != "error_login_device" {
		t.Fatalf(`Login(ctx, "otp", "glue", "") = '%s', %v, want '', 'error_login_device'`, sid, err)
	}
}

func TestLoginReturnsErrorWhenOtpIsMissing(t *testing.T) {
	ctx := tests.Context()
	glue := "glue"
	sid, err := Login(ctx, "", glue, "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil || err.Error() != "input_otp_required" {
		t.Fatalf(`Login(ctx, "", "glue", "Mozilla/5.0 Gecko/20100101 Firefox/115.0") = '%s', %v, want '', 'input_otp_required'`, sid, err)
	}
}

func TestLoginReturnsErrorWhenOtpDoesNotExist(t *testing.T) {
	ctx := tests.Context()
	glue := "glue"
	sid, err := Login(ctx, "titi", glue, "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil {
		t.Fatalf("Login(ctx, 'titi') = '%s', %v, want '', 'input_otp_required'", sid, err)
	}
}

func TestLogoutRetunsEmptyWhenSuccess(t *testing.T) {
	ctx := tests.Context()

	db.Redis.Set(ctx, "auth:"+"will-logout", "1", conf.SessionDuration)

	err := Logout(ctx, "will-logout")
	if err != nil {
		t.Fatalf("Logout(ctx, u.SID) = %v, want nil", err)
	}
}

func TestLogoutRetunsErrorWhenSidIsMissing(t *testing.T) {
	ctx := tests.Context()
	err := Logout(ctx, "")
	if err == nil || err.Error() != "error_http_unauthorized" {
		t.Fatalf("Logout(ctx, '') = %v, want 'error_user_logout'", err)
	}
}

func TestLogoutRetunsErrorWhenSessionIsExpired(t *testing.T) {
	ctx := tests.Context()
	err := Logout(ctx, "expired")
	if err == nil || err.Error() != "error_http_unauthorized" {
		t.Fatalf(`Logout(ctx, "expired") = %v, want 'error_user_logout'`, err)
	}
}

func TestLogoutRetunsErrorWhenSessionIsNotFound(t *testing.T) {
	ctx := tests.Context()
	err := Logout(ctx, "122")
	if err == nil || err.Error() != "error_http_unauthorized" {
		t.Fatalf("Logout(ctx, 124, '122') = %v, want nil", err)
	}
}

func TestSessionsReturnsSessionsWhenUserHasSession(t *testing.T) {
	ctx := tests.Context()
	sessions, err := user.Sessions(ctx)
	if len(sessions) == 0 || err != nil {
		t.Fatalf("user.Session(ctx) = %v, %v, want []Session, nil", sessions, err)
	}

	session := sessions[0]
	if session.ID == "" || session.Device == "" || session.TTL == 0 {
		t.Fatalf("sessions[0] = %v, want Session", session)
	}
}

func TestSessionsReturnsEmptySliceWhenUserDoesNotHaveSession(t *testing.T) {
	ctx := tests.Context()
	sessions, err := User{ID: 2}.Sessions(ctx)
	if len(sessions) != 0 || err != nil {
		t.Fatalf("User{ID: 2}.Session(ctx) = %v, %v, want []Session, nil", sessions, err)
	}
}

func TestSaveAddressReturnNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	err := user.SaveAddress(ctx, address)
	if err != nil {
		t.Fatalf("user.SaveAddress(ctx, address) = %v, want nil", err)
	}
}

func TestSaveAddressReturnNilWhenNoComplementary(t *testing.T) {
	a := address
	a.Complementary = ""

	ctx := tests.Context()
	err := user.SaveAddress(ctx, a)
	if err != nil {
		t.Fatalf("user.SaveAddress(ctx, a) = %v, want nil", err)
	}
}

func TestSaveAddressReturnErrorWhenUidIsEmpty(t *testing.T) {
	ctx := tests.Context()
	err := User{ID: 0}.SaveAddress(ctx, address)
	if err == nil || err.Error() != "error_http_general" {
		t.Fatalf("User{ID: 0}.SaveAddress(ctx, a) = %v, want 'error_http_general'", err)
	}
}

func TestSaveAddressReturnErrorWhenFirstnameIsEmpty(t *testing.T) {
	a := address
	a.Firstname = ""

	ctx := tests.Context()
	err := user.SaveAddress(ctx, a)
	if err == nil || err.Error() != "input_firstname_required" {
		t.Fatalf("user.SaveAddress(ctx, a) = %v, want 'input_firstname_required'", err)
	}
}

func TestSaveAddressReturnErrorWhenLastnameIsEmpty(t *testing.T) {
	a := address
	a.Lastname = ""

	ctx := tests.Context()
	err := user.SaveAddress(ctx, a)
	if err == nil || err.Error() != "input_lastname_required" {
		t.Fatalf("user.SaveAddress(ctx, a) = %v, want 'input_lastname_required'", err)
	}
}

func TestSaveAddressReturnErrorWhenStreeIsEmpty(t *testing.T) {
	a := address
	a.Street = ""

	ctx := tests.Context()
	err := user.SaveAddress(ctx, a)
	if err == nil || err.Error() != "input_street_required" {
		t.Fatalf("user.SaveAddress(ctx, a) = %v, want 'input_street_required'", err)
	}
}

func TestSaveAddressReturnErrorWhenCityIsEmpty(t *testing.T) {
	a := address
	a.City = ""

	ctx := tests.Context()
	err := user.SaveAddress(ctx, a)
	if err == nil || err.Error() != "input_city_required" {
		t.Fatalf("user.SaveAddress(ctx, a) = %v, want 'input_city_required'", err)
	}
}

func TestSaveAddressReturnErrorWhenZipcodeIsEmpty(t *testing.T) {
	a := address
	a.Zipcode = ""

	ctx := tests.Context()
	err := user.SaveAddress(ctx, a)
	if err == nil || err.Error() != "input_zipcode_required" {
		t.Fatalf("user.SaveAddress(ctx, a) = %v, want 'input_zipcode_required'", err)
	}
}

func TestSaveAddressReturnErrorWhenPhoneIsEmpty(t *testing.T) {
	a := address
	a.Phone = ""

	ctx := tests.Context()
	err := user.SaveAddress(ctx, a)
	if err == nil || err.Error() != "input_phone_required" {
		t.Fatalf("user.SaveAddress(ctx, a) = %v, want 'input_phone_required'", err)
	}
}

func TestGetReturnsUserWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	user, err := Get(ctx, 1)
	if err != nil || user.ID == 0 {
		t.Fatalf("users.Get(ctx, 1) = %v, %v, want User, nil", user, err)
	}
}

func TestGetReturnsErrorWhenIdIsEmpty(t *testing.T) {
	ctx := tests.Context()
	user, err := Get(ctx, 0)
	if err == nil || err.Error() != "error_session_notfound" || user.ID != 0 {
		t.Fatalf("users.Get(ctx, 0) = %v, %v, wan t User{}, 'error_session_notfound'", user, err)
	}
}

func TestGetReturnsErrorWhenUserDoesNotExist(t *testing.T) {
	ctx := tests.Context()
	user, err := Get(ctx, 123)
	if err == nil || err.Error() != "error_session_notfound" || user.ID != 0 {
		t.Fatalf("users.Get(ctx, 0) = %v, %v, want User{}, 'error_session_notfound'", user, err)
	}
}

func TestIsAdminReturnsFalseWhenUserIsNotAdmin(t *testing.T) {
	ctx := tests.Context()

	if IsAdmin(ctx, user.Email) {
		t.Fatalf("user.IsAdmin(ctx)= true, want false")
	}
}

func TestIsAdminReturnsTrueWhenUserIsAdmin(t *testing.T) {
	ctx := tests.Context()
	if !IsAdmin(ctx, "hello@world.com") {
		t.Fatalf("user.IsAdmin(ctx)= false, want true")
	}
}
