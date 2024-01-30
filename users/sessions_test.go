package users

import (
	"artisons/tests"
	"testing"
)

func TestSessionsReturnsSessionsWhenUserHasSession(t *testing.T) {
	ctx := tests.Context()
	sessions, err := user.Sessions(ctx)
	if len(sessions) == 0 || err != nil {
		t.Fatalf("user.Session(ctx) = %v, %v, want not empty, nil", sessions, err)
	}

	session := sessions[0]
	if session.ID == "" || session.Device == "" || session.TTL == 0 {
		t.Fatalf("sessions[0] = %v, want Session", session)
	}
}

func TestSessionsReturnsEmptySliceWhenUserDoesNotHaveSession(t *testing.T) {
	ctx := tests.Context()
	sessions, err := User{ID: 98}.Sessions(ctx)
	if len(sessions) != 0 || err != nil {
		t.Fatalf("User{ID: 98}.Session(ctx) = %v, %v, want []Session, nil", sessions, err)
	}
}
