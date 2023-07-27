package users

import (
	"context"
	"gifthub/conf"
	"gifthub/db"
	"log"
	"net/http"
)

func findBySessionID(sessionID string) (User, error) {
	ctx := context.Background()
	id, err := db.Redis.Get(ctx, "auth:"+sessionID).Result()
	if err != nil {
		log.Printf("WARN: authz_fail: error when looking for session %s %s", sessionID, err.Error())
		return User{}, err
	}

	m, err := db.Redis.HGetAll(ctx, "user:"+id).Result()
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when loading redis data for session %s %s", sessionID, err.Error())
		return User{}, err
	}

	u, err := parseUser(m)
	if err != nil {
		log.Printf("ERROR: sequence_fail: error parsing data for user %s %s ", id, err.Error())
		return User{}, err
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
