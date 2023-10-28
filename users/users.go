// Package users provide everything around users
package users

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/httputil"
	"gifthub/locales"
	"gifthub/string/stringutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"golang.org/x/text/language"
)

// ContextKey is the context key used to store the lang
const ContextKey httputil.ContextKey = "user"

// User is the user representation in the application
type User struct {
	// The current session ID
	SID string

	Email     string
	ID        int64
	Address   Address
	CreatedAt time.Time
	UpdatedAt time.Time
	MagicCode string

	// The key is the session id and the value
	// if the device information, the user agent.
	Devices map[string]string

	// The web push tokens
	WPTokens []string

	Lang language.Tag
}

type Session struct {
	ID     string
	Device string
	TTL    time.Duration
}

type Address struct {
	Lastname      string `validate:"required"`
	Firstname     string `validate:"required"`
	City          string `validate:"required"`
	Address       string `validate:"required"`
	Complementary string
	Zipcode       string `validate:"required"`
	Phone         string `validate:"required"`
}

func updateMagicIfExists(email, magic string) (bool, error) {
	ctx := context.Background()

	exists, err := db.Redis.Exists(ctx, "user:"+email).Result()
	if err != nil {
		log.Printf("ERROR: sequence_fail when checking existings email %s : %s", email, err.Error())
		return false, errors.New("something_went_wrong")
	}

	if exists == 0 {
		return false, nil
	}

	uid, err := db.Redis.HGet(ctx, "user:"+email, "id").Result()
	if err != nil {
		log.Printf("ERROR: sequence_fail when getting the user id %s : %s", email, err.Error())
		return false, errors.New("something_went_wrong")
	}

	if _, err := db.Redis.Set(ctx, "magic:"+magic, uid, conf.MagicCodeDuration).Result(); err != nil {
		log.Printf("ERROR: sequence_fail when storing the magic code : %s", err.Error())
		return false, errors.New("something_went_wrong")
	}

	return true, nil
}

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
	if done, err := updateMagicIfExists(email, magic); err != nil || done {
		return magic, err
	}

	id, err := db.Redis.Incr(ctx, "user_next_id").Result()
	if err != nil {
		log.Printf("ERROR: sequence_fail when creating a new id: %s", err.Error())
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
		log.Printf("ERROR: sequence_fail: error when parsing id %s %s", m["id"], err.Error())
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

	devices := make(map[string]string)
	wptokens := []string{}
	for key, value := range m {
		if strings.HasPrefix("auth:", key) {
			k := strings.Replace(key, "auth:", "", 1)
			devices[k] = value
		}

		if strings.HasPrefix("wptoken:", key) {
			wptokens = append(wptokens, value)
		}
	}

	return User{
		ID:        id,
		SID:       m["sid"],
		Email:     m["email"],
		MagicCode: m["magic"],
		Lang:      language.Make(m["lang"]),
		Address: Address{
			Lastname:      m["lastname"],
			Firstname:     m["firstname"],
			Address:       m["street"],
			Complementary: m["complementary"],
			Zipcode:       m["zipcode"],
			City:          m["city"],
			Phone:         m["phone"],
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Devices:   devices,
		WPTokens:  wptokens,
	}, nil
}

// List returns the users list in the application
func List(page int64) ([]User, error) {
	key := "users"
	ctx := context.Background()

	start, end := conf.Pagination(page)
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

// SaveAddress attachs an address to an user
func (u User) SaveAddress(a Address) error {
	v := validator.New()
	if err := v.Struct(a); err != nil {
		log.Printf("input_validation_fail: error when validation user %s", err.Error())
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("user_%s_required", low)
	}

	if u.ID == 0 {
		log.Printf("sequence_fail: the user id is empty")
		return errors.New("something_went_wrong")
	}

	ctx := context.Background()
	if _, err := db.Redis.HSet(ctx, fmt.Sprintf("user:%d", u.ID),
		"firstname", a.Firstname,
		"lastname", a.Lastname,
		"complementary", a.Complementary,
		"city", a.City,
		"phone", a.Phone,
		"zipcode", a.Zipcode,
		"street", a.Address,
	).Result(); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return errors.New("something_went_wrong")
	}

	return nil
}

// Sessions retrieve the active user sessions.
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
			continue
		}

		ttl := cmd.(*redis.DurationCmd).Val()
		if ttl.Nanoseconds() < 0 {
			log.Printf("session_expired: the session %s is expired %s", key, ttl)
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
func Login(magic, device string) (string, error) {
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
	if err != nil || uid == "" {
		log.Printf("authn_login_fail: the magic code %s is not valid", magic)
		return "", errors.New("user_magic_code_invalid")
	}

	id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when parsing id %s", uid)
		return "", errors.New("something_went_wrong")
	}

	sid, err := stringutil.Random()
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when generating a new session ID %s", err.Error())
		return "", errors.New("something_went_wrong")
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(ctx, "auth:"+sid, id, conf.SessionDuration)
		rdb.HSet(ctx, fmt.Sprintf("user:%d", id), "auth:"+sid, device)
		rdb.HSet(ctx, fmt.Sprintf("user:%d", id), "lang", locales.Default.String())
		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return "", errors.New("something_went_wrong")
	}

	log.Printf("authn_login_success: user ID %s did a successful login on device %s", uid, device)
	log.Printf("authn_token_created: session ID generated %s for user ID %s\n", sid, uid)

	return sid, nil
}

// Logout destroys the user session.
func Logout(sid string) error {
	if sid == "" {
		log.Printf("input_validation_fail: the session id is required")
		return errors.New("unauthorized")
	}

	ctx := context.Background()
	uid, err := db.Redis.Get(ctx, "auth:"+sid).Result()
	if err != nil || uid == "" {
		log.Printf("input_validation_fail: the session id does not exist")
		return errors.New("unauthorized")
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Del(ctx, "auth:"+sid)
		rdb.HDel(ctx, "user:%s"+uid, "auth:"+sid)

		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return errors.New("something_went_wrong")
	}

	log.Printf("authn_token_revoked: user session %s destroyed", sid)

	return nil
}

// Get the user information from its id
func Get(id int64) (User, error) {
	if id == 0 {
		log.Printf("input_validation_fail: the user id is required")
		return User{}, errors.New("user_not_found")
	}

	ctx := context.Background()
	data, err := db.Redis.HGetAll(ctx, fmt.Sprintf("user:%d", id)).Result()
	if err != nil {
		log.Printf("sequence_fail: error when getting data from redis %s", err.Error())
		return User{}, errors.New("something_went_wrong")
	}

	if data["id"] == "" {
		log.Printf("input_validation_fail: the user is not found")
		return User{}, errors.New("user_not_found")
	}

	return parseUser(data)
}
