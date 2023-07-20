// Package populate provide script to populate date into Redis
package populate

import (
	"context"
	"gifthub/db"
	"gifthub/users"
	"strings"

	"github.com/go-faker/faker/v4"
)

// Run the populate script. It will flush the database first
func Run() error {
	ctx := context.Background()

	db.Redis.FlushDB(ctx)

	u := users.User{
		Email:    faker.Email(),
		Username: strings.ToLower(faker.Username()),
	}

	_, err := u.Persist("test")

	if err != nil {
		return err
	}

	u = users.User{
		Email:    faker.Email(),
		Username: "toto",
	}

	_, err = u.Persist("test")

	return err
}
