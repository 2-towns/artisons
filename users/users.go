// Package users provide everything around users
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
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

// ContextKey is the context key used to store the lang
const ContextKey httputil.ContextKey = "user"

// User is the user representation in the application
type User struct {
	Email     string
	ID        int64
	Address   Address
	CreatedAt time.Time
	UpdatedAt time.Time
	MagicCode string
	Devices   map[string]string
}

type Session struct {
	ID     string
	Device string
	TTL    time.Duration
}

type Address struct {
	Lastname      string
	Firstname     string
	City          string
	Complementary string
	Zipcode       string
	Phone         string
}

/*	sessionID, err := stringutil.Random()
	if err != nil {
		log.Printf("ERROR: sequence_fail when generating session ID : %s", err.Error())
		return 0, errors.New("something_went_wrong")
	}*/
/*
		content := p.Sprintf("user_magik_link_email", magicLink)
	if err := mails.Send(email, content); err != nil {
		log.Printf("ERROR: sequence_fail: error wehn sending email %s", err.Error())
		return "", err
	}*/

// MagicLink generates a login code.
// If the email does not exist, an user is created with an incremented
// ID.
//
// The link between the magic and the user id is stored in Redis.
// The link between the email and the user id is also stored in Redis.
// The user is also added in a sorted set with the timestamp as score.
//
// All the operations are executed in a single transaction.
//
// The email does not pass the validation, an error occurs.
// The user ID is an incremented field in Redis and returned.
func MagicCode(email string) (string, error) {
	v := validator.New()
	if err := v.Var(email, "required,email"); err != nil {
		log.Printf("input_validation_fail: error when validation user %s", err.Error())
		return "", errors.New("user_email_invalid")
	}

	magic, err := stringutil.Random()
	if err != nil {
		log.Printf("ERROR: sequence_fail when generating magic link: %s", err.Error())
		return "", errors.New("something_went_wrong")
	}

	log.Printf("WARN: sensitive_create: a new magic code is created with email %s\n", email)

	ctx := context.Background()
	uid, err := db.Redis.HGet(ctx, "user:"+email, "id").Result()
	if uid != "" && err != nil {
		if _, err := db.Redis.Set(ctx, "magic:"+magic, uid, conf.MagicCodeDuration).Result(); err != nil {
			log.Printf("ERROR: sequence_fail when storing the magic code : %s", err.Error())
			return "", errors.New("something_went_wrong")
		}

		return magic, nil
	}

	id, err := db.Redis.Incr(ctx, "user_next_id").Result()
	if err != nil {
		log.Printf("ERROR: sequence_fail: %s", err.Error())
		return "", errors.New("something_went_wrong")
	}

	now := time.Now()
	key := fmt.Sprintf("user:%d", id)

	if _, err = db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, key,
			"id", id,
			"email", email,
			"updated_at", now.Format(time.RFC3339),
			"created_at", now.Format(time.RFC3339),
		)
		rdb.Set(ctx, "magic:"+magic, id, conf.MagicCodeDuration)
		rdb.HSet(ctx, "user", email, id)
		rdb.ZAdd(ctx, "users", redis.Z{
			Score:  float64(now.Unix()),
			Member: id,
		})
		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return "", errors.New("something_went_wrong")
	}

	log.Printf("WARN: sensitive_create: a new user is created with id %d\n", id)

	return magic, nil
}

// Delete all the user data
func (u User) Delete() error {
	ctx := context.Background()

	key := fmt.Sprintf("user:%d", u.ID)
	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Del(ctx, key)
		if u.MagicCode != "" {
			rdb.Del(ctx, "magic:"+u.MagicCode)
		}
		rdb.HDel(ctx, "user", u.Email)
		rdb.ZRem(ctx, "users", u.ID)
		for k := range u.Devices {
			rdb.Del(ctx, k)
		}

		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return errors.New("something_went_wrong")
	}

	log.Printf("WARN: sensitive_delete: the user %d is deleted\n", u.ID)

	return nil
}

func parseUser(m map[string]string) (User, error) {
	id, err := strconv.ParseInt(m["id"], 10, 32)
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when parsing id %s", m["id"])
		return User{}, errors.New("something_went_wrong")
	}

	createdAt, err := time.Parse(time.RFC3339, m["created_at"])
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when parsing created_at %s", m["created_at"])
		return User{}, errors.New("something_went_wrong")
	}

	updatedAt, err := time.Parse(time.RFC3339, m["updated_at"])
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when parsing created_at %s", m["updated_at"])
		return User{}, errors.New("something_went_wrong")
	}

	return User{
		ID:        id,
		Email:     m["email"],
		MagicCode: m["magic"],
		/*	Lastname:      m["lastname"],
			Firstname:     m["firstname"],
			Address:       m["address"],
			Complementary: m["complementary"],
			Zipcode:       m["zipcode"],
			City:          m["city"],
			Phone:         m["phone"],*/
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// List returns the users list in the application
func List(page int64) ([]User, error) {
	key := "users"
	ctx := context.Background()

	var start int64
	var end int64

	if page == -1 {
		start = 0
		end = -1
	} else {
		start = page * conf.ItemsPerPage
		end = page*conf.ItemsPerPage + conf.ItemsPerPage
	}

	users := []User{}
	ids := db.Redis.ZRange(ctx, key, start, end).Val()
	pipe := db.Redis.Pipeline()

	for _, v := range ids {
		k := "user:" + v
		pipe.HGetAll(ctx, k).Val()
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("ERROR: sequence_fail: go error from redis %s", err.Error())
		return users, errors.New("something_went_wrong")
	}

	for _, cmd := range cmds {
		m := cmd.(*redis.MapStringStringCmd).Val()
		user, err := parseUser(m)
		if err != nil {
			continue
		}

		users = append(users, user)
	}

	return users, nil
}

// Sessions retrieve the active user sessions.
// If a session is expired, it will be removed from the user session ids.
func (u User) Sessions() ([]Session, error) {
	ctx := context.Background()
	pipe := db.Redis.Pipeline()
	var keys []string
	var devices []string
	for key, device := range u.Devices {
		pipe.TTL(ctx, key)
		keys = append(keys, key)
		devices = append(devices, device)
	}

	var sessions []Session

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when getting sessions details %s", err.Error())
		return sessions, errors.New("something_went_wrong")
	}

	for index, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])
		if cmd.Err() != nil {
			log.Printf("ERROR: sequence_fail: error when retrieving session TTL %s %s", key, err.Error())
			pipe.HDel(ctx, fmt.Sprintf("user:%d", u.ID), key)
			continue
		}

		ttl := cmd.(*redis.DurationCmd).Val()
		if ttl.Nanoseconds() < 0 {
			log.Printf("session_expired: the session %s is expired %s", key, ttl)
			pipe.HDel(ctx, fmt.Sprintf("user:%d", u.ID), key)
			continue
		}

		id := strings.Replace(key, "auth:", "", 0)
		device := devices[index]
		sessions = append(sessions, Session{
			ID:     id,
			Device: device,
			TTL:    ttl,
		})
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("ERROR: sequence_fail: error deleting expired session %s", err.Error())
	}

	return sessions, nil
}

// Login authenicate user with a magic code.
// If the magic is empty, an error occurs.
// If the login is successful, a session ID is created.
// The user ID is stored with the key auth:sessionID with an expiration time.
// An extra data is stored in order to retreive all the sessions for an user.
func Login(magic string, device string) (string, error) {
	if magic == "" {
		log.Printf("input_validation_fail: the magic code is required")
		return "", errors.New("user_magic_code_required")
	}

	if device == "" {
		log.Printf("input_validation_fail: the device is required")
		return "", errors.New("user_device_required")
	}

	ctx := context.Background()
	uid, err := db.Redis.Get(ctx, "magic:"+magic).Result()
	if err != nil {
		log.Printf("authn_login_fail: the magic code %s is not valid", magic)
		return "", errors.New("user_magic_code_invalid")
	}

	id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when parsing id %s", uid)
		return "", errors.New("something_went_wrong")
	}

	sessionID, err := stringutil.Random()
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when generating a new session ID %s", err.Error())
		return "", errors.New("something_went_wrong")
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(ctx, "auth:"+sessionID, id, conf.SessionDuration)
		rdb.HSet(ctx, fmt.Sprintf("user:%d", id), "auth:"+sessionID, device)
		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return "", errors.New("something_went_wrong")
	}

	log.Printf("authn_login_success: user ID %s did a successful login on device %s", uid, device)
	log.Printf("authn_token_created: session ID generated %s for user ID %s\n", sessionID, uid)

	return sessionID, nil
}

// Logout destroys the user session.
func Logout(id int64, sid string) error {
	if sid == "" || id == 0 {
		log.Printf("input_validation_fail: the id and session id are required")
		return errors.New("user_logout_invalid")
	}

	ctx := context.Background()
	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Del(ctx, "auth:"+sid)
		rdb.Del(ctx, "session:"+sid)
		rdb.SRem(ctx, fmt.Sprintf("sessions:%d", id), sid)

		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return errors.New("something_went_wrong")
	}

	log.Printf("authn_token_revoked: user session %s destroyed", sid)

	return nil
}

/*pipe.HSet(ctx, key, "lastname", u.Address.Lastname)
pipe.HSet(ctx, key, "firstname", u.Firstname)
pipe.HSet(ctx, key, "address", u.Address)
pipe.HSet(ctx, key, "city", u.City)
pipe.HSet(ctx, key, "complemnetary", u.Complementary)
pipe.HSet(ctx, key, "zipcode", u.Zipcode)
pipe.HSet(ctx, key, "phone", u.Phone)*/
