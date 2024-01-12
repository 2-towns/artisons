package locales

import (
	"gifthub/tests"
	"testing"
)

var value = Value{
	Locale: "en",
	Key:    "test",
	Value:  "coucou",
}

func TestValidateReturnsErrorWhenKeyIsEmpty(t *testing.T) {
	c := tests.Context()

	v := value
	v.Key = ""

	if err := v.Validate(c); err == nil || err.Error() != "input_key_invalid" {
		t.Fatalf(`v.Validate(c) = %v, want not "input_key_invalid"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenValueIsEmpty(t *testing.T) {
	c := tests.Context()

	v := value
	v.Value = ""

	if err := v.Validate(c); err == nil || err.Error() != "input_value_invalid" {
		t.Fatalf(`v.Validate(c) = %v, want not "input_value_invalid"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenLocaleIsEmpty(t *testing.T) {
	c := tests.Context()

	v := value
	v.Locale = ""

	if err := v.Validate(c); err == nil || err.Error() != "input_locale_invalid" {
		t.Fatalf(`v.Validate(c) = %v, want not "input_locale_invalid"`, err.Error())
	}
}

func TestValidateReturnsErrorWhenLocaleIsInvalid(t *testing.T) {
	c := tests.Context()

	v := value
	v.Locale = "!!!"

	if err := v.Validate(c); err == nil || err.Error() != "input_locale_invalid" {
		t.Fatalf(`v.Validate(c) = %v, want not "input_locale_invalid"`, err.Error())
	}
}

func TestSaveReturnsNoError(t *testing.T) {
	c := tests.Context()

	v := value

	if err := v.Save(c); err != nil {
		t.Fatalf(`v.Validate(c) = %s, want not nil`, err.Error())
	}
}
