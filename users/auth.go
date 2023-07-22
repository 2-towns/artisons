package users

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/httputil"
	"gifthub/string/stringutil"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// ContextKey is the context key used to store the lang
const ContextKey httputil.ContextKey = "user"

// Login makes user authentication.
// The password hash will be checked with the one in the database.
// If the username does no exist or the password does not match,
// an error occurs.
// If the login is successful, a session ID is created.
// The user ID is stored with the key auth:sessionID with an expiration time,
// and the session ID in the user data.
func Login(username string, password string) (string, error) {
	ctx := context.Background()
	id, err := db.Redis.HGet(ctx, "user", username).Result()
	if err != nil {
		log.Printf("authn_login_fail: the username %s does not exist", username)
		return "", errors.New("user_login_failed")
	}

	key := "user:" + id
	hash, err := db.Redis.HGet(ctx, key, "hash").Result()
	if err != nil {
		log.Printf("authn_login_fail: error when retrieving the user id %s %v", id, err)
		return "", errors.New("user_login_failed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		log.Printf("authn_login_fail: error when comparing password hash %s", err.Error())
		return "", errors.New("user_login_failed")
	}

	var sessionID string
	var previousSessionID string

	if conf.IsMultipleSessionEnabled == true {
		// TODO: Manage multiple session ID
	} else {
		sessionID, err = stringutil.Random()
		if err != nil {
			log.Printf("ERROR: sequence_fail: error when generating a new session ID %s", err.Error())
			return "", errors.New("something_went_wrong")
		}

		previousSessionID, err = db.Redis.HGet(ctx, "user:"+id, "session_id").Result()
		if err != nil {
			log.Printf("WARN: sequence_fail: error when getting previous session ID %s", err.Error())
		}
	}

	if _, err = db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		if previousSessionID != "" {
			rdb.Del(ctx, "auth:"+previousSessionID)
		}
		rdb.Set(ctx, "auth:"+sessionID, id, conf.SessionDuration)
		rdb.HSet(ctx, "user:"+id, "session_id", sessionID)
		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return "", errors.New("something_went_wrong")
	}

	log.Printf("authn_login_success: user ID %s did a successful login", id)
	log.Printf("authn_token_created: session ID generated %s for user ID %s\n", sessionID, id)

	return sessionID, nil
}

// Logout destroys the user session.
func (u User) Logout() error {
	if u.SessionID == "" {
		return nil
	}

	ctx := context.Background()
	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Del(ctx, "auth:"+u.SessionID)
		rdb.HDel(ctx, fmt.Sprintf("user:%d", u.ID), "auth")
		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return errors.New("something_went_wrong")
	}

	return nil
}

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

	if u.SessionID != sessionID {
		log.Printf("WARN: authz_fail: the redis session ID %s does not match the cookie %s ", u.SessionID, sessionID)
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
