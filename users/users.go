// Package products provide everything around users
package users

import (
	"context"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"log"
	"strconv"
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

// Persist an user into redis
func (u User) Persist(password string) error {
	if password == "" {
		return fmt.Errorf("input_validation_fail: password is required")
	}

	v := validator.New()
	err := v.Struct(u)
	if err != nil {
		return fmt.Errorf("input_validation_fail: error when validation user %s", err.Error())
	}

	ctx := context.Background()

	existing, err := db.Redis.HGet(ctx, "user", u.Username).Result()
	if existing != "" && err != nil {
		return fmt.Errorf("input_validation_fail:username already exists")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return fmt.Errorf("sequence_fail: error when generating password hash %s", err.Error())
	}

	hash := string(bytes)
	id, err := db.Redis.Incr(ctx, "user_next_id").Result()
	if err != nil {
		id = 1
	}

	now := time.Now()
	key := fmt.Sprintf("user:%d", id)
	pipe := db.Redis.Pipeline()
	pipe.HSet(ctx, key, "id", id)
	pipe.HSet(ctx, key, "username", u.Username)
	pipe.HSet(ctx, key, "email", u.Email)
	/*pipe.HSet(ctx, key, "lastname", u.Address.Lastname)
	pipe.HSet(ctx, key, "firstname", u.Firstname)
	pipe.HSet(ctx, key, "address", u.Address)
	pipe.HSet(ctx, key, "city", u.City)
	pipe.HSet(ctx, key, "complemnetary", u.Complementary)
	pipe.HSet(ctx, key, "zipcode", u.Zipcode)
	pipe.HSet(ctx, key, "phone", u.Phone)*/
	pipe.HSet(ctx, key, "hash", hash)
	pipe.HSet(ctx, key, "updated_at", now.Format(time.RFC3339))
	pipe.HSet(ctx, key, "created_at", now.Format(time.RFC3339))
	pipe.HSet(ctx, "user", u.Username, id)
	pipe.ZAdd(ctx, "users", redis.Z{
		Score:  float64(now.Unix()),
		Member: id,
	})

	_, err = pipe.Exec(ctx)

	return err
}

// List returns the user list in the application
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
		return users, fmt.Errorf("sequence_fail: go error from redis %s", err.Error())
	}

	for _, cmd := range cmds {
		m := cmd.(*redis.MapStringStringCmd).Val()

		id, err := strconv.ParseInt(m["id"], 10, 32)
		if err != nil {
			log.Printf("sequence_fail: error when parsing id %s", m["id"])

			continue
		}

		createdAt, err := time.Parse(time.RFC3339, m["created_at"])
		if err != nil {
			log.Printf("sequence_fail: error when parsing created_at %s", m["created_at"])

			continue
		}

		updatedAt, err := time.Parse(time.RFC3339, m["updated_at"])
		if err != nil {
			log.Printf("sequence_fail: error when parsing created_at %s", m["updated_at"])

			continue
		}

		users = append(users, User{
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
		})
	}

	return users, nil
}
