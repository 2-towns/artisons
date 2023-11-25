package security

import (
	"gifthub/http/httperrors"
	"log/slog"
	"net/http"
	"time"
)

func Csrf(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.Header.Get("HX-Request") != "true" {
			slog.Info(r.Method + " " + r.Header.Get("HX-Request"))
			httperrors.Page(w, r.Context(), "error_http_csrf", 400)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Powered-By", "WordPress")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set(
			"Accept-CH",
			"Sec-CH-Prefers-Color-Scheme, Device-Memory, Downlink, ECT",
		)
		w.Header().Set("Referrer-Policy", "strict-origin")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set(
			"Strict-Transport-Security",
			"max-age=63072000; includeSubDomains; preload",
		)
		w.Header().Set("X-XSS-Protection", "1")
		w.Header().Set("Date", time.Now().Format(time.RFC1123))
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		next.ServeHTTP(w, r)
	})

}
