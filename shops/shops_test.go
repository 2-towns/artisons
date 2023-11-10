package shops

import (
	"gifthub/tests"
	"gifthub/users"
	"testing"

	"github.com/go-faker/faker/v4"
)

var ra faker.RealAddress = faker.GetRealAddress()
var shop Shop = Shop{
	Logo: "../web/images/123/1",
	Address: users.Address{
		Lastname:      faker.Name(),
		Firstname:     faker.Name(),
		Street:        ra.Address,
		City:          ra.City,
		Complementary: ra.Address,
		Zipcode:       ra.PostalCode,
		Phone:         faker.Phonenumber(),
	},
}

func TestSaveReturnErrorWhenFirstnameIsEmpty(t *testing.T) {
	s := shop
	s.Address.Firstname = ""

	ctx := tests.Context()
	err := s.Save(ctx)
	if err == nil || err.Error() != "input_firstname_required" {
		t.Fatalf("s.Save(ctx, a) = '%v', want 'input_firstname_required'", err)
	}
}

func TestSaveReturnErrorWhenStreetIsEmpty(t *testing.T) {
	s := shop
	s.Address.Street = ""

	ctx := tests.Context()
	err := s.Save(ctx)
	if err == nil || err.Error() != "input_street_required" {
		t.Fatalf("s.Save(ctx, a) = '%v', want 'input_street_required'", err)
	}
}

func TestSaveReturnErrorWhenCityIsEmpty(t *testing.T) {
	s := shop
	s.Address.City = ""

	ctx := tests.Context()
	err := s.Save(ctx)
	if err == nil || err.Error() != "input_city_required" {
		t.Fatalf("s.Save(ctx, a) = '%v', want 'input_city_required'", err)
	}
}

func TestSaveReturnErrorWhenZipcodeIsEmpty(t *testing.T) {
	s := shop
	s.Address.Zipcode = ""

	ctx := tests.Context()
	err := s.Save(ctx)
	if err == nil || err.Error() != "input_zipcode_required" {
		t.Fatalf("s.Save(ctx, a) = '%v', want 'input_zipcode_required'", err)
	}
}

func TestSaveReturnErrorWhenLogoIsEmpty(t *testing.T) {
	s := shop
	s.Logo = ""

	ctx := tests.Context()
	err := s.Save(ctx)
	if err == nil || err.Error() != "input_logo_required" {
		t.Fatalf("s.Save(ctx, a) = '%v', want 'input_logo_required'", err)
	}
}

func TestGetReturnShopInfoErrorWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	shop, err := Get(ctx)
	if err != nil {
		t.Fatalf("Get(ctx) = %v, want nil", err)
	}

	if shop.Logo == "" {
		t.Fatalf("shop.Logo = '%v', want not empty", err)
	}

	if shop.Address.City == "" {
		t.Fatalf("shop.Address.City = '%v', want not empty", err)
	}

	if shop.Address.Firstname == "" {
		t.Fatalf("shop.Address.Firstname = '%v', want not empty", err)
	}

	if shop.Address.Lastname != "None" {
		t.Fatalf("shop.Address.Lastname = '%v', want not empty", err)
	}

	if shop.Address.Phone == "" {
		t.Fatalf("shop.Address.Phone = '%v', want not empty", err)
	}

	if shop.Address.Street == "" {
		t.Fatalf("shop.Address.Street = '%v', want not empty", err)
	}

	if shop.Address.Zipcode == "" {
		t.Fatalf("shop.Address.Zipcode = '%v', want not empty", err)
	}
}
