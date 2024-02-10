package cache

import "testing"

func TestBusting(t *testing.T) {
	Busting()

	if Buster("admin.js") == "" {
		t.Fatal("cachebuster is empty from admin.js")
	}

	if Buster("admin.css") == "" {
		t.Fatal("cachebuster is empty from admin.css")
	}
}
