package db

import (
	"context"
	"testing"
)

func TestReturnsEncodedStringWhenUsingQuote(t *testing.T) {
	exp := "c\\'est"
	arg := "c'est"
	if s := Escape(arg); s != exp {
		t.Fatalf(`Escape('%s') = '%s', want '%s'`, arg, s, exp)
	}
}

func TestUnescapeWhenContainsQuote(t *testing.T) {
	exp := "c'est"
	arg := "c\\'est"

	if s := Unescape(arg); s != exp {
		t.Fatalf(`Unescape('%s') = '%s', want '%s'`, arg, s, exp)
	}
}

func TestSearchValueWhenDoesNotContainSpace(t *testing.T) {
	exp := "hello"
	arg := "hello"

	if s := SearchValue(arg); s != exp {
		t.Fatalf(`Unescape('%s') = '%s', want '%s'`, arg, s, exp)
	}
}

func TestSearchValueWhenDoesContainSpace(t *testing.T) {
	exp := "hello|world"
	arg := "hello world"

	if s := SearchValue(arg); s != exp {
		t.Fatalf(`Unescape('%s') = '%s', want '%s'`, arg, s, exp)
	}
}

func TestSplitQueryReturnsStringSplitted(t *testing.T) {
	ctx := context.Background()
	s := "hello \"hello with space\""

	if args, err := SplitQuery(ctx, s); err != nil || len(args) != 2 {
		t.Fatalf(`SplitQuery(ctx, s) = '%v' %v, want [hello "hello with space"], nil`, args, err)
	}

}
