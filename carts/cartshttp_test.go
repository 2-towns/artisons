package carts

import (
	"artisons/conf"
	"artisons/http/cookies"
	"artisons/string/stringutil"
	"artisons/tests"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareRefreshCartIdWhenExisting(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	cid, err := stringutil.Random()
	if err != nil {
		t.Fatalf(`err = %s, want nil`, err.Error())
	}

	cookie := &http.Cookie{
		Name:     cookies.CartID,
		Value:    cid,
		MaxAge:   int(conf.Cookie.MaxAge),
		Path:     "/",
		HttpOnly: true,
		Secure:   conf.Cookie.Secure,
		Domain:   conf.Cookie.Domain,
		SameSite: http.SameSiteStrictMode,
	}

	req.AddCookie(cookie)

	ctx := tests.Context()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()

	handler := Middleware(testHandler)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf(`status = %d, want %d`, status, http.StatusOK)
	}

	cks := rr.Result().Cookies()
	if len(cks) != 1 {
		t.Fatalf(`len(cookies) = %d, want 1`, len(cks))
	}

	c := cks[0]

	if c.Name != cookies.CartID {
		t.Fatalf(`c.Name = %s, want %s`, c.Name, cookies.CartID)
	}

	if c.Value != cid {
		t.Fatalf(`c.Value = %s, want %s`, c.Value, cid)
	}
}
