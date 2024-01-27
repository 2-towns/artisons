package users

import (
	"context"
	"errors"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/contexts"
	"gifthub/http/cookies"
	"gifthub/http/httperrors"
	"gifthub/string/stringutil"
	"log/slog"
	"net/http"
)

func findBySessionID(ctx context.Context, sid string) (User, error) {
	l := slog.With(slog.String("sid", sid))
	l.LogAttrs(ctx, slog.LevelInfo, "finding the user")

	if sid == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate the session id")
		return User{}, errors.New("your are not authorized to process this request")
	}

	id, err := db.Redis.Get(ctx, "auth:"+sid).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the auth id from redis", slog.String("error", err.Error()))
		return User{}, errors.New("your are not authorized to process this request")
	}

	m, err := db.Redis.HGetAll(ctx, "user:"+id).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the session from redis", slog.String("error", err.Error()))
		return User{}, errors.New("your are not authorized to process this request")
	}

	m["sid"] = sid
	u, err := parseUser(ctx, m)
	if err != nil {
		return User{}, errors.New("your are not authorized to process this request")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "user found", slog.Int("user_id", u.ID))

	return u, err
}

// Middleware detects the session ID in the cookies.
// If the session ID exists, it will load the current
// user into the context.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cid string
		c, err := r.Cookie(cookies.CartID)
		if err != nil || c.Value == "" {
			cid, err = stringutil.Random()
			if err != nil {
				slog.LogAttrs(r.Context(), slog.LevelError, "cannot generate cart id", slog.String("error", err.Error()))
				httperrors.Page(w, r.Context(), err.Error(), 500)
				return
			}
		} else {
			cid = c.Value
		}

		cookie := &http.Cookie{
			Name:     cookies.CartID,
			Value:    cid,
			MaxAge:   int(conf.Cookie.MaxAge),
			Path:     "/",
			HttpOnly: true,
			Secure:   conf.Cookie.Secure,
			Domain:   conf.Cookie.Domain,
		}

		http.SetCookie(w, cookie)
		ctx := context.WithValue(r.Context(), contexts.Cart, cid)

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
			}
			http.SetCookie(w, cookie)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		} else {
			cookie := &http.Cookie{
				Name:     cookies.SessionID,
				Value:    sid.Value,
				MaxAge:   int(conf.Cookie.MaxAge),
				Path:     "/",
				HttpOnly: true,
				Secure:   conf.Cookie.Secure,
				Domain:   conf.Cookie.Domain,
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
			httperrors.Catch(w, ctx, "your are not authorized to process this request", 401)
			return
		}

		ctx = context.WithValue(ctx, contexts.Demo, user.Demo)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
