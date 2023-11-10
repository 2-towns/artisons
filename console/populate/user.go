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
		"updated_at", now.Format(time.RFC3339),
		"created_at", now.Format(time.RFC3339),
	).Result()
	if err != nil {
		return users.User{}, err
	}

	magic, err := stringutil.Random()
	if err != nil {
		return users.User{}, err
	}

	_, err = db.Redis.Set(ctx, "magic:"+magic, id, conf.MagicCodeDuration).Result()
	if err != nil {
		return users.User{}, err
	}

	_, err = db.Redis.HSet(ctx, "user", faker.Email(), id).Result()
	if err != nil {
		return users.User{}, err
	}
	_, err = db.Redis.ZAdd(ctx, "users", redis.Z{
		Score:  float64(now.Unix()),
		Member: id,
	}).Result()
	if err != nil {
		return users.User{}, err
	}

	if alive {
		_, err = db.Redis.Set(ctx, "auth:"+sid, id, conf.SessionDuration).Result()
		if err != nil {
			return users.User{}, err
		}

		_, err = db.Redis.HSet(ctx, fmt.Sprintf("auth:%s:session", sid), "device", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0").Result()
		if err != nil {
			return users.User{}, err
		}

		_, err = db.Redis.SAdd(ctx, fmt.Sprintf("user:%d:sessions", id), sid).Result()
		if err != nil {
			return users.User{}, err
		}
	}

	_, err = db.Redis.HSet(ctx, fmt.Sprintf("user:%d", id), "lang", locales.Default.String()).Result()
	if err != nil {
		return users.User{}, err
	}

	return users.User{
		ID: id,
	}, err
}
