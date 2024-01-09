package shops

import (
	"gifthub/tests"
	"testing"

	"github.com/go-faker/faker/v4"
)

var ra faker.RealAddress = faker.GetRealAddress()
var shop Settings = Settings{
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
	if err == nil || err.Error() != "input_name_invalid" {
		t.Fatalf("s.Validate(ctx, a) = '%v', want 'input_name_invalid'", err)
	}
}

func TestSaveReturnErrorWhenAddressIsEmpty(t *testing.T) {
	s := shop
	s.Address = ""

	ctx := tests.Context()
	err := s.Validate(ctx)
	if err == nil || err.Error() != "input_street_invalid" {
		t.Fatalf("s.Validate(ctx, a) = '%v', want 'input_street_invalid'", err)
	}
}

func TestSaveReturnErrorWhenCityIsEmpty(t *testing.T) {
	s := shop
	s.City = ""

	ctx := tests.Context()
	err := s.Validate(ctx)
	if err == nil || err.Error() != "input_city_invalid" {
		t.Fatalf("s.Validate(ctx, a) = '%v', want 'input_city_invalid'", err)
	}
}

func TestSaveReturnErrorWhenZipcodeIsEmpty(t *testing.T) {
	s := shop
	s.Zipcode = ""

	ctx := tests.Context()
	err := s.Validate(ctx)
	if err == nil || err.Error() != "input_zipcode_invalid" {
		t.Fatalf("s.Validate(ctx, a) = '%v', want 'input_zipcode_invalid'", err)
	}
}

func TestSaveReturnErrorWhenLogoIsEmpty(t *testing.T) {
	s := shop
	s.Logo = ""

	ctx := tests.Context()
	err := s.Validate(ctx)
	if err == nil || err.Error() != "input_logo_invalid" {
		t.Fatalf("s.Validate(ctx, a) = '%v', want 'input_logo_invalid'", err)
	}
}

func TestSaveReturnErrorWhenEmailIsEmpty(t *testing.T) {
	s := shop
	s.Email = ""

	ctx := tests.Context()
	err := s.Validate(ctx)
	if err == nil || err.Error() != "input_email_invalid" {
		t.Fatalf("s.Validate(ctx, a) = '%v', want 'input_email_invalid'", err)
	}
}

func TestSaveReturnErrorWhenEmailIsInvalid(t *testing.T) {
	s := shop
	s.Email = "hello"

	ctx := tests.Context()
	err := s.Validate(ctx)
	if err == nil || err.Error() != "input_email_invalid" {
		t.Fatalf("s.Validate(ctx, a) = '%v', want 'input_email_invalid'", err)
	}
}
