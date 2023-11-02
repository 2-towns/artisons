package users

import (
	"context"
	"errors"
	"gifthub/conf"
	"gifthub/db"
	"net/http"

	"golang.org/x/exp/slog"
)

func findBySessionID(c context.Context, sid string) (User, error) {
	l := slog.With(slog.String("sid", sid))
	l.LogAttrs(c, slog.LevelInfo, "finding the user")

	if sid == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the session id")
		return User{}, errors.New("unauthorized")
	}

	ctx := context.Background()
	id, err := db.Redis.Get(ctx, "auth:"+sid).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot get the auth if from redis", slog.String("error", err.Error()))
		return User{}, errors.New("unauthorized")
	}

	m, err := db.Redis.HGetAll(ctx, "user:"+id).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot get the session from redis", slog.String("error", err.Error()))
		return User{}, errors.New("unauthorized")
	}

	m["sid"] = sid
	u, err := parseUser(c, m)
	if err != nil {
		return User{}, errors.New("unauthorized")
	}

	l.LogAttrs(c, slog.LevelInfo, "user found", slog.Int64("user_id", u.ID))

	return u, err
}

// Middleware detects the session ID in the cookies.
// If the session ID exists, it will load the current
// user into the context.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sid, err := r.Cookie(conf.SessionIDCookie)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := findBySessionID(r.Context(), sid.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// TODO display session id in lgos
		ctx := context.WithValue(r.Context(), ContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
