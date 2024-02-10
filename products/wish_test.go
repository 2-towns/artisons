package products

import (
	"artisons/tests"
	"slices"
	"testing"
)

func TestWish(t *testing.T) {
	ctx := tests.Context()
	tests.ImportData(ctx, cur+"testdata/wish.redis")

	if err := Wish(ctx, 1, "PDT1"); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}
}

func TestUnWish(t *testing.T) {
	ctx := tests.Context()
	tests.ImportData(ctx, cur+"testdata/wish.redis")

	if err := UnWish(ctx, 1, "PDT1"); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}
}

func TestWishes(t *testing.T) {
	ctx := tests.Context()
	tests.ImportData(ctx, cur+"testdata/wish.redis")

	wishes, err := Wishes(ctx, 1)

	if err != nil {
		t.Fatalf(`err = %v, %v, want nil`, wishes, err)
	}

	if len(wishes) == 0 {
		t.Fatalf(`len(wishes) = %d, want 0`, len(wishes))
	}

	if !slices.Contains(wishes, "PDT2") {
		t.Fatalf(`slices contains "PDT2" = %v, want true`, slices.Contains(wishes, "PDT2"))
	}
}

func TestHasWish(t *testing.T) {
	ctx := tests.Context()
	tests.ImportData(ctx, cur+"testdata/wish.redis")

	var cases = []struct {
		name string
		pid  string
		b    bool
	}{
		{"success", "PDT1", true},
		{"pid=idontexist", "idontexist", false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if has := HasWish(ctx, 1, tt.pid); has != tt.b {
				t.Fatalf(`has = %v, want %v`, has, tt.b)
			}
		})
	}
}
