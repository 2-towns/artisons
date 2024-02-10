package seo

import (
	"artisons/tests"
	"reflect"
	"testing"
)

var content = Content{
	Key:         "test",
	URL:         "/test.html",
	Title:       "The social networks are evil.",
	Description: "Buh the social networks.",
}

func TestValidate(t *testing.T) {
	ctx := tests.Context()

	var tests = []struct{ name, field, value, want string }{
		{"key", "Key", "", "input:key"},
		{"title", "Title", "", "input:title"},
		{"url", "URL", "", "input:url"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := content

			reflect.ValueOf(&c).Elem().FieldByName(tt.field).SetString(tt.value)

			if err := c.Validate(ctx); err == nil || err.Error() != tt.want {
				t.Fatalf(`err = %v, want %s`, err, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	ctx := tests.Context()

	if _, err := Find(ctx, "idontexist"); err == nil || err.Error() != "the data is not found" {
		t.Fatalf(`err = %v, want the data is not found`, err)
	}
}
