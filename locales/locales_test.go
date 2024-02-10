package locales

import (
	"artisons/tests"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

var value = Value{
	Key:   "test",
	Value: "coucou",
}

func TestValidate(t *testing.T) {
	ctx := tests.Context()

	var tests = []struct {
		name  string
		field string
		value string
		err   error
	}{
		{"key", "Key", "", errors.New("input:key")},
		{"value", "Value", "", errors.New("input:value")},
		{"success", "", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := value

			if tt.field != "" {
				reflect.ValueOf(&v).Elem().FieldByName(tt.field).SetString(tt.value)
			}

			if err := v.Validate(ctx); fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tt.err) {
				t.Fatalf(`err = %v, want %s`, err, tt.err)
			}
		})
	}
}
