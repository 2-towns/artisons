package shops

import (
	"artisons/tests"
	"fmt"
	"path"
	"reflect"
	"runtime"
	"testing"

	"github.com/go-faker/faker/v4"
)

var cur string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	cur = path.Dir(filename) + "/"
}

var ra faker.RealAddress = faker.GetRealAddress()
var shop Contact = Contact{
	Logo:    "../shops/123/1.jpeg",
	Banner1: "../shops/123/1.jpeg",
	Name:    faker.Name(),
	Email:   faker.Email(),
	Address: ra.Address,
	City:    ra.City,
	Zipcode: ra.PostalCode,
	Phone:   faker.Phonenumber(),
}

func TestSave(t *testing.T) {
	ctx := tests.Context()

	var tests = []struct{ name, field, value, want string }{
		{"name=", "Name", "", "input:name"},
		{"phone=", "Phone", "", "input:phone"},
		{"email=", "Email", "", "input:email"},
		{"email=hello", "Email", "hello", "input:email"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := shop

			reflect.ValueOf(&s).Elem().FieldByName(tt.field).SetString(tt.value)

			if err := s.Validate(ctx); err == nil || err.Error() != tt.want {
				t.Fatalf(`err = %v, want %s`, err, tt.want)
			}
		})
	}
}

func TestDeliveries(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/deliveries.redis")

	del, err := Deliveries(ctx)
	if err != nil || len(del) == 0 {
		t.Fatalf(`err = %v, want nil`, err)
	}
}

func TestPayments(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/payments.redis")

	del, err := Payments(ctx)
	if err != nil || len(del) == 0 {
		t.Fatalf(`err = %v, want nil`, err)
	}
}

func TestPay(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/payments.redis")

	var tests = []struct {
		name    string
		oid     string
		payment string
		err     error
	}{
		{"cash", "ORD1", "cash", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Pay(ctx, tt.oid, tt.payment); fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.err) {
				t.Fatalf(`err = %v, want %v`, err, tt.err)
			}
		})
	}
}

func TestIsValidDelivery(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/deliveries.redis")

	var tests = []struct {
		name  string
		value string
		valid bool
	}{
		{"collect", "collect", true},
		{"idontexist", "idontexist", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if valid := IsValidDelivery(ctx, tt.value); valid != tt.valid {
				t.Fatalf(`valid = %v, want %v`, valid, tt.valid)
			}
		})
	}
}

func TestIsValidPayment(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/payments.redis")

	var tests = []struct {
		name  string
		value string
		valid bool
	}{
		{"cash", "cash", true},
		{"idontexist", "idontexist", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if valid := IsValidPayment(ctx, tt.value); valid != tt.valid {
				t.Fatalf(`valid = %v, want %v`, valid, tt.valid)
			}
		})
	}
}
