package users

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/string/stringutil"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// Login makes user authentication.
// The password hash will be checked with the one in the database.
// If the username does no exist or the password does not match,
// an error occurs.
// If the login is successful, a session ID is created.
// The user ID is stored with the key auth:sessionID with an expiration time,
// and the session ID in the user data. So the user can have only one session active.
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

	sessionID, err := stringutil.Random()
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when generating a new session ID %s", err.Error())
		return "", errors.New("something_went_wrong")
	}

	if _, err = db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(ctx, "auth:"+sessionID, id, conf.SessionDuration)
		rdb.HSet(ctx, "user:"+id, "auth", sessionID)
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
