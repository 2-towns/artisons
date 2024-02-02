package users

import (
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/cookies"
	"artisons/tests"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindBySessionIDReturnsSessionWhenSuccess(t *testing.T) {
	ctx := tests.Context()
	u, err := findBySessionID(ctx, user.SID)

	if err != nil || u.SID == "" || u.Email == "" {
		t.Fatalf(`findBySessionID(ctx, user.SID) = %v, %v, want User, nil`, u, err)
	}
}

func TestFindBySessionIDReturnsErrorWhenSidIsEmpty(t *testing.T) {
	ctx := tests.Context()
	u, err := findBySessionID(ctx, "")

	if err == nil || err.Error() != "you are not authorized to process this request" || u.Email != "" {
		t.Fatalf("findBySessionID('') = %v, %v, want User{}, 'unauthorized'", u, err)
	}
}

func TestFindBySessionIDReturnsErrorWhenSessionIsExpired(t *testing.T) {
	ctx := tests.Context()
	u, err := findBySessionID(ctx, "expired")

	if err == nil || err.Error() != "you are not authorized to process this request" || u.Email != "" {
		t.Fatalf(`findBySessionID("expired") = %v, %v, want User, nil`, u, err)
	}
}

func TestMiddlewareDestroySessionIdWhenItIsNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	cookie := &http.Cookie{
		Name:     cookies.SessionID,
		Value:    "hello",
		MaxAge:   int(conf.Cookie.MaxAge),
		Path:     "/",
		HttpOnly: true,
		Secure:   conf.Cookie.Secure,
		Domain:   conf.Cookie.Domain,
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
	var c *http.Cookie
	for _, val := range cks {
		if val.Name == cookies.SessionID {
			c = val
		}
	}

	if c.Name == "" {
		t.Fatalf(`c.Name is empty, want 'wsid'`)
	}

	if c.Name != cookies.SessionID {
		t.Fatalf(`c.Name = %s, want %s`, c.Name, cookies.SessionID)
	}

	if c.Value != "hello" {
		t.Fatalf(`c.Value = %s, want 'hello'`, c.Value)
	}

	if c.Secure != conf.Cookie.Secure {
		t.Fatalf(`c.Secure = false, want %v`, conf.Cookie.Secure)
	}

	if c.HttpOnly == false {
		t.Fatalf(`c.HttpOnly = false, want true`)
	}

	if c.Path != "/" {
		t.Fatalf(`c.Path = %s, want '/'`, c.Path)
	}

	if c.MaxAge > 0 {
		t.Fatalf(`c.MaxAge = %d, want 0`, c.MaxAge)
	}
}

func TestMiddlewareRefreshesSessionIdWhenItExists(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	cookie := &http.Cookie{
		Name:     cookies.SessionID,
		Value:    user.SID,
		MaxAge:   int(conf.Cookie.MaxAge),
		Path:     "/",
		HttpOnly: true,
		Secure:   conf.Cookie.Secure,
		Domain:   conf.Cookie.Domain,
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
	if len(cks) != 2 {
		t.Fatalf(`len(cookies) = %d, want 2`, len(cks))
	}

	c := cks[1]

	if c.Name != cookies.SessionID {
		t.Fatalf(`c.Name = %s, want %s`, c.Name, cookies.SessionID)
	}

	if c.Value != user.SID {
		t.Fatalf(`c.Value = %s, want '%s'`, c.Value, user.SID)
	}

	if c.Secure != conf.Cookie.Secure {
		t.Fatalf(`c.Secure = false, want %v`, conf.Cookie.Secure)
	}

	if c.HttpOnly == false {
		t.Fatalf(`c.HttpOnly = false, want true`)
	}

	if c.Path != "/" {
		t.Fatalf(`c.Path = %s, want '/'`, c.Path)
	}

	if c.MaxAge == 0 {
		t.Fatalf(`c.MaxAge = %d, want 0`, c.MaxAge)
	}
}

func TestAdminOnlyRedirectWhenNoUserInContext(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := tests.Context()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()

	handler := AdminOnly(testHandler)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusFound {
		t.Fatalf(`status = %d, want %d`, status, http.StatusFound)
	}

	if rr.Header().Get("Location") != "/sso.html" {
		t.Fatalf(`Location = %s, want %s`, rr.Header().Get("Location"), "/sso.html")
	}
}

func TestAdminOnlyRedirectWhenUserIsNotAdmin(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.User, User{Role: "user"})

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()

	handler := AdminOnly(testHandler)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Fatalf(`status = %d, want %d`, status, http.StatusUnauthorized)
	}
}

func TestAdminOnlyContinuesWhenUserIsAdmin(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := tests.Context()
	ctx = context.WithValue(ctx, contexts.User, User{Role: "admin"})

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()

	handler := AdminOnly(testHandler)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf(`status = %d, want %d`, status, http.StatusOK)
	}
}
