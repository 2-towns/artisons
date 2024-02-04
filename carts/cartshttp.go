package carts

import (
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/cookies"
	"artisons/http/httperrors"
	"context"
	"log/slog"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cid, err := r.Cookie(cookies.CartID)
		if err != nil {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		cookie := &http.Cookie{
			Name:     cookies.CartID,
			Value:    cid.Value,
			MaxAge:   int(conf.Cookie.MaxAge),
			Path:     "/",
			HttpOnly: true,
			Secure:   conf.Cookie.Secure,
			Domain:   conf.Cookie.Domain,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, cookie)

		ctx = context.WithValue(ctx, contexts.Cart, cid.Value)

		slog.LogAttrs(ctx, slog.LevelInfo, "cart id detected", slog.String("cid", cid.Value))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Redirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cart, err := Get(ctx)
		if err != nil {
			httperrors.Catch(w, ctx, err.Error(), 500)
			return
		}

		if len(cart.Products) == 0 {
			slog.LogAttrs(ctx, slog.LevelInfo, "the cart is empty")
			http.Redirect(w, r, "/cart.html", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)

	})
}
