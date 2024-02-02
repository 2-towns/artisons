package htmx

import (
	"artisons/http/contexts"
	"context"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contexts.HX, r.Header.Get("HX-Request") == "true")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
