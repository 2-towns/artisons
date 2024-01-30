package httpext

import (
	"artisons/conf"
	"net/http"
	"testing"
)

func TestPaginationReturnsData(t *testing.T) {
	req, err := http.NewRequest("GET", "/?q=hello&page=2", nil)
	if err != nil {
		t.Fatal(err)
	}

	p := Pagination(req)
	if p.Offset != conf.ItemsPerPage {
		t.Fatalf(`offset = %d, want %d`, p.Offset, conf.ItemsPerPage)
	}

	if p.Num != conf.ItemsPerPage*2 {
		t.Fatalf(`offset = %d, want %d`, p.Offset, conf.ItemsPerPage*2)
	}

	if p.Query != "hello" {
		t.Fatalf(`q = %s, want hello`, p.Query)
	}
}

func TestPaginationReturnsZeroWhenPageIsMissing(t *testing.T) {
	req, err := http.NewRequest("GET", "/?q=id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	p := Pagination(req)
	if p.Offset != 0 {
		t.Fatalf(`offset = %d, want 0`, p.Offset)
	}

	if p.Num != conf.ItemsPerPage {
		t.Fatalf(`offset = %d, want %d`, p.Offset, conf.ItemsPerPage*2)
	}
}

func TestPaginationReturnsZeroWhenPageIsInvalid(t *testing.T) {
	req, err := http.NewRequest("GET", "/?q=id=1&page=hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	p := Pagination(req)
	if p.Offset != 0 {
		t.Fatalf(`offset = %d, want 0`, p.Offset)
	}

	if p.Num != conf.ItemsPerPage {
		t.Fatalf(`offset = %d, want %d`, p.Offset, conf.ItemsPerPage*2)
	}
}
