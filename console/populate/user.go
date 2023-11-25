package populate

import (
	"context"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/locales"
	"gifthub/string/stringutil"
	"gifthub/users"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/redis/go-redis/v9"
)

func User(ctx context.Context, sid string, alive bool) (users.User, error) {
	id, err := db.Redis.Incr(ctx, "user_next_id").Result()
	if err != nil {
		return users.User{}, err
	}

	now := time.Now()
	key := fmt.Sprintf("user:%d", id)

	_, err = db.Redis.HSet(ctx, key,
		"id", id,
		"email", faker.Email(),
		"updated_at", now.Unix(),
		"created_at", now.Unix(),
	).Result()
	if err != nil {
		return users.User{}, err
	}

	otp, err := stringutil.Random()
	if err != nil {
		return users.User{}, err
	}

	if _, err = db.Redis.Set(ctx, "otp:"+otp, id, conf.OtpDuration).Result(); err != nil {
		return users.User{}, err
	}

	if _, err = db.Redis.HSet(ctx, "user", faker.Email(), id).Result(); err != nil {
		return users.User{}, err
	}

	if _, err = db.Redis.ZAdd(ctx, "users", redis.Z{
		Score:  float64(now.Unix()),
		Member: id,
	}).Result(); err != nil {
		return users.User{}, err
	}

	if alive {
		if _, err = db.Redis.Set(ctx, "auth:"+sid, id, conf.SessionDuration).Result(); err != nil {
			return users.User{}, err
		}

		if _, err = db.Redis.HSet(ctx, fmt.Sprintf("auth:%s:session", sid), "device", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0").Result(); err != nil {
			return users.User{}, err
		}

		if _, err = db.Redis.SAdd(ctx, fmt.Sprintf("user:%d:sessions", id), sid).Result(); err != nil {
			return users.User{}, err
		}
	}

	if _, err = db.Redis.HSet(ctx, fmt.Sprintf("user:%d", id), "lang", locales.Default.String()).Result(); err != nil {
		return users.User{}, err
	}

	if _, err = db.Redis.SAdd(ctx, "admins", "hello@world.com", "lock@world.com").Result(); err != nil {
		return users.User{}, err
	}

	return users.User{
		ID: id,
	}, err
}
