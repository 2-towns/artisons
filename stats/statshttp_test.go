package stats

import (
	"gifthub/db"
	"gifthub/tests"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMiddlewareDoestNotUpdateStatWhenAdmin(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/admin/index.html", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := tests.Context()
	now := time.Now().Format("20060102")
	pageviews, _ := db.Redis.ZScore(ctx, "stats:pageviews:"+now, "/index.html").Result()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := Middleware(testHandler)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	time.Sleep(time.Millisecond * 10)

	p, _ := db.Redis.ZScore(ctx, "stats:pageviews:"+now, "/index.html").Result()
	if pageviews != p {
		t.Fatalf(`pageviews = %f, want %f`, p, pageviews)
	}
}
