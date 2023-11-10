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
	"log/slog"
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
	Lang      language.Tag
	Role      string
}

type Session struct {
	ID      string
	Device  string
	WPToken string
	TTL     time.Duration
}

type Address struct {
	Lastname      string `validate:"required"`
	Firstname     string `validate:"required"`
	City          string `validate:"required"`
	Street        string `validate:"required"`
	Complementary string
	Zipcode       string `validate:"required"`
	Phone         string `validate:"required"`
}

func updateMagicIfExists(c context.Context, email, magic string) (bool, error) {
	l := slog.With(slog.String("email", email), slog.String("magic", magic))
	l.LogAttrs(c, slog.LevelInfo, "updating the magic code")

	ctx := context.Background()

	exists, err := db.Redis.Exists(ctx, "user:"+email).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot find the email", slog.String("error", err.Error()))
		return false, errors.New("something_went_wrong")
	}

	if exists == 0 {
		l.LogAttrs(c, slog.LevelInfo, "cannot find the email")
		return false, nil
	}

	uid, err := db.Redis.HGet(ctx, "user:"+email, "id").Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot find the user id", slog.String("user_id", uid), slog.String("error", err.Error()))
		return false, errors.New("something_went_wrong")
	}

	if _, err := db.Redis.Set(ctx, "magic:"+magic, uid, conf.MagicCodeDuration).Result(); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot store the magic code", slog.String("error", err.Error()))
		return false, errors.New("something_went_wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "magic code updated")

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
// The keys are:
// - user:id => the user date
// - magic:code => the magic code
// - user email => the email email link with the id
// - users => the user id list
// The email does not pass the validation, an error occurs.
// The user ID is an incremented field in Redis and returned.
func MagicCode(c context.Context, email string) (string, error) {
	l := slog.With(slog.String("email", email))
	l.LogAttrs(c, slog.LevelInfo, "generating a magic code")

	v := validator.New()
	if err := v.Var(email, "required,email"); err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the email", slog.String("error", err.Error()))
		return "", errors.New("input_email_invalid")
	}

	magic, err := stringutil.Random()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot generate the magic link", slog.String("error", err.Error()))
		return "", errors.New("something_went_wrong")
	}

	ctx := context.Background()
	if done, err := updateMagicIfExists(c, email, magic); err != nil || done {
		return magic, err
	}

	id, err := db.Redis.Incr(ctx, "user_next_id").Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot get the next id", slog.String("error", err.Error()))
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
		l.LogAttrs(c, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return "", errors.New("something_went_wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "magic code updated", slog.String("magic", magic), slog.Int64("user_id", id))

	return magic, nil
}

// Delete all the user data.
// The keys to delete:
// - user:id => the user data
// - magic:code => the magic code if it exits
// - user email => the email link with the magic code
// - auth:sid:session => the session data
// - user:id:sessions => the session ids list
func (u User) Delete(c context.Context) error {
	l := slog.With(slog.String("sid", u.SID))
	l.LogAttrs(c, slog.LevelInfo, "deleting the user")

	ctx := context.Background()

	ids, err := db.Redis.SMembers(ctx, fmt.Sprintf("user:%d:sessions", u.ID)).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot retrieve the session id list", slog.String("error", err.Error()))
		return errors.New("something_went_wrong")
	}

	key := fmt.Sprintf("user:%d", u.ID)
	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Del(ctx, key)

		if u.MagicCode != "" {
			rdb.Del(ctx, "magic:"+u.MagicCode)
		}

		rdb.HDel(ctx, "user", u.Email)
		rdb.ZRem(ctx, "users", u.ID)

		for _, sid := range ids {
			rdb.Del(ctx, "auth:"+sid+":session")
		}

		rdb.Del(ctx, fmt.Sprintf("user:%d:sessions", u.ID))

		return nil
	}); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return errors.New("something_went_wrong")
	}

	l.LogAttrs(c, slog.LevelWarn, "the user is deleted")

	return nil
}

func parseUser(c context.Context, m map[string]string) (User, error) {
	l := slog.With(slog.String("user_id", m["id"]))
	l.LogAttrs(c, slog.LevelInfo, "parsing the user data")

	id, err := strconv.ParseInt(m["id"], 10, 64)
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot parse the id", slog.String("error", err.Error()))
		return User{}, errors.New("something_went_wrong")
	}

	createdAt, err := time.Parse(time.RFC3339, m["created_at"])
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot parse the created_at", slog.String("created_at", m["created_at"]), slog.String("error", err.Error()))
		return User{}, errors.New("something_went_wrong")
	}

	updatedAt, err := time.Parse(time.RFC3339, m["updated_at"])
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot parse the updated_at", slog.String("updated_at", m["updated_at"]), slog.String("error", err.Error()))
		return User{}, errors.New("something_went_wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "the user is parse", slog.String("sid", m["sid"]))

	return User{
		ID:        id,
		SID:       m["sid"],
		Email:     m["email"],
		MagicCode: m["magic"],
		Lang:      language.Make(m["lang"]),
		Address: Address{
			Lastname:      m["lastname"],
			Firstname:     m["firstname"],
			Street:        m["street"],
			Complementary: m["complementary"],
			Zipcode:       m["zipcode"],
			City:          m["city"],
			Phone:         m["phone"],
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Role:      m["role"],
	}, nil
}

// List returns the users list in the application
func List(c context.Context, page int) ([]User, error) {
	l := slog.With(slog.Int("page", page))
	l.LogAttrs(c, slog.LevelInfo, "listing the users")

	key := "users"
	ctx := context.Background()

	start, end := conf.Pagination(page)
	users := []User{}
	ids := db.Redis.ZRange(ctx, key, int64(start), int64(end)).Val()
	pipe := db.Redis.Pipeline()

	for _, v := range ids {
		k := "user:" + v
		pipe.HGetAll(ctx, k).Val()
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot get the user list", slog.String("error", err.Error()))
		return users, errors.New("something_went_wrong")
	}

	for _, cmd := range cmds {
		m := cmd.(*redis.MapStringStringCmd).Val()
		user, err := parseUser(c, m)
		if err != nil {
			continue
		}

		users = append(users, user)
	}

	l.LogAttrs(c, slog.LevelInfo, "got user list", slog.Int("users", len(users)))

	return users, nil
}

// SaveAddress attachs an address to an user.
// The data are stored with:
// - user:id => the address
func (u User) SaveAddress(c context.Context, a Address) error {
	slog.LogAttrs(c, slog.LevelInfo, "saving the address")

	v := validator.New()
	if err := v.Struct(a); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot validate the user", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input_%s_required", low)
	}

	if u.ID == 0 {
		slog.LogAttrs(c, slog.LevelInfo, "cannot validate the user id while it is empty")
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
		"street", a.Street,
	).Result(); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot store the user", slog.String("error", err.Error()))
		return errors.New("something_went_wrong")
	}

	slog.LogAttrs(c, slog.LevelInfo, "the address is saved")

	return nil
}

// Sessions retrieve the active user sessions.
func (u User) Sessions(c context.Context) ([]Session, error) {
	slog.LogAttrs(c, slog.LevelInfo, "listing the sessions")

	ctx := context.Background()
	var sessions []Session

	ids, err := db.Redis.SMembers(ctx, fmt.Sprintf("user:%d:sessions", u.ID)).Result()
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get the session ids", slog.String("error", err.Error()))
		return sessions, errors.New("something_went_wrong")
	}

	pipe := db.Redis.Pipeline()

	for _, id := range ids {
		pipe.TTL(ctx, "auth:"+id)
	}

	ttls, err := pipe.Exec(ctx)
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get the session details", slog.String("error", err.Error()))
		return sessions, errors.New("something_went_wrong")
	}

	for _, id := range ids {
		pipe.HGetAll(ctx, "auth:"+id+":session")
	}

	scmds, err := pipe.Exec(ctx)
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get the sessions", slog.String("error", err.Error()))
		return sessions, errors.New("something_went_wrong")
	}

	for index, cmd := range ttls {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil {
			slog.LogAttrs(c, slog.LevelError, "cannot get the session ttl", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		ttl := cmd.(*redis.DurationCmd).Val()
		if ttl.Nanoseconds() < 0 {
			slog.LogAttrs(c, slog.LevelInfo, "the session is expired", slog.String("key", key), slog.Duration("ttl", ttl))
			continue
		}

		scmd := scmds[index]
		if scmd.Err() != nil {
			slog.LogAttrs(c, slog.LevelError, "cannot get the session data", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		data := scmd.(*redis.MapStringStringCmd).Val()

		id := strings.Replace(key, "auth:", "", 1)
		sessions = append(sessions, Session{
			ID:      id,
			Device:  data["device"],
			WPToken: data["wptoken"],
			TTL:     ttl,
		})
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		slog.LogAttrs(c, slog.LevelWarn, "cannot delete the expired session", slog.String("error", err.Error()))
	}

	slog.LogAttrs(c, slog.LevelInfo, "got the sessions", slog.Int("sessions", len(sessions)))

	return sessions, nil
}

// Login authenicate user with a magic code.
// If the magic is empty, an error occurs.
// If the login is successful, a session ID is created.
// The data are stored with:
// - auth:sid => the user id with an expiration key
// - auth:sid:session device => the device related to the session
// - user:id:sessions => the session id set (list)
// An extra data is stored in order to retreive all the sessions for an user.
func Login(c context.Context, magic, device string) (string, error) {
	l := slog.With(slog.String("magic", magic))
	l.LogAttrs(c, slog.LevelInfo, "trying to login", slog.String("device", device))

	if magic == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the magic code")
		return "", errors.New("input_magic_code_required")
	}

	if device == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the device")
		return "", errors.New("input_device_required")
	}

	ctx := context.Background()
	uid, err := db.Redis.Get(ctx, "magic:"+magic).Result()
	if err != nil || uid == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot find the magic code")
		return "", errors.New("user_magic_code_invalid")
	}

	id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot parse the uid", slog.String("user_id", uid), slog.String("error", err.Error()))
		return "", errors.New("something_went_wrong")
	}

	sid, err := stringutil.Random()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot parse the session id", slog.String("sid", sid), slog.String("user_id", uid), slog.String("error", err.Error()))
		return "", errors.New("something_went_wrong")
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(ctx, "auth:"+sid, id, conf.SessionDuration)
		rdb.HSet(ctx, fmt.Sprintf("auth:%s:session", sid), "device", device)
		rdb.HSet(ctx, fmt.Sprintf("user:%d", id), "lang", locales.Default.String())
		rdb.SAdd(ctx, fmt.Sprintf("user:%d:sessions", id), sid)

		return nil
	}); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot store the data", slog.String("sid", sid), slog.String("user_id", uid), slog.String("error", err.Error()))
		return "", errors.New("something_went_wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "the login is successful", slog.String("device", device), slog.String("sid", sid), slog.String("user_id", uid))

	return sid, nil
}

// Logout destroys the user session.
func Logout(c context.Context, sid string) error {
	l := slog.With(slog.String("sid", sid))
	l.LogAttrs(c, slog.LevelInfo, "trying to logout")

	if sid == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the session id")
		return errors.New("unauthorized")
	}

	ctx := context.Background()
	uid, err := db.Redis.Get(ctx, "auth:"+sid).Result()
	if err != nil || uid == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot find the session id")
		return errors.New("unauthorized")
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Del(ctx, "auth:"+sid)
		rdb.HDel(ctx, "user:%s"+uid, "auth:"+sid)

		return nil
	}); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return errors.New("something_went_wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "the logout is successful")

	return nil
}

// Get the user information from its id
func Get(c context.Context, id int64) (User, error) {
	l := slog.With(slog.Int64("user_id", id))
	l.LogAttrs(c, slog.LevelInfo, "trying to get the user")

	if id == 0 {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the user id")
		return User{}, errors.New("user_not_found")
	}

	ctx := context.Background()
	data, err := db.Redis.HGetAll(ctx, fmt.Sprintf("user:%d", id)).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot get the user from redis", slog.String("error", err.Error()))
		return User{}, errors.New("something_went_wrong")
	}

	if data["id"] == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot find the user")
		return User{}, errors.New("user_not_found")
	}

	u, err := parseUser(c, data)

	l.LogAttrs(c, slog.LevelInfo, "the user is found")

	return u, err
}
