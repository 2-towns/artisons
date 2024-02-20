package addresses

import (
	"artisons/tests"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-faker/faker/v4"
)

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

func TestSaveAddress(t *testing.T) {
	ctx := tests.Context()
	err := address.Save(ctx, "user:1")

	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}

	err = address.Save(ctx, "")

	if err == nil {
		t.Fatal("err = nil, want something went wrong")
	}
}

func TestValidateAddress(t *testing.T) {
	ctx := tests.Context()

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
