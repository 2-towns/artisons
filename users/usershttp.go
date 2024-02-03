package users

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/http/cookies"
	"artisons/http/httperrors"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

func findBySessionID(ctx context.Context, sid string) (User, error) {
	l := slog.With(slog.String("sid", sid))
	l.LogAttrs(ctx, slog.LevelInfo, "finding the user")

	if sid == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate the session id")
		return User{}, errors.New("you are not authorized to process this request")
	}

	id, err := db.Redis.HGet(ctx, "session:"+sid, "uid").Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the auth id from redis", slog.String("error", err.Error()))
		return User{}, errors.New("you are not authorized to process this request")
	}

	m, err := db.Redis.HGetAll(ctx, "user:"+id).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the session from redis", slog.String("error", err.Error()))
		return User{}, errors.New("you are not authorized to process this request")
	}

	m["sid"] = sid
	u, err := parse(ctx, m)
	if err != nil {
		return User{}, errors.New("you are not authorized to process this request")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "user found", slog.Int("user_id", u.ID))

	return u, err
}

// Middleware detects the session ID in the cookies.
// If the session ID exists, it will load the current
// user into the context.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sid, err := r.Cookie(cookies.SessionID)
		if err != nil {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		user, err := findBySessionID(ctx, sid.Value)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelInfo, "session id not found so destroying it", slog.String("sid", sid.Value))

			cookie := &http.Cookie{
				Name:     cookies.SessionID,
				Value:    sid.Value,
				MaxAge:   -1,
				Path:     "/",
				HttpOnly: true,
				Secure:   conf.Cookie.Secure,
				Domain:   conf.Cookie.Domain,
				SameSite: http.SameSiteStrictMode,
			}
			http.SetCookie(w, cookie)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		} else {
			err := user.RefreshSession(ctx)
			if err != nil {
				httperrors.Page(w, r.Context(), err.Error(), 500)
				return
			}

			cookie := &http.Cookie{
				Name:     cookies.SessionID,
				Value:    sid.Value,
				MaxAge:   int(conf.Cookie.MaxAge),
				Path:     "/",
				HttpOnly: true,
				Secure:   conf.Cookie.Secure,
				Domain:   conf.Cookie.Domain,
				SameSite: http.SameSiteStrictMode,
			}
			http.SetCookie(w, cookie)
		}

		ctx = context.WithValue(ctx, contexts.User, user)
		ctx = context.WithValue(ctx, contexts.UserID, user.ID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AccountOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, ok := ctx.Value(contexts.User).(User)

		if !ok {
			slog.LogAttrs(ctx, slog.LevelInfo, "no session cookie found")
			http.Redirect(w, r, "/otp.html", http.StatusFound)
			return
		}

		ctx = context.WithValue(ctx, contexts.Demo, user.Demo)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Domain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if strings.HasPrefix(r.URL.Path, "/admin") {
			ctx = context.WithValue(ctx, contexts.Domain, "back")
		} else {
			ctx = context.WithValue(ctx, contexts.Domain, "front")
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, ok := ctx.Value(contexts.User).(User)

		if !ok {
			slog.LogAttrs(ctx, slog.LevelInfo, "no session cookie found")
			http.Redirect(w, r, "/sso.html", http.StatusFound)
			return
		}

		if user.Role != "admin" {
			slog.LogAttrs(ctx, slog.LevelInfo, "the user is not admin", slog.Int("id", user.ID))
			httperrors.Catch(w, ctx, "you are not authorized to process this request", 401)
			return
		}

		ctx = context.WithValue(ctx, contexts.Demo, user.Demo)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
