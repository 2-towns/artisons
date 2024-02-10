package validators

import (
	"testing"
)

func TestValidate(t *testing.T) {

	var tests = []struct {
		name  string
		value string
		b     bool
	}{
		{"hello ", "hello ", false},
		{"hello(", "hello(", false},
		{"hello)", "hello)", false},
		{"Hello", "Hello", false},
		{"Hello1", "Hello1", false},
		{"hello-", "hello-", false},
		{"hello!", "hello!", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := V.Var(tt.value, "title"); (!tt.b && err != nil) || (tt.b && err == nil) {
				t.Fatalf(`err = %v`, err.Error())
			}
		})
	}
}
