// Package users provide everything around users
package users

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// User is the user representation in the application
type User struct {
	Email     string `validate:"required,email"`
	Username  string `validate:"required,alpha,lowercase"`
	ID        int64
	Address   Address
	Hash      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Address struct {
	Lastname      string
	Firstname     string
	City          string
	Complementary string
	Zipcode       string
	Phone         string
}

// Persist an user into HSET with the storage key: user:ID.
// An extra link is create between the username and the ID,
// with the key user:USERNAME.
// The user is also added in a sorted set with the timestamp as score.
// All the operations are executed in a single transaction.
//
// The email, username and password are required to create an user.
// If one of those fields is empty, an error occurs.
//
// A hash is generated from this password and stored in Redis.
//
// The user ID is an incremented field in Redis and returned.
func (u User) Persist(password string) (int64, error) {
	if password == "" {
		log.Printf("input_validation_fail: password is required")
		return 0, errors.New("user_password_required")
	}

	v := validator.New()
	err := v.Struct(u)
	if err != nil {
		log.Printf("input_validation_fail: error when validation user %s", err.Error())
		e := err.(validator.ValidationErrors)[0]
		f := strings.ToLower(e.StructField())
		return 0, fmt.Errorf("user_%s_invalid", f)
	}

	ctx := context.Background()
	existing, err := db.Redis.HGet(ctx, "user", u.Username).Result()
	if existing != "" && err == nil {
		log.Printf("input_validation_fail: username already exists")
		return 0, errors.New("user_username_exists")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when generating password hash %s", err.Error())
		return 0, errors.New("something_went_wrong")
	}

	hash := string(bytes)
	id, err := db.Redis.Incr(ctx, "user_next_id").Result()
	if err != nil {
		log.Printf("ERROR: sequence_fail: %s", err.Error())
		return 0, errors.New("something_went_wrong")
	}

	now := time.Now()
	key := fmt.Sprintf("user:%d", id)

	if _, err = db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, key,
			"id", id, "username",
			u.Username, "email",
			u.Email, "hash", hash,
			"updated_at", now.Format(time.RFC3339),
			"created_at", now.Format(time.RFC3339),
		)
		rdb.HSet(ctx, "user", u.Username, id)
		rdb.ZAdd(ctx, "users", redis.Z{
			Score:  float64(now.Unix()),
			Member: id,
		})
		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return 0, errors.New("something_went_wrong")
	}

	/*pipe.HSet(ctx, key, "lastname", u.Address.Lastname)
	pipe.HSet(ctx, key, "firstname", u.Firstname)
	pipe.HSet(ctx, key, "address", u.Address)
	pipe.HSet(ctx, key, "city", u.City)
	pipe.HSet(ctx, key, "complemnetary", u.Complementary)
	pipe.HSet(ctx, key, "zipcode", u.Zipcode)
	pipe.HSet(ctx, key, "phone", u.Phone)*/

	log.Printf("WARN: sensitive_create: a new user is created with id %d\n", u.ID)

	return id, nil
}

// Delete an user at three levels in a single transation:
//   - the hset data
//   - the ID link
//   - the member in user list
func (u User) Delete() error {
	key := fmt.Sprintf("user:%d", u.ID)
	ctx := context.Background()

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Del(ctx, key)
		rdb.HDel(ctx, "user", u.Username)
		rdb.ZRem(ctx, "users", u.ID)
		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return errors.New("something_went_wrong")
	}

	log.Printf("WARN: sensitive_delete: the user %d deleted\n", u.ID)

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
		ID:       id,
		Email:    m["email"],
		Username: m["username"],
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
