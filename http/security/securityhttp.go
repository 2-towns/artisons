package security

import (
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"log/slog"
	"net/http"
)

func Csrf(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		isHX, _ := ctx.Value(contexts.HX).(bool)

		if !isHX {
			slog.Info(r.Method + " " + r.Header.Get("HX-Request"))
			httperrors.Page(w, r.Context(), "you are not authorized to process this request", 400)
			return
		}

		next.ServeHTTP(w, r)
	})
}
