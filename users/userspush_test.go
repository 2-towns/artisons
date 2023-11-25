package users

import (
	"gifthub/tests"
	"testing"
)

func TestAddWPTokenReturnsNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	if err := user.AddWPToken(ctx, "{}"); err != nil {
		t.Fatalf("user.AddWPToken('{}') = %v, want nil", err)
	}
}

func TestAddWPTokenReturnsErrorWhenTokenIsEmpty(t *testing.T) {
	ctx := tests.Context()
	if err := user.AddWPToken(ctx, ""); err == nil || err.Error() != "input_wptoken_required" {
		t.Fatalf("user.AddWPToken('') = %v, want input_wptoken_required", err)
	}
}

func TestATestDeleteWPTokenReturnsNilWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	if err := user.DeleteWPToken(ctx, user.SID); err != nil {
		t.Fatalf("user.TestDeleteWPToken(user.SID) = %v, want nil", err)
	}
}

func TestATestDeleteWPTokenReturnsErrorWhenTokenIsEmpty(t *testing.T) {
	ctx := tests.Context()
	if err := user.DeleteWPToken(ctx, ""); err == nil || err.Error() != "error_http_unauthorized" {
		t.Fatalf("u.DeleteWPToken('') = %v, want unauthorized", err)
	}
}
