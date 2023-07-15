// Package populate provide script to populate date into Redis
package populate

import (
	"context"
	"gifthub/db"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/redis/go-redis/v9"
)

// Run the populate script. It will flush the database first
func Run() error {
	ctx := context.Background()

	db.Redis.FlushDB(ctx)

	a := faker.GetRealAddress()

	now := time.Now()
	key := "user:1"

	_, err := db.Redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, "users:next_id", 2, 0)
		pipe.HSet(ctx, key, "id", "1")
		pipe.HSet(ctx, key, "email", faker.Email())
		pipe.HSet(ctx, key, "firstname", faker.FirstName())
		pipe.HSet(ctx, key, "lastname", faker.LastName())
		pipe.HSet(ctx, key, "address", a.Address)
		pipe.HSet(ctx, key, "city", a.City)
		pipe.HSet(ctx, key, "complementary", a.Address)
		pipe.HSet(ctx, key, "zipcode", a.PostalCode)
		pipe.HSet(ctx, key, "phone", faker.Phonenumber())
		pipe.HSet(ctx, key, "updated_at", now.Format(time.RFC3339))
		pipe.HSet(ctx, key, "created_at", now.Format(time.RFC3339))
		pipe.ZAdd(ctx, "users", redis.Z{
			Score:  float64(now.Unix()),
			Member: "1",
		})

		return nil
	})

	return err
}
