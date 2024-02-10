package shops

import (
	"artisons/tests"
	"reflect"
	"testing"

	"github.com/go-faker/faker/v4"
)

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

	del, err := Deliveries(ctx)
	if err != nil || len(del) == 0 {
		t.Fatalf(`err = %v, want nil`, err)
	}
}
