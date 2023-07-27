// Package populate provide script to populate date into Redis
package populate

import (
	"context"
	"gifthub/db"
	"gifthub/users"
	"github.com/go-faker/faker/v4"
)

// Run the populate script. It will flush the database first
func Run() error {
	ctx := context.Background()

	db.Redis.FlushDB(ctx)

	magic, err := users.MagicCode(faker.Email())
	if err != nil {
		return err
	}

	_, err = users.Login(magic, "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if err != nil {
		return err
	}

	return nil
}
