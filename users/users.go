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

	"github.com/redis/go-redis/v9"
)

// User is the user representation in the application
type User struct {
	Email         string
	ID            int64
	Lastname      string
	Firstname     string
	Address       string
	City          string
	Complementary string
	Zipcode       string
	Phone         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
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
			ID:            id,
			Email:         m["email"],
			Lastname:      m["lastname"],
			Firstname:     m["firstname"],
			Address:       m["address"],
			Complementary: m["complementary"],
			Zipcode:       m["zipcode"],
			City:          m["city"],
			Phone:         m["phone"],
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
		})
	}

	return users, nil
}
