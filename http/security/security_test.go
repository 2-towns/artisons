package security

import (
	"gifthub/tests"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareSetHeadersWhenRequestIsComing(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := tests.Context()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	handler := Headers(testHandler)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf(`status = %d, want %d`, status, http.StatusOK)
	}

	if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf(`X-Content-Type-Options = %s, want 'nosniff'`, rr.Header().Get("X-Content-Type-Options"))
	}

	if rr.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatalf(`Access-Control-Allow-Credentials = %s, want 'nosniff'`, rr.Header().Get("Access-Control-Allow-Credentials"))
	}

	if rr.Header().Get("Referrer-Policy") != "strict-origin" {
		t.Fatalf(`Referrer-Policy = %s, want 'strict-origin'`, rr.Header().Get("Referrer-Policy"))
	}

	if rr.Header().Get("Strict-Transport-Security") != "max-age=63072000; includeSubDomains; preload" {
		t.Fatalf(`Strict-Transport-Security = %s, want 'max-age=63072000; includeSubDomains; preload'`, rr.Header().Get("Strict-Transport-Security"))
	}

	if rr.Header().Get("X-XSS-Protection") != "1" {
		t.Fatalf(`X-XSS-Protection = %s, want "1"`, rr.Header().Get("X-XSS-Protection"))
	}

	if rr.Header().Get("Content-Security-Policy") != "default-src 'self'" {
		t.Fatalf(`Content-Security-Policy = %s, want "default-src 'self'"`, rr.Header().Get("Content-Security-Policy"))
	}
}

func TestCsrfReturns200WhenGetMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := tests.Context()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	handler := Csrf(testHandler)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf(`status = %d, want %d`, status, http.StatusOK)
	}
}

func TestCsrfReturns400WhenPostMethodWithoutHxHeader(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := tests.Context()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	handler := Csrf(testHandler)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf(`status = %d, want %d`, status, http.StatusBadRequest)
	}
}

func TestCsrfReturns200WhenPostMethodWithHxHeader(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("HX-Request", "true")

	ctx := tests.Context()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	handler := Csrf(testHandler)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf(`status = %d, want %d`, status, http.StatusOK)
	}
}
