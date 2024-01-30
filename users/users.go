// Package users provide everything around users
package users

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/notifications/mails"
	"artisons/string/stringutil"
	"artisons/tracking"
	"artisons/validators"
	"context"
	"errors"
	"fmt"
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
	ID        int
	Address   Address
	CreatedAt time.Time
	UpdatedAt time.Time
	Otp       string
	Lang      language.Tag

	// admin or user
	Role string

	Demo bool
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

type Query struct {
	Email string
}

type SearchResults struct {
	Total int
	Users []User
}

// OTP generates a login code.
// An error is raised if an otp was already generated in the ttl period.
// All the operations are executed in a single transaction.
// The keys are:
// - {{email}}:attempts => set to 0 the OTP attempts
// The email does not pass the validation, an error occurs.
func Otp(ctx context.Context, email string) error {
	l := slog.With(slog.String("email", email))
	l.LogAttrs(ctx, slog.LevelInfo, "generating a otp code")

	if err := validators.V.Var(email, "required,email"); err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate the email", slog.String("error", err.Error()))
		return errors.New("input:email")
	}

	ttl, err := db.Redis.TTL(ctx, "otp:"+email).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the ttl", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	if conf.OtpDuration-ttl < conf.OtpInterval {
		l.LogAttrs(ctx, slog.LevelInfo, "the ttl exceed the otp interval", slog.Duration("ttl", ttl))
		return errors.New("you need to wait before asking another otp")
	}

	otp := rand.Intn(999999-100000) + 100000

	if _, err = db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, fmt.Sprintf("otp:%s", email), "otp", otp, "attempts", 0)
		rdb.Expire(ctx, fmt.Sprintf("otp:%s", email), conf.OtpDuration)
		return nil
	}); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	go func() {
		lang := ctx.Value(contexts.Locale).(language.Tag)
		p := message.NewPrinter(lang)
		mails.Send(ctx, email, p.Sprintf("email_otp_subject"), p.Sprintf("email_otp_login", fmt.Sprintf("%d", otp)))
	}()

	l.LogAttrs(ctx, slog.LevelInfo, "otp code updated", slog.Int("otp", otp))

	return nil
}

// Delete all the user data.
// The keys to delete:
// - user:id => the user data
// - user email => the email link with the otp code
func (u User) Delete(ctx context.Context) error {
	l := slog.With(slog.String("sid", u.SID))
	l.LogAttrs(ctx, slog.LevelInfo, "deleting the user")

	sessions, err := u.Sessions(ctx)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot retrieve the session id list", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	key := fmt.Sprintf("user:%d", u.ID)
	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Del(ctx, key)

		for _, s := range sessions {
			rdb.Del(ctx, "session:"+s.ID)
		}

		return nil
	}); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(ctx, slog.LevelWarn, "the user is deleted")

	return nil
}

func parse(ctx context.Context, m map[string]string) (User, error) {
	l := slog.With(slog.String("user_id", m["id"]))

	if m["id"] == "" {
		l.LogAttrs(ctx, slog.LevelError, "cannot continue with empty id")
		return User{}, errors.New("something went wrong")
	}

	id, err := strconv.ParseInt(m["id"], 10, 64)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.String("error", err.Error()))
		return User{}, errors.New("something went wrong")
	}

	createdAt, err := strconv.ParseInt(m["created_at"], 10, 64)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot parse the created_at", slog.String("created_at", m["created_at"]), slog.String("error", err.Error()))
		return User{}, errors.New("something went wrong")
	}

	updatedAt, err := strconv.ParseInt(m["updated_at"], 10, 64)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot parse the updated_at", slog.String("updated_at", m["updated_at"]), slog.String("error", err.Error()))
		return User{}, errors.New("something went wrong")
	}

	return User{
		ID:    int(id),
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
		CreatedAt: time.Unix(createdAt, 0),
		UpdatedAt: time.Unix(updatedAt, 0),
		Role:      m["role"],
		Demo:      m["demo"] == "1",
	}, nil
}

// SaveAddress attachs an address to an user.
// The data are stored with:
// - user:id => the address
func (u User) SaveAddress(ctx context.Context, a Address) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "saving the address")

	if err := validators.V.Struct(a); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot validate the user", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input:%s", low)
	}

	if u.ID == 0 {
		slog.LogAttrs(ctx, slog.LevelError, "cannot validate the user id while it is empty")
		return errors.New("something went wrong")
	}

	if _, err := db.Redis.HSet(ctx, fmt.Sprintf("user:%d", u.ID),
		"firstname", a.Firstname,
		"lastname", a.Lastname,
		"complementary", a.Complementary,
		"city", a.City,
		"phone", a.Phone,
		"zipcode", a.Zipcode,
		"street", a.Street,
	).Result(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot store the user", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "the address is saved")

	return nil
}

// Login authenicate user with a otp code.
// If the otp or the device is empty, an error occurs.
// If the otp does not exist for the email, an error occurs.
// If the otp does not match the otp generated, an errors occurs.
// If the maximum attempts is reached, an error occured and the otp related
// to the data are deleted.
// If the id does not exist for the email, the id sequence is incremented
// and the value is linked to the email.
// If the login is successful, a session ID is created.
// The data are stored with:
// - user:id => the user data
// An extra data is stored in order to retreive all the sessions for an user.
func Login(ctx context.Context, email, otp, device string) (string, error) {
	l := slog.With(slog.String("otp", otp))
	l.LogAttrs(ctx, slog.LevelInfo, "trying to login", slog.String("device", device))

	if err := validators.V.Var(email, "required,email"); err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate the email", slog.String("error", err.Error()))
		return "", errors.New("input:email")
	}

	if otp == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate the otp code")
		return "", errors.New("input:otp")
	}

	if device == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate the device")
		return "", errors.New("your are not authorized to access to this page")
	}

	val, err := db.Redis.HGet(ctx, "otp:"+email, "otp").Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the existing otp", slog.String("error", err.Error()))
		return "", errors.New("your are not authorized to process this request")
	}

	if val != otp && !(conf.OtpDemo && otp == "111111") {
		l.LogAttrs(ctx, slog.LevelInfo, "the otp do not match", slog.String("val", val), slog.String("otp", otp))

		attempts, err := db.Redis.HIncrBy(ctx, "otp:"+email, "attempts", 1).Result()
		if err != nil {
			l.LogAttrs(ctx, slog.LevelError, "cannot increment the otp attempt", slog.String("error", err.Error()))
			return "", errors.New("something went wrong")
		}

		if attempts >= conf.OtpAttempts {
			if _, err := db.Redis.Del(ctx, "otp:"+email).Result(); err != nil {
				l.LogAttrs(ctx, slog.LevelError, "cannot destory the otp", slog.String("error", err.Error()))
				return "", errors.New("something went wrong")
			}

			l.LogAttrs(ctx, slog.LevelInfo, "max attempts reached", slog.Int64("attempts", attempts))
			return "", errors.New("you reached the max tentatives")
		}

		return "", errors.New("the OTP does not match")
	}

	q := Query{Email: email}
	res, err := Search(ctx, q, 0, 1)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot verify email existence", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	var uid int

	if res.Total == 0 {
		val, err := db.Redis.Incr(ctx, "user_next_id").Result()

		if err != nil {
			l.LogAttrs(ctx, slog.LevelError, "cannot get the next id", slog.String("error", err.Error()))
			return "", errors.New("something went wrong")
		}

		uid = int(val)
	} else {
		uid = res.Users[0].ID
	}

	sid, err := stringutil.Random()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot generated the session id", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	now := time.Now()

	role := "user"
	if admin := IsAdmin(ctx, email); admin {
		role = "admin"
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, "session:"+sid, "uid", uid, "id", sid, "device", device, "type", "session")
		rdb.Expire(ctx, "session:"+sid, conf.SessionDuration)
		key := fmt.Sprintf("user:%d", uid)
		rdb.HSet(ctx, key,
			"updated_at", now.Unix(),
			// @todo get the lang from the browser and match with the ones on the server
			"lang", conf.DefaultLocale.String(),
		)
		rdb.HSetNX(ctx, key, "id", uid)
		rdb.HSetNX(ctx, key, "email", email)
		rdb.HSetNX(ctx, key, "role", role)
		rdb.HSetNX(ctx, key, "type", "user")
		rdb.HSetNX(ctx, key, "created_at", now.Unix())
		rdb.Del(ctx, "otp:"+email)

		return nil
	}); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot store the data", slog.String("sid", sid), slog.Int("user_id", uid), slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	if role != "admin" && conf.EnableTrackingLog {
		data := map[string]string{
			"sid":   sid,
			"otp":   otp,
			"email": email,
		}

		go tracking.Log(ctx, "login", data)
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the login is successful", slog.String("device", device), slog.String("sid", sid), slog.Int("user_id", uid))

	return sid, nil
}

// Logout destroys the user session.
func (u User) Logout(ctx context.Context) error {
	l := slog.With(slog.String("sid", u.SID))
	l.LogAttrs(ctx, slog.LevelInfo, "trying to logout")

	if u.SID == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate the session id")
		return errors.New("your are not authorized to process this request")
	}

	exist, err := db.Redis.Exists(ctx, "session:"+u.SID).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, err.Error())
		return errors.New("something went wrong")
	}

	if exist == 0 {
		l.LogAttrs(ctx, slog.LevelInfo, "the session does not exist")
		return errors.New("your are not authorized to process this request")
	}

	_, err = db.Redis.Del(ctx, "session:"+u.SID).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot delete the session")
		return errors.New("your are not authorized to process this request")
	}

	if conf.EnableTrackingLog {
		data := map[string]string{
			"sid": u.SID,
		}

		go tracking.Log(ctx, "logout", data)
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the logout is successful")

	return nil
}

// Get the user information from its id
func Get(ctx context.Context, id int) (User, error) {
	l := slog.With(slog.Int("user_id", id))
	l.LogAttrs(ctx, slog.LevelInfo, "trying to get the user")

	if id == 0 {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate the user id")
		return User{}, errors.New("the user is not found")
	}

	data, err := db.Redis.HGetAll(ctx, fmt.Sprintf("user:%d", id)).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot get the user from redis", slog.String("error", err.Error()))
		return User{}, errors.New("something went wrong")
	}

	if data["id"] == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot find the user")
		return User{}, errors.New("the user is not found")
	}

	u, err := parse(ctx, data)

	l.LogAttrs(ctx, slog.LevelInfo, "the user is found")

	return u, err
}

func IsAdmin(ctx context.Context, email string) bool {
	l := slog.With(slog.String("email", email))
	l.LogAttrs(ctx, slog.LevelInfo, "trying to known if the user is admin")

	is, err := db.Redis.SIsMember(ctx, "admins", email).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot retrieve admins from redis", slog.String("error", err.Error()))
		return false
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the user is admin", slog.Bool("yes", is))

	return is
}

func (u User) ToggleDemo(ctx context.Context) (bool, error) {
	l := slog.With(slog.Int("uid", u.ID), slog.Bool("demo", u.Demo))
	l.LogAttrs(ctx, slog.LevelInfo, "toggle demo mode")

	v := "1"
	if u.Demo {
		v = "0"
	}

	_, err := db.Redis.HSet(context.Background(), fmt.Sprintf("user:%d", u.ID), "demo", v).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot toggle demo mode", slog.String("error", err.Error()))
		return u.Demo, err
	}

	l.LogAttrs(ctx, slog.LevelInfo, "demo modo toggled", slog.Bool("demo", !u.Demo))

	return !u.Demo, nil
}

func Search(ctx context.Context, q Query, offset, num int) (SearchResults, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "searching articles", slog.Int("offset", offset), slog.Int("num", num))

	qs := fmt.Sprintf("FT.SEARCH %s @type:{user}", db.UserIdx)

	if q.Email != "" {
		qs += fmt.Sprintf("(@email:{%s})", db.Escape(q.Email))
	}

	qs += fmt.Sprintf(" SORTBY updated_at desc LIMIT %d %d DIALECT 2", offset, num)

	slog.LogAttrs(ctx, slog.LevelInfo, "preparing redis request", slog.String("query", qs))

	args, err := db.SplitQuery(ctx, qs)
	if err != nil {
		return SearchResults{}, err
	}

	cmds, err := db.Redis.Do(ctx, args...).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot run the search query", slog.String("error", err.Error()))
		return SearchResults{}, err
	}

	res := cmds.(map[interface{}]interface{})
	total := res["total_results"].(int64)
	results := res["results"].([]interface{})
	users := []User{}

	for _, value := range results {
		m := value.(map[interface{}]interface{})
		attributes := m["extra_attributes"].(map[interface{}]interface{})
		data := db.ConvertMap(attributes)

		product, err := parse(ctx, data)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the blog", slog.Any("blog", data), slog.String("error", err.Error()))
			continue
		}

		users = append(users, product)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "search done", slog.Int64("results", total))

	return SearchResults{
		Total: int(total),
		Users: users,
	}, nil
}
