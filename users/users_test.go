package users

import (
	"artisons/conf"
	"artisons/db"
	"artisons/tests"
	"errors"
	"fmt"
	"path"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
)

var cur string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	cur = path.Dir(filename) + "/"
}

const ua = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"

var user User = User{
	ID:    1,
	SID:   "123456789",
	Email: "arnaud@artisons.me",
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

func TestOtp(t *testing.T) {
	ctx := tests.Context()

	var tests = []struct {
		name  string
		email string
		err   error
	}{
		{"success", faker.Email(), nil},
		{"email=", "", errors.New("input:email")},
		{"email=idontexist", "", errors.New("input:email")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Otp(ctx, tt.email)
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.err) {
				t.Fatalf("err = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/users.redis")

	var tests = []struct {
		name string
		id   int
		err  error
	}{
		{"success", 2, nil},
		{"id=", 0, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := User{ID: tt.id}.Delete(ctx)
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.err) {
				t.Fatalf("err = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/users.redis")

	var cases = []struct {
		name  string
		email string
		otp   string
		ua    string
		err   error
	}{
		{"success", "hello@artisons.me", "123456", ua, nil},
		{"email=", "", "123456", ua, errors.New("input:email")},
		{"otp=111111", "otp@artisons.me", "111111", ua, errors.New("the OTP does not match")},
		{"otpattempts=3", "blocked@artisons.me", "111111", ua, errors.New("you reached the max tentatives")},
		{"email=idontexist@artisons.me", "idontexist@artisons.me", "111111", ua, errors.New("you are not authorized to process this request")},
		{"device=", "otp@artisons.me", "123456", "", errors.New("your are not authorized to access to this page")},
		{"otp=333333", "hello@artisons.me", "333333", ua, errors.New("you are not authorized to process this request")},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Login(ctx, tt.email, tt.otp, tt.ua)
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.err) {
				t.Fatalf(`err = %v, want %s`, err, tt.err)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/users.redis")

	var cases = []struct {
		name string
		sid  string
		err  error
	}{
		{"success", "123456789", nil},
		{"sid=", "", errors.New("you are not authorized to process this request")},
		{"sid=idontexist", "idontexist", errors.New("you are not authorized to process this request")},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := User{SID: tt.sid}.Logout(ctx)
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.err) {
				t.Fatalf(`err = %v, want %s`, err, tt.err)
			}
		})
	}
}

func TestSaveAddress(t *testing.T) {
	ctx := tests.Context()
	err := address.Save(ctx, user.ID)

	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}

	err = address.Save(ctx, 0)

	if err == nil {
		t.Fatal("err = nil, want something went wrong")
	}
}

func TestValidateAddress(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/users.redis")

	var cases = []struct {
		name  string
		uid   int
		field string
		value string
		err   error
	}{
		{"success", 1, "", "", nil},
		{"complementary=", 1, "Complementary", "", nil},
		{"firstname=", 1, "Firstname", "", errors.New("input:firstname")},
		{"lastname=", 1, "Lastname", "", errors.New("input:lastname")},
		{"street=", 1, "Street", "", errors.New("input:street")},
		{"city=", 1, "City", "", errors.New("input:city")},
		{"zipcode=", 1, "Zipcode", "", errors.New("input:zipcode")},
		{"phone=", 1, "Phone", "", errors.New("input:phone")},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			a := address

			if tt.field != "" {
				reflect.ValueOf(&a).Elem().FieldByName(tt.field).SetString(tt.value)
			}

			err := a.Validate(ctx)
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.err) {
				t.Fatalf("err = %v, want nil", err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/users.redis")

	var cases = []struct {
		name string
		uid  int
		err  error
	}{
		{"success", 1, nil},
		{"id=0", 0, errors.New("the user is not found")},
		{"id=99999", 99999, errors.New("the user is not found")},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FindByUID(ctx, tt.uid)
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.err) {
				t.Fatalf("err = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestIsAdmin(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/users.redis")

	var cases = []struct {
		name  string
		email string
		admin bool
	}{
		{"admin", "hello@artisons.me", true},
		{"not  admin", "arnaud@artisons.me", false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			b := IsAdmin(ctx, tt.email)
			if b != tt.admin {
				t.Fatalf("err = %v, want %v", b, tt.admin)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	ctx := tests.Context()

	tests.Del(ctx, "user")
	tests.ImportData(ctx, cur+"testdata/users.redis")

	var cases = []struct {
		name  string
		email string
		role  string
		total int
	}{
		{"email=arnaud@artisons.me", "arnaud@artisons.me", "", 1},
		{"role=user", "arnaud@artisons.me", "user", 1},
		{"email=idontexist@artisons.me", "idontexist@artisons.me", "", 0},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			q := Query{Email: tt.email}

			if tt.role != "" {
				q.Role = tt.role
			}

			u, err := Search(ctx, q, 0, 10)
			if err != nil {
				t.Fatalf(`err = %v, want nil`, err.Error())
			}

			if u.Total != tt.total {
				t.Fatalf(`total = %d, want %d`, u.Total, tt.total)
			}

			if len(u.Users) != tt.total {
				t.Fatalf(`len(articles) = %d, want %d`, len(u.Users), tt.total)
			}
		})
	}
}

func TestRefresh(t *testing.T) {
	ctx := tests.Context()
	u := User{SID: "123456789"}

	if err := u.RefreshSession(ctx); err != nil {
		t.Fatalf("err = %s, want nil", err.Error())
	}

	ttl, _ := db.Redis.TTL(ctx, "session:"+"123456789").Result()
	if ttl <= time.Minute {
		t.Fatalf(`ttl = %d want %d`, ttl, conf.SessionDuration)
	}
}
