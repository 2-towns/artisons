package validators

import (
	"testing"
)

func TestValidatesReturnsNilTitleWhenVarContainsASpace(t *testing.T) {
	if err := V.Var("hello ", "title"); err != nil {
		t.Fatalf(`V.Var("hello ", "title") = %v, want nil`, err.Error())
	}
}

func TestValidatesReturnsNilTitleWhenVarContainsAOpeningBrace(t *testing.T) {
	if err := V.Var("hello(", "title"); err != nil {
		t.Fatalf(`V.Var("hello(", "title") = %v, want nil`, err.Error())
	}
}

func TestValidatesReturnsNilTitleWhenVarContainsAClosingBrace(t *testing.T) {
	if err := V.Var("hello)", "title"); err != nil {
		t.Fatalf(`V.Var("hello)", "title") = %v, want nil`, err.Error())
	}
}

func TestValidatesReturnsNilTitleWhenVarContainsUppercase(t *testing.T) {
	if err := V.Var("Hello", "title"); err != nil {
		t.Fatalf(`V.Var("Hello", "title") = %v, want nil`, err.Error())
	}
}

func TestValidatesReturnsNilTitleWhenVarContainsNumeric(t *testing.T) {
	if err := V.Var("Hello1", "title"); err != nil {
		t.Fatalf(`V.Var("hello1", "title") = %v, want nil`, err.Error())
	}
}

func TestValidatesReturnsNilTitleWhenVarContainsDash(t *testing.T) {
	if err := V.Var("hello-", "title"); err != nil {
		t.Fatal(`V.Var("hello-", "title") = nil, want nil`)
	}
}

func TestValidatesReturnsErrorTitleWhenVarContainsSpecialCharacters(t *testing.T) {
	if err := V.Var("hello!", "title"); err == nil {
		t.Fatal(`V.Var("hello!", "title") = nil, want not nil`)
	}
}
