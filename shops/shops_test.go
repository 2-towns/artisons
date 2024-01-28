package shops

import (
	"artisons/tests"
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

func TestSaveReturnErrorWhenNameIsEmpty(t *testing.T) {
	s := shop
	s.Name = ""

	ctx := tests.Context()
	err := s.Validate(ctx)
	if err == nil || err.Error() != "input:name" {
		t.Fatalf("s.Validate(ctx, a) = '%v', want 'input:name'", err)
	}
}

func TestSaveReturnErrorWhenPhoneIsEmpty(t *testing.T) {
	s := shop
	s.Phone = ""

	ctx := tests.Context()
	err := s.Validate(ctx)
	if err == nil || err.Error() != "input:phone" {
		t.Fatalf("s.Validate(ctx, a) = '%v', want 'input:phone'", err)
	}
}

func TestSaveReturnErrorWhenEmailIsEmpty(t *testing.T) {
	s := shop
	s.Email = ""

	ctx := tests.Context()
	err := s.Validate(ctx)
	if err == nil || err.Error() != "input:email" {
		t.Fatalf("s.Validate(ctx, a) = '%v', want 'input:email'", err)
	}
}

func TestSaveReturnErrorWhenEmailIsInvalid(t *testing.T) {
	s := shop
	s.Email = "hello"

	ctx := tests.Context()
	err := s.Validate(ctx)
	if err == nil || err.Error() != "input:email" {
		t.Fatalf("s.Validate(ctx, a) = '%v', want 'input:email'", err)
	}
}
