package users

import (
	"artisons/conf"
	"artisons/tests"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
)

var user User = User{
	ID:    tests.UserID1,
	SID:   tests.UserSID,
	Email: tests.UserEmail,
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

func TestOtpCodeReturnsCodeWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	err := Otp(ctx, faker.Email())
	if err != nil {
		t.Fatalf("OtpCode(ctx, faker.Email()) = %v, want nil", err)
	}
}

func TestOtpCodeReturnsCodeWhenUsedTwice(t *testing.T) {
	ctx := tests.Context()
	Otp(ctx, faker.Email())
	err := Otp(ctx, faker.Email())
	if err != nil {
		t.Fatalf("OtpCode(ctx, faker.Email()) = %v, want nil", err)
	}
}

func TestOtpCodeReturnsErrorWhenEmailIsEmpty(t *testing.T) {
	ctx := tests.Context()
	err := Otp(ctx, "")
	if err == nil || err.Error() != "input:email" {
		t.Fatalf("OtpCode(ctx, '') = %v, input:email", err)
	}
}

func TestOtpCodeReturnsErrorWhenEmailIsInvalid(t *testing.T) {
	ctx := tests.Context()
	err := Otp(ctx, tests.DoesNotExist)
	if err == nil || err.Error() != "input:email" {
		t.Fatalf("OtpCode(ctx, tests.DoesNotExist) = %v, want input:email", err)
	}
}

func TestDeleteReturnsNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	err := User{ID: tests.UserToDeleteID}.Delete(ctx)
	if err != nil {
		t.Fatalf("User{ID: test.UserToDeleteID}.Delete(ctx) = %v, want nil", err)
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

	sid, uid, err := Login(ctx, tests.AdminEmail, tests.Otp, tests.UA)
	if sid == "" || uid == 0 || err != nil {
		t.Fatalf(`Login(ctx, tests.AdminEmail,  tests.Otp, tests.UA) = '%s',%d,  %v, want string, int, nil`, sid, uid, err)
	}
}

func TestLoginReturnsErrorWhenEmailIsEmpty(t *testing.T) {
	ctx := tests.Context()

	sid, uid, err := Login(ctx, "", tests.Otp, tests.UA)
	if sid != "" || uid != 0 || err == nil || err.Error() != "input:email" {
		t.Fatalf(`Login(ctx, "", tests.Otp, tests.UA) = '%s', %d, %v, want '', 0, nil`, sid, uid, err)
	}
}

func TestLoginReturnsErrorWhenOtpDoesNotMatch(t *testing.T) {
	ctx := tests.Context()

	sid, uid, err := Login(ctx, tests.OtpNotMatching, "123457", tests.UA)
	if sid != "" || uid != 0 || err == nil || err.Error() != "the OTP does not match" {
		t.Fatalf(`Login(ctx, tests.OtpNotMatching, "123457", tests.UA) = '%s',%d, %v, want '', 0, nil`, sid, uid, err)
	}
}

func TestLoginReturnsErrorWhenOtpIsBlocked(t *testing.T) {
	ctx := tests.Context()

	sid, uid, err := Login(ctx, tests.AdminBlockedEmail, "1234567", tests.UA)
	if sid != "" || uid != 0 || err == nil || err.Error() != "you reached the max tentatives" {
		t.Fatalf(`Login(ctx, tests.AdminBlockedEmail, "1234567", tests.UA) = '%s', %d, %v, want '', 0, nil`, sid, uid, err)
	}
}

func TestLoginReturnsErrorWhenEmailOtpIsNotFound(t *testing.T) {
	ctx := tests.Context()

	sid, uid, err := Login(ctx, tests.EmailDoesNotExist, tests.Otp, tests.UA)
	if sid != "" || uid != 0 || err == nil || err.Error() != "you are not authorized to process this request" {
		t.Fatalf(`Login(ctx, tests.EmailDoesNotExist, tests.Otp,  tests.UA) = '%s',%d,  %v, want '', 0, nil`, sid, uid, err)
	}
}

func TestLoginReturnsErrorWhenDeviceIsMissing(t *testing.T) {
	ctx := tests.Context()
	sid, uid, err := Login(ctx, tests.AdminEmail, tests.Otp, "")
	if sid != "" || uid != 0 || err == nil || err.Error() != "your are not authorized to access to this page" {
		t.Fatalf(`Login(ctx, tests.AdminEmail, tests.Otp, "") = '%s', %d, %v, want '', 0, 'your are not authorized to access to this page'`, sid, uid, err)
	}
}

func TestLoginReturnsErrorWhenOtpIsMissing(t *testing.T) {
	ctx := tests.Context()
	sid, uid, err := Login(ctx, tests.AdminEmail, "", tests.UA)
	if sid != "" || uid != 0 || err == nil || err.Error() != "input:otp" {
		t.Fatalf(`Login(ctx, tests.AdminEmail, "", tests.UA) = '%s', %d,%v, want '', 0, 'input:otp'`, sid, uid, err)
	}
}

func TestLoginReturnsErrorWhenOtpDoesNotExist(t *testing.T) {
	ctx := tests.Context()
	sid, uid, err := Login(ctx, tests.AdminEmail, tests.DoesNotExist, tests.UA)
	if sid != "" || uid != 0 || err == nil {
		t.Fatalf(`Login(ctx, tests.AdminEmail, tests.DoesNotExist) = '%s', %d, %v, want '', 0, 'input:otp'`, sid, uid, err)
	}
}

func TestLogoutRetunsEmptyWhenSuccess(t *testing.T) {
	ctx := tests.Context()

	err := User{SID: tests.UserSIDSignedIn}.Logout(ctx)
	if err != nil {
		t.Fatalf("Logout(ctx, u.SID) = %v, want nil", err)
	}
}

func TestLogoutRetunsErrorWhenSidIsMissing(t *testing.T) {
	ctx := tests.Context()
	err := User{}.Logout(ctx)
	if err == nil || err.Error() != "you are not authorized to process this request" {
		t.Fatalf("Logout(ctx, '') = %v, want 'you are not authorized to process this request'", err)
	}
}

func TestLogoutRetunsErrorWhenSessionDoesNotExist(t *testing.T) {
	ctx := tests.Context()
	err := User{SID: tests.DoesNotExist}.Logout(ctx)
	if err == nil || err.Error() != "you are not authorized to process this request" {
		t.Fatalf(`Logout(ctx, tests.DoesNotExist) = %v, want 'you are not authorized to process this request'`, err)
	}
}

func TestSaveAddressReturnNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	err := address.Save(ctx, user.ID)
	if err != nil {
		t.Fatalf("user.SaveAddress(ctx, address) = %v, want nil", err)
	}
}

func TestSaveAddressReturnNilWhenNoComplementary(t *testing.T) {
	a := address
	a.Complementary = ""

	ctx := tests.Context()
	err := a.Save(ctx, user.ID)
	if err != nil {
		t.Fatalf("user.SaveAddress(ctx, a) = %v, want nil", err)
	}
}

func TestSaveAddressReturnErrorWhenUidIsEmpty(t *testing.T) {
	ctx := tests.Context()
	err := address.Save(ctx, 0)
	if err == nil || err.Error() != "something went wrong" {
		t.Fatalf("User{ID: 0}.SaveAddress(ctx, a) = %v, want 'something went wrong'", err)
	}
}

func TestSaveAddressReturnErrorWhenFirstnameIsEmpty(t *testing.T) {
	a := address
	a.Firstname = ""

	ctx := tests.Context()
	err := a.Validate(ctx)
	if err == nil || err.Error() != "input:firstname" {
		t.Fatalf("user.SaveAddress(ctx, a) = %v, want 'input:firstname'", err)
	}
}

func TestSaveAddressReturnErrorWhenLastnameIsEmpty(t *testing.T) {
	a := address
	a.Lastname = ""

	ctx := tests.Context()
	err := a.Validate(ctx)
	if err == nil || err.Error() != "input:lastname" {
		t.Fatalf("user.SaveAddress(ctx, a) = %v, want 'input:lastname'", err)
	}
}

func TestSaveAddressReturnErrorWhenStreeIsEmpty(t *testing.T) {
	a := address
	a.Street = ""

	ctx := tests.Context()
	err := a.Validate(ctx)
	if err == nil || err.Error() != "input:street" {
		t.Fatalf("user.SaveAddress(ctx, a) = %v, want 'input:street'", err)
	}
}

func TestSaveAddressReturnErrorWhenCityIsEmpty(t *testing.T) {
	a := address
	a.City = ""

	ctx := tests.Context()
	err := a.Validate(ctx)
	if err == nil || err.Error() != "input:city" {
		t.Fatalf("user.SaveAddress(ctx, a) = %v, want 'input:city'", err)
	}
}

func TestSaveAddressReturnErrorWhenZipcodeIsEmpty(t *testing.T) {
	a := address
	a.Zipcode = ""

	ctx := tests.Context()
	err := a.Validate(ctx)
	if err == nil || err.Error() != "input:zipcode" {
		t.Fatalf("user.SaveAddress(ctx, a) = %v, want 'input:zipcode'", err)
	}
}

func TestSaveAddressReturnErrorWhenPhoneIsEmpty(t *testing.T) {
	a := address
	a.Phone = ""

	ctx := tests.Context()
	err := a.Validate(ctx)
	if err == nil || err.Error() != "input:phone" {
		t.Fatalf("user.SaveAddress(ctx, a) = %v, want 'input:phone'", err)
	}
}

func TestGetReturnsUserWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	user, err := FindByUID(ctx, 1)
	if err != nil || user.ID == 0 {
		t.Fatalf("users.Get(ctx, 1) = %v, %v, want User, nil", user, err)
	}
}

func TestGetReturnsErrorWhenIdIsEmpty(t *testing.T) {
	ctx := tests.Context()
	user, err := FindByUID(ctx, 0)
	if err == nil || err.Error() != "the user is not found" || user.ID != 0 {
		t.Fatalf("users.Get(ctx, 0) = %v, %v, wan t User{}, 'the user is not found'", user, err)
	}
}

func TestGetReturnsErrorWhenUserDoesNotExist(t *testing.T) {
	ctx := tests.Context()
	user, err := FindByUID(ctx, tests.IDDoesNotExist)
	if err == nil || err.Error() != "the user is not found" || user.ID != 0 {
		t.Fatalf("users.Get(ctx, tests.IDDoesNotExist) = %v, %v, want User{}, 'the user is not found'", user, err)
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
	if !IsAdmin(ctx, tests.AdminEmail) {
		t.Fatalf("user.IsAdmin(ctx, tests.AdminEmail) = false, want true")
	}
}

func TestSearchReturnsUserWhenEmailMatching(t *testing.T) {
	c := tests.Context()
	u, err := Search(c, Query{Email: user.Email}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Email: user.Email}) = %v, want nil`, err.Error())
	}

	if u.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, u.Total)
	}

	if len(u.Users) == 0 {
		t.Fatalf(`len(p.Articles) = %d, want > 0`, len(u.Users))
	}

	if u.Users[0].ID != user.ID {
		t.Fatalf(`%d != %d`, u.Users[0].ID, user.ID)
	}
}

func TestSearchReturnsUserWhenEmailAndRoleching(t *testing.T) {
	c := tests.Context()
	u, err := Search(c, Query{Email: user.Email, Role: "user"}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Email: user.Email, Role: "user"}}) = %v, want nil`, err.Error())
	}

	if u.Total == 0 {
		t.Fatalf(`p.Total = %d, want > 0`, u.Total)
	}

	if len(u.Users) == 0 {
		t.Fatalf(`len(p.Articles) = %d, want > 0`, len(u.Users))
	}

	if u.Users[0].ID != user.ID {
		t.Fatalf(`%d != %d`, u.Users[0].ID, user.ID)
	}
}

func TestSearchReturnsNoUserWhenNoMatching(t *testing.T) {
	c := tests.Context()
	a, err := Search(c, Query{Email: tests.EmailDoesNotExist}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: tests.EmailDoesNotExist}) = %v, want nil`, err.Error())
	}

	if a.Total > 0 {
		t.Fatalf(`p.Total = %d, want == 0`, a.Total)
	}

	if len(a.Users) > 0 {
		t.Fatalf(`len(p.Articles) = %d, want == 0`, len(a.Users))
	}
}

func TestRefreshSessionReturnsNilWhenOK(t *testing.T) {
	ctx := tests.Context()
	u := User{SID: tests.UserRefreshSID}

	if err := u.RefreshSession(ctx); err != nil {
		t.Fatalf("u.RefreshSession(ctx) = %s, want nil", err.Error())
	}

	ttl := tests.TTL(ctx, "session:"+tests.UserRefreshSID)
	if ttl <= time.Minute {
		t.Fatalf(`ttl = %d want %d`, ttl, conf.SessionDuration)
	}
}
