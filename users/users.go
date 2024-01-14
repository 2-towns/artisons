// Package users provide everything around users
package users

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/contexts"
	"gifthub/notifications/mails"
	"gifthub/string/stringutil"
	"gifthub/tracking"
	"gifthub/validators"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// User is the user representation in the application
type User struct {
	// The current session ID
	SID string

	Email     string
	ID        int64
	Address   Address
	CreatedAt int64
	UpdatedAt int64
	Otp       string
	Lang      language.Tag

	// admin or user
	Role string

	Demo bool
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

// OTP generates a login code.
// An error is raised if an otp was already generated in the ttl period.
// A glue is created to restrict the otp attempt to the source device only.
// All the operations are executed in a single transaction.
// The keys are:
// - {{email}}:otp => the otp
// - {{email}}:attempts => set to 0 the OTP attempts
// - otp:{{glue}} => the email
// The email does not pass the validation, an error occurs.
func Otp(c context.Context, email string) (string, error) {
	l := slog.With(slog.String("email", email))
	l.LogAttrs(c, slog.LevelInfo, "generating a otp code")

	if err := validators.V.Var(email, "required,email"); err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the email", slog.String("error", err.Error()))
		return "", errors.New("input:email")
	}

	ctx := context.Background()

	ttl, err := db.Redis.TTL(ctx, email+":otp").Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot get the ttl", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	if conf.OtpDuration-ttl < conf.OtpInterval {
		l.LogAttrs(c, slog.LevelInfo, "the ttl exceed the otp interval", slog.Duration("ttl", ttl))
		return "", errors.New("you need to wait before asking another otp")
	}

	otp := rand.Intn(999999-100000) + 100000

	glue, err := stringutil.Random()
	if err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot generate the glue", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	if _, err = db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(ctx, fmt.Sprintf("%s:otp", email), otp, conf.OtpDuration)
		rdb.Set(ctx, fmt.Sprintf("%s:otp:attempts", email), 0, conf.OtpDuration)
		rdb.Set(ctx, fmt.Sprintf("otp:%s", glue), email, conf.OtpDuration)
		return nil
	}); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	go func() {
		lang := c.Value(contexts.Locale).(language.Tag)
		p := message.NewPrinter(lang)
		mails.Send(c, email, p.Sprintf("email_otp_subject"), p.Sprintf("email_otp_login", fmt.Sprintf("%d", otp)))
	}()

	l.LogAttrs(c, slog.LevelInfo, "otp code updated", slog.Int("otp", otp), slog.String("glue", glue))

	return glue, nil
}

// Delete all the user data.
// The keys to delete:
// - user:id => the user data
// - otp:code => the otp code if it exits
// - user email => the email link with the otp code
// - auth:sid:session => the session data
// - user:id:sessions => the session ids list
func (u User) Delete(c context.Context) error {
	l := slog.With(slog.String("sid", u.SID))
	l.LogAttrs(c, slog.LevelInfo, "deleting the user")

	ctx := context.Background()

	ids, err := db.Redis.SMembers(ctx, fmt.Sprintf("user:%d:sessions", u.ID)).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot retrieve the session id list", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	key := fmt.Sprintf("user:%d", u.ID)
	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Del(ctx, key)

		if u.Otp != "" {
			rdb.Del(ctx, "otp:"+u.Otp)
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
		return errors.New("something went wrong")
	}

	l.LogAttrs(c, slog.LevelWarn, "the user is deleted")

	return nil
}

func parseUser(c context.Context, m map[string]string) (User, error) {
	l := slog.With(slog.String("user_id", m["id"]))
	l.LogAttrs(c, slog.LevelInfo, "parsing the user data")

	if m["id"] == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot continue with empty id")
		return User{}, errors.New("something went wrong")
	}

	id, err := strconv.ParseInt(m["id"], 10, 64)
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot parse the id", slog.String("error", err.Error()))
		return User{}, errors.New("something went wrong")
	}

	createdAt, err := strconv.ParseInt(m["created_at"], 10, 64)
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot parse the created_at", slog.String("created_at", m["created_at"]), slog.String("error", err.Error()))
		return User{}, errors.New("something went wrong")
	}

	updatedAt, err := strconv.ParseInt(m["updated_at"], 10, 64)
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot parse the updated_at", slog.String("updated_at", m["updated_at"]), slog.String("error", err.Error()))
		return User{}, errors.New("something went wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "the user is parsed", slog.String("sid", m["sid"]))

	return User{
		ID:    id,
		SID:   m["sid"],
		Email: m["email"],
		Otp:   m["otp"],
		Lang:  language.Make(m["lang"]),
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
		Demo:      m["demo"] == "1",
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
		return users, errors.New("something went wrong")
	}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil {
			slog.LogAttrs(c, slog.LevelError, "cannot get the user", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

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

	if err := validators.V.Struct(a); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot validate the user", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input:%s", low)
	}

	if u.ID == 0 {
		slog.LogAttrs(c, slog.LevelInfo, "cannot validate the user id while it is empty")
		return errors.New("something went wrong")
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
		return errors.New("something went wrong")
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
		return sessions, errors.New("something went wrong")
	}

	pipe := db.Redis.Pipeline()

	for _, id := range ids {
		pipe.TTL(ctx, "auth:"+id)
	}

	ttls, err := pipe.Exec(ctx)
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get the session details", slog.String("error", err.Error()))
		return sessions, errors.New("something went wrong")
	}

	for _, id := range ids {
		pipe.HGetAll(ctx, "auth:"+id+":session")
	}

	scmds, err := pipe.Exec(ctx)
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get the sessions", slog.String("error", err.Error()))
		return sessions, errors.New("something went wrong")
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

// Login authenicate user with a otp code.
// If the otp or the device is empty, an error occurs.
// If the glue does not match any email, an error occurs.
// If the otp does not exist for the email, an error occurs.
// If the otp does not match the otp generated, an errors occurs.
// If the maximum attempts is reached, an error occured and the otp related
// to the data are deleted.
// If the id does not exist for the email, the id sequence is incremented
// and the value is linked to the email.
// If the login is successful, a session ID is created.
// The data are stored with:
// - auth:sid => the user id with an expiration key
// - auth:sid:session device => the device related to the session
// - user:id:sessions => the session id set (list)
// - user:id => the user data
// An extra data is stored in order to retreive all the sessions for an user.
func Login(c context.Context, otp, glue, device string) (string, error) {
	l := slog.With(slog.String("otp", otp), slog.String("glue", glue))
	l.LogAttrs(c, slog.LevelInfo, "trying to login", slog.String("device", device))

	if otp == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the otp code")
		return "", errors.New("input:otp")
	}

	if device == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the device")
		return "", errors.New("your are not authorized to access to this page")
	}

	if glue == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the glue")
		return "", errors.New("your are not authorized to process this request")
	}

	ctx := context.Background()

	email, err := db.Redis.Get(ctx, "otp:"+glue).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot find the glue", slog.String("error", err.Error()))
		return "", errors.New("your are not authorized to process this request")
	}

	val, err := db.Redis.Get(ctx, email+":otp").Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot get the existing otp", slog.String("error", err.Error()))
		return "", errors.New("your are not authorized to process this request")
	}

	if val != otp && !(conf.OtpDemo && otp == "111111") {
		l.LogAttrs(c, slog.LevelInfo, "the otp do not match", slog.String("val", val), slog.String("otp", otp))

		cnt, err := db.Redis.Incr(ctx, email+":otp:attempts").Result()
		if err != nil {
			l.LogAttrs(c, slog.LevelError, "cannot increment the otp attempt", slog.String("error", err.Error()))
			return "", errors.New("something went wrong")
		}

		if cnt >= conf.OtpAttempts {
			if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
				rdb.Del(ctx, fmt.Sprintf("%s:otp", email))
				rdb.Del(ctx, fmt.Sprintf("%s:otp:attempts", email))
				rdb.Del(ctx, fmt.Sprintf("otp:%s", glue))

				return nil
			}); err != nil {
				l.LogAttrs(c, slog.LevelError, "cannot destory the otp", slog.String("error", err.Error()))
				return "", errors.New("something went wrong")
			}

			l.LogAttrs(c, slog.LevelInfo, "max attempts reached", slog.Int64("attempts", cnt))
			return "", errors.New("you reached the max tentatives")
		}

		return "", errors.New("the OTP does not match")
	}

	eid, err := db.Redis.Get(ctx, email+":id").Result()
	if err != nil && err.Error() != "redis: nil" {
		l.LogAttrs(c, slog.LevelError, "cannot verify id existence", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	var uid int64

	if eid == "" {
		uid, err = db.Redis.Incr(ctx, "user_next_id").Result()

		if err != nil {
			l.LogAttrs(c, slog.LevelError, "cannot get the next id", slog.String("error", err.Error()))
			return "", errors.New("something went wrong")
		}
	} else {
		uid, err = strconv.ParseInt(eid, 10, 64)
		if err != nil {
			l.LogAttrs(c, slog.LevelError, "cannot parse the uid", slog.String("user_id", eid), slog.String("error", err.Error()))
			return "", errors.New("something went wrong")
		}
	}

	sid, err := stringutil.Random()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot generated the session id", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	now := time.Now()

	role := "user"
	if admin := IsAdmin(ctx, email); admin {
		role = "admin"
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(ctx, "auth:"+sid, uid, conf.SessionDuration)
		rdb.HSet(ctx, fmt.Sprintf("auth:%s:session", sid), "device", device)
		// @todo get the lang from the browser and match with the ones on the server
		rdb.HSet(ctx, fmt.Sprintf("user:%d", uid), "lang", conf.DefaultLocale.String())
		rdb.SAdd(ctx, fmt.Sprintf("user:%d:sessions", uid), sid)
		rdb.HSet(ctx, fmt.Sprintf("user:%d", uid),
			"id", uid,
			"email", email,
			"updated_at", now.Unix(),
			"created_at", now.Unix(),
			"role", role,
		)
		rdb.Del(ctx, fmt.Sprintf("%s:otp", email))
		rdb.Del(ctx, fmt.Sprintf("%s:otp:attempts", email))
		rdb.Del(ctx, fmt.Sprintf("otp:%s", glue))

		return nil
	}); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot store the data", slog.String("sid", sid), slog.Int64("user_id", uid), slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	data := map[string]string{
		"sid":   sid,
		"otp":   otp,
		"email": email,
	}

	if role != "admin" {
		go tracking.Log(c, "login", data)
	}

	l.LogAttrs(c, slog.LevelInfo, "the login is successful", slog.String("device", device), slog.String("sid", sid), slog.Int64("user_id", uid))

	return sid, nil
}

// Logout destroys the user session.
func Logout(c context.Context, sid string) error {
	l := slog.With(slog.String("sid", sid))
	l.LogAttrs(c, slog.LevelInfo, "trying to logout")

	if sid == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the session id")
		return errors.New("your are not authorized to process this request")
	}

	ctx := context.Background()
	uid, err := db.Redis.Get(ctx, "auth:"+sid).Result()
	if err != nil || uid == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot find the session id")
		return errors.New("your are not authorized to process this request")
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Del(ctx, "auth:"+sid)
		rdb.HDel(ctx, "user:%s"+uid, "auth:"+sid)

		return nil
	}); err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	data := map[string]string{
		"sid": sid,
	}

	go tracking.Log(c, "logout", data)

	l.LogAttrs(c, slog.LevelInfo, "the logout is successful")

	return nil
}

// Get the user information from its id
func Get(c context.Context, id int64) (User, error) {
	l := slog.With(slog.Int64("user_id", id))
	l.LogAttrs(c, slog.LevelInfo, "trying to get the user")

	if id == 0 {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate the user id")
		return User{}, errors.New("the user is not found")
	}

	ctx := context.Background()
	data, err := db.Redis.HGetAll(ctx, fmt.Sprintf("user:%d", id)).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot get the user from redis", slog.String("error", err.Error()))
		return User{}, errors.New("something went wrong")
	}

	if data["id"] == "" {
		l.LogAttrs(c, slog.LevelInfo, "cannot find the user")
		return User{}, errors.New("the user is not found")
	}

	u, err := parseUser(c, data)

	l.LogAttrs(c, slog.LevelInfo, "the user is found")

	return u, err
}

func IsAdmin(c context.Context, email string) bool {
	l := slog.With(slog.String("email", email))
	l.LogAttrs(c, slog.LevelInfo, "trying to known if the user is admin")

	ctx := context.Background()
	is, err := db.Redis.SIsMember(ctx, "admins", email).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot retrieve admins from redis", slog.String("error", err.Error()))
		return false
	}

	l.LogAttrs(c, slog.LevelInfo, "the user is admin", slog.Bool("yes", is))

	return is
}

func (u User) ToggleDemo(c context.Context) (bool, error) {
	l := slog.With(slog.Int64("uid", u.ID), slog.Bool("demo", u.Demo))
	l.LogAttrs(c, slog.LevelInfo, "toggle demo mode")

	v := "1"
	if u.Demo {
		v = "0"
	}

	_, err := db.Redis.HSet(context.Background(), fmt.Sprintf("user:%d", u.ID), "demo", v).Result()
	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot toggle demo mode", slog.String("error", err.Error()))
		return u.Demo, err
	}

	l.LogAttrs(c, slog.LevelInfo, "demo modo toggled", slog.Bool("demo", !u.Demo))

	return !u.Demo, nil
}
