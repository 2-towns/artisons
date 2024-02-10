package users

import (
	"artisons/tests"
	"errors"
	"fmt"
	"testing"
)

func TestAddWPToken(t *testing.T) {
	ctx := tests.Context()

	var cases = []struct {
		name  string
		token string
		err   error
	}{
		{"success", "abc", nil},
		{"token=", "", errors.New("input:wptoken")},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if err := user.AddWPToken(ctx, tt.token); fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.err) {
				t.Fatalf("err = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestDeleteWPToken(t *testing.T) {
	ctx := tests.Context()

	tests.ImportData(ctx, cur+"testdata/users.redis")

	var cases = []struct {
		name string
		sid  string
		err  error
	}{
		{"success", "987654321", nil},
		{"success", "", errors.New("you are not authorized to process this request")},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if err := user.DeleteWPToken(ctx, tt.sid); fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.err) {
				t.Fatalf("err = %v, want %v", err, tt.err)
			}
		})
	}
}
