package db

import (
	"context"
	"testing"
)

func TestEscape(t *testing.T) {
	exp := "c\\'est"
	arg := "c'est"
	if s := Escape(arg); s != exp {
		t.Fatalf(`s = %s, want c\\'est`, s)
	}
}

func TestUnescape(t *testing.T) {
	exp := "c'est"
	arg := "c\\'est"

	if s := Unescape(arg); s != exp {
		t.Fatalf(`s = %s, want c'est`, s)
	}
}

func TestSearchValue(t *testing.T) {
	exp := "hello"
	arg := "hello"

	t.Run("hello", func(t *testing.T) {
		if s := SearchValue(arg); s != exp {
			t.Fatalf(`s = %s, want hello`, s)
		}
	})

	t.Run("hello world", func(t *testing.T) {
		exp := "hello|world"
		arg := "hello world"

		if s := SearchValue(arg); s != exp {
			t.Fatalf(`s = %s, want hello|world`, s)
		}
	})

}

func TestSplitQuery(t *testing.T) {
	ctx := context.Background()
	s := "hello \"hello with space\""

	if args, err := SplitQuery(ctx, s); err != nil || len(args) != 2 {
		t.Fatalf(`len(args) = %d, %v, want 2, nil`, len(args), err)
	}

}
