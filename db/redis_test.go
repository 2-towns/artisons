package db

import (
	"testing"
)

func TestReturnsEncodedStringWhenSpecialCharacters(t *testing.T) {
	exp := "s\\,\\.\\<\\>\\{\\}\\[\\]\\\"\\:\\;\\!\\@\\#\\$\\%\\^\\&\\*\\(\\)\\-\\+\\=\\~"
	arg := "s,.<>{}[]\":;!@#$%^&*()-+=~"
	if s := Escape(arg); s != exp {
		t.Fatalf(`EscapeSearchQuery('%s') = '%s', want '%s'`, arg, s, exp)
	}
}

func TestUnescapeWhenContainsEscapedContent(t *testing.T) {
	exp := "s,.<>{}[]\":;!@#$%^&*()-+=~"
	arg := "s\\,\\.\\<\\>\\{\\}\\[\\]\\\\\"\\:\\;\\!\\@\\#\\$\\%\\^\\&\\*\\(\\)\\-\\+\\=\\~"

	if s := Unescape(arg); s != exp {
		t.Fatalf(`Unescape('%s') = '%s', want '%s'`, arg, s, exp)
	}
}
