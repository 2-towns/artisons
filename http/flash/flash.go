package flash

import (
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/cookies"
	"context"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		c, err := r.Cookie(cookies.FlashMessage)
		if err == nil && c != nil {
			flash := c.Value

			cookie := &http.Cookie{
				Name:     cookies.FlashMessage,
				Value:    flash,
				MaxAge:   -1,
				Path:     "/",
				HttpOnly: true,
				Secure:   conf.Cookie.Secure,
				Domain:   conf.Cookie.Domain,
				SameSite: http.SameSiteStrictMode,
			}

			http.SetCookie(w, cookie)

			ctx = context.WithValue(ctx, contexts.Flash, flash)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
