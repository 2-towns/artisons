package security

import (
	"artisons/tests"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCsrf(t *testing.T) {
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
