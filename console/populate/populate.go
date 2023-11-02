// Package populate provide script to populate date into Redis
package populate

import (
	"context"
	"gifthub/db"
	"gifthub/orders"
	"gifthub/string/stringutil"
	"gifthub/tests"
	"gifthub/users"

	"github.com/go-faker/faker/v4"
)

// Run the populate script. It will flush the database first
func Run() error {
	ctx := context.Background()
	c := tests.Context()

	db.Redis.FlushDB(ctx)

	pid, _ := stringutil.Random()
	db.Redis.HSet(ctx, "product:"+pid, "status", "online")

	magic, err := users.MagicCode(c, faker.Email())
	if err != nil {
		return err
	}

	_, err = users.Login(c, magic, "Mozilla/5.0 Gecko/20100101 Firefox/115.0")
	if err != nil {
		return err
	}

	u, err := users.List(c, 0)
	if err != nil {
		return err
	}

	user := u[0]

	o := orders.Order{
		UID:      user.ID,
		Delivery: "collect",
		Payment:  "cash",
		Products: map[string]int64{pid: 1},
	}

	_, err = o.Save(c)
	if err != nil {
		return err
	}

	return nil
}
