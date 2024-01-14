package shops

import (
	"gifthub/tests"
	"testing"

	"github.com/go-faker/faker/v4"
)

var ra faker.RealAddress = faker.GetRealAddress()
var shop Contact = Contact{
	Logo:    "../web/images/123/1",
	Banner1: "../web/images/123/1",
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

func TestSaveReturnErrorWhenAddressIsEmpty(t *testing.T) {
	s := shop
	s.Address = ""

	ctx := tests.Context()
	err := s.Validate(ctx)
	if err == nil || err.Error() != "input:address" {
		t.Fatalf("s.Validate(ctx, a) = '%v', want 'input:address'", err)
	}
}

func TestSaveReturnErrorWhenCityIsEmpty(t *testing.T) {
	s := shop
	s.City = ""

	ctx := tests.Context()
	err := s.Validate(ctx)
	if err == nil || err.Error() != "input:city" {
		t.Fatalf("s.Validate(ctx, a) = '%v', want 'input:city'", err)
	}
}

func TestSaveReturnErrorWhenZipcodeIsEmpty(t *testing.T) {
	s := shop
	s.Zipcode = ""

	ctx := tests.Context()
	err := s.Validate(ctx)
	if err == nil || err.Error() != "input:zipcode" {
		t.Fatalf("s.Validate(ctx, a) = '%v', want 'input:zipcode'", err)
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
