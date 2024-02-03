package users

import (
	"artisons/conf"
	"artisons/db"
	"artisons/tests"
	"fmt"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
)

var user User = User{
	ID:    99,
	SID:   "SES99",
	Email: faker.Email(),
}

func init() {
	ctx := tests.Context()

	now := time.Now()

	db.Redis.HSet(ctx, fmt.Sprintf("user:%d", user.ID),
		"id", user.ID,
		"email", user.Email,
		"created_at", now.Unix(),
		"updated_at", now.Unix(),
		"type", "user",
		"role", "user",
	)

	db.Redis.HSet(ctx, "user:98",
		"id", 98,
		"email", faker.Email(),
		"created_at", now.Unix(),
		"updated_at", now.Unix(),
		"type", "user",
		"role", "user",
	)

	db.Redis.HSet(ctx, "user:100",
		"id", 100,
		"email", "hello@world.com",
		"created_at", now.Unix(),
		"updated_at", now.Unix(),
		"type", "user",
		"role", "admin",
	)

	db.Redis.HSet(ctx, "session:"+user.SID, "id", user.SID, "uid", user.ID, "device", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0", "type", "session")
	db.Redis.Expire(ctx, "session:"+user.SID, conf.SessionDuration)
	db.Redis.HSet(ctx, "session:will-logout", "id", "will-logout", "uid", "1", "type", "session")
	db.Redis.Expire(ctx, "session:will-logout", conf.SessionDuration)
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
	err := Otp(ctx, "toto")
	if err == nil || err.Error() != "input:email" {
		t.Fatalf("OtpCode(ctx, 'toto') = %v, want input:email", err)
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

	db.Redis.HSet(ctx, "otp:hellow@world.com", "otp", otp, "attempts", 0)
	db.Redis.Expire(ctx, "otp:hellow@world.com", conf.SessionDuration)

	sid, err := Login(ctx, "hellow@world.com", otp, "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid == "" || err != nil {
		t.Fatalf(`Login(ctx, "hellow@world.com", "hello-world", 'Mozilla/5.0 Gecko/20100101 Firefox/115.0') = '%s', %v, want string, nil`, sid, err)
	}
}

func TestLoginReturnsErrorWhenEmailIsEmpty(t *testing.T) {
	ctx := tests.Context()

	otp := "hello-world"

	sid, err := Login(ctx, "", otp, "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil || err.Error() != "input:email" {
		t.Fatalf(`Login(ctx, "", otp, 'Mozilla/5.0 Gecko/20100101 Firefox/115.0') = '%s', %v, want string, nil`, sid, err)
	}
}

func TestLoginReturnsErrorWhenOtpDoesNotMatch(t *testing.T) {
	ctx := tests.Context()

	otp := "hello-world"

	db.Redis.HSet(ctx, "otp:hellow@world.com", "otp", otp, "attempts", 0)
	db.Redis.Expire(ctx, "otp:hellow@world.com", conf.SessionDuration)

	sid, err := Login(ctx, "hellow@world.com", "hello", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil || err.Error() != "the OTP does not match" {
		t.Fatalf(`Login(ctx, "hellow@world.com", "hello", 'Mozilla/5.0 Gecko/20100101 Firefox/115.0') = '%s', %v, want string, nil`, sid, err)
	}
}

func TestLoginReturnsErrorWhenOtpIsBlocked(t *testing.T) {
	ctx := tests.Context()

	otp := "hello-world"

	db.Redis.HSet(ctx, "otp:hellow@world.com", "otp", otp, "attempts", 2)
	db.Redis.Expire(ctx, "otp:hellow@world.com", conf.SessionDuration)

	sid, err := Login(ctx, "hellow@world.com", "hello", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil || err.Error() != "you reached the max tentatives" {
		t.Fatalf(`Login(ctx, "hellow@world.com", "hello", 'Mozilla/5.0 Gecko/20100101 Firefox/115.0') = '%s', %v, want string, nil`, sid, err)
	}
}

func TestLoginReturnsErrorWhenEmailOtpIsNotFound(t *testing.T) {
	ctx := tests.Context()

	otp := "hello-world"

	db.Redis.HSet(ctx, "otp:hellow@world.com", "otp", otp, "attempts", 0)
	db.Redis.Expire(ctx, "otp:hellow@world.com", conf.SessionDuration)

	sid, err := Login(ctx, "crazy@world.com", otp, "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil || err.Error() != "you are not authorized to process this request" {
		t.Fatalf(`Login(ctx, "crazy@world.com", "hello-world",  'Mozilla/5.0 Gecko/20100101 Firefox/115.0') = '%s', %v, want string, nil`, sid, err)
	}
}

func TestLoginReturnsErrorWhenDeviceIsMissing(t *testing.T) {
	ctx := tests.Context()
	sid, err := Login(ctx, "hello@world.com", "otp", "")
	if sid != "" || err == nil || err.Error() != "your are not authorized to access to this page" {
		t.Fatalf(`Login(ctx, "hello@world.com", "otp", "") = '%s', %v, want '', 'your are not authorized to access to this page'`, sid, err)
	}
}

func TestLoginReturnsErrorWhenOtpIsMissing(t *testing.T) {
	ctx := tests.Context()
	sid, err := Login(ctx, "hello@world.com", "", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil || err.Error() != "input:otp" {
		t.Fatalf(`Login(ctx, "hello@world.com", "", "Mozilla/5.0 Gecko/20100101 Firefox/115.0") = '%s', %v, want '', 'input:otp'`, sid, err)
	}
}

func TestLoginReturnsErrorWhenOtpDoesNotExist(t *testing.T) {
	ctx := tests.Context()
	sid, err := Login(ctx, "hello@world.com", "titi", "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if sid != "" || err == nil {
		t.Fatalf(`Login(ctx, "hello@world.com", 'titi') = '%s', %v, want '', 'input:otp'`, sid, err)
	}
}

func TestLogoutRetunsEmptyWhenSuccess(t *testing.T) {
	ctx := tests.Context()

	err := User{SID: "will-logout"}.Logout(ctx)
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
	err := User{SID: "iamnotexisting"}.Logout(ctx)
	if err == nil || err.Error() != "you are not authorized to process this request" {
		t.Fatalf(`Logout(ctx, "iamnotexisting") = %v, want 'you are not authorized to process this request'`, err)
	}
}

func TestLogoutRetunsErrorWhenSessionIsNotFound(t *testing.T) {
	ctx := tests.Context()
	err := User{SID: "122"}.Logout(ctx)
	if err == nil || err.Error() != "you are not authorized to process this request" {
		t.Fatalf("Logout(ctx, 124, '122') = %v, want nil", err)
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
	user, err := FindByUID(ctx, 123)
	if err == nil || err.Error() != "the user is not found" || user.ID != 0 {
		t.Fatalf("users.Get(ctx, 0) = %v, %v, want User{}, 'the user is not found'", user, err)
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
	a, err := Search(c, Query{Email: "crazy@world.com"}, 0, 10)
	if err != nil {
		t.Fatalf(`Search(c, Query{Keywords: ""crazy world"}) = %v, want nil`, err.Error())
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
	u := User{SID: "crazy"}

	db.Redis.Expire(ctx, "session:crazy", 10)

	if err := u.RefreshSession(ctx); err != nil {
		t.Fatalf("u.RefreshSession(ctx) = %s, want nil", err.Error())
	}

	ttl, err := db.Redis.TTL(ctx, "session:crazy").Result()
	if ttl >= conf.SessionDuration || err != nil {
		t.Fatalf(`db.Redis.TTL(ctx,"session:crazy") = %d, %v, want %d, nil`, ttl, err, conf.SessionDuration)
	}
}
