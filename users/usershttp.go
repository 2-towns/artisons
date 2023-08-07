package users

import (
	"context"
	"errors"
	"gifthub/conf"
	"gifthub/db"
	"log"
	"net/http"
)

func findBySessionID(sid string) (User, error) {
	if sid == "" {
		log.Printf("input_validation_fail: the session id is required")
		return User{}, errors.New("unauthorized")
	}

	ctx := context.Background()
	id, err := db.Redis.Get(ctx, "auth:"+sid).Result()
	if err != nil {
		log.Printf("WARN: authz_fail: error when looking for session %s %s", sid, err.Error())
		return User{}, errors.New("unauthorized")
	}

	m, err := db.Redis.HGetAll(ctx, "user:"+id).Result()
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when loading redis data for session %s %s", sid, err.Error())
		return User{}, errors.New("unauthorized")
	}

	u, err := parseUser(m)
	if err != nil {
		return User{}, errors.New("unauthorized")
	}

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

		user, err := findBySessionID(sid.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
