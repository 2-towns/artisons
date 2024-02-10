package users

import (
	"artisons/tests"
	"testing"
)

func TestSessions(t *testing.T) {
	ctx := tests.Context()

	tests.Del(ctx, "session")
	tests.ImportData(ctx, cur+"testdata/sessions.redis")

	var tests = []struct {
		name     string
		id       int
		sessions int
	}{
		{"id=1", 1, 1},
		{"id=2", 2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{ID: tt.id}

			sessions, err := u.Sessions(ctx)

			if err != nil {
				t.Fatalf("err =  %v, want nil", err)
			}

			if len(sessions) != tt.sessions {
				t.Fatalf("len(sessions) = %d , want %d", len(sessions), tt.sessions)
			}
		})
	}
}
