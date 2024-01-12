package populate

import (
	"context"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/string/stringutil"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/redis/go-redis/v9"
)

func user(ctx context.Context, pipe redis.Pipeliner, sid string, alive bool) {
	id, err := db.Redis.Incr(ctx, "user_next_id").Result()
	if err != nil {
		return
	}

	now := time.Now()
	key := fmt.Sprintf("user:%d", id)

	pipe.HSet(ctx, key,
		"id", id,
		"email", faker.Email(),
		"updated_at", now.Unix(),
		"created_at", now.Unix(),
	)

	otp, err := stringutil.Random()
	if err != nil {
		return
	}

	pipe.Set(ctx, "otp:"+otp, id, conf.OtpDuration)
	pipe.HSet(ctx, "user", faker.Email(), id)
	pipe.ZAdd(ctx, "users", redis.Z{
		Score:  float64(now.Unix()),
		Member: id,
	})
	pipe.HSet(ctx, fmt.Sprintf("user:%d", id), "lang", conf.DefaultLocale.String())

	if alive {
		pipe.Set(ctx, "auth:"+sid, id, conf.SessionDuration)
		pipe.HSet(ctx, fmt.Sprintf("auth:%s:session", sid), "device", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0")
		pipe.SAdd(ctx, fmt.Sprintf("user:%d:sessions", id), sid)
	}
}
