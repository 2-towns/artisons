package users

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/string/stringutil"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func saveSID(id int64) (string, error) {
	sessionID, err := stringutil.Random()
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when generating a new session ID %s", err.Error())
		return "", errors.New("something_went_wrong")
	}

	ctx := context.Background()
	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(ctx, "auth:"+sessionID, id, conf.SessionDuration)
		/*rdb.Set(ctx, fmt.Sprintf("user:%d", id), "session:"+sessionID, conf.SessionDuration)*/
		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return "", errors.New("something_went_wrong")
	}

	return sessionID, nil
}

func findBySessionID(sessionID string) (User, error) {
	ctx := context.Background()
	id, err := db.Redis.Get(ctx, "auth:"+sessionID).Result()
	if err != nil {
		log.Printf("WARN: authz_fail: error when looking for session %s %s", sessionID, err.Error())
		return User{}, err
	}

	m, err := db.Redis.HGetAll(ctx, "user:"+id).Result()
	if err != nil {
		log.Printf("ERROR: sequence_fail: error when loading redis data for session %s %s", sessionID, err.Error())
		return User{}, err
	}

	u, err := parseUser(m)
	if err != nil {
		log.Printf("ERROR: sequence_fail: error parsing data for user %s %s ", id, err.Error())
		return User{}, err
	}

	return u, err
}

func saveUser(email string, magic string) (int64, error) {
	ctx := context.Background()

	id, err := db.Redis.Incr(ctx, "user_next_id").Result()
	if err != nil {
		log.Printf("ERROR: sequence_fail: %s", err.Error())
		return 0, errors.New("something_went_wrong")
	}

	now := time.Now()
	key := fmt.Sprintf("user:%d", id)

	if _, err = db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, key,
			"id", id,
			"email", email,
			"updated_at", now.Format(time.RFC3339),
			"created_at", now.Format(time.RFC3339),
		)
		rdb.Set(ctx, "magic:"+magic, id, conf.MagicCodeDuration)
		rdb.HSet(ctx, "user", email, id)
		rdb.ZAdd(ctx, "users", redis.Z{
			Score:  float64(now.Unix()),
			Member: id,
		})
		return nil
	}); err != nil {
		log.Printf("ERROR: sequence_fail: error when storing in redis %s", err.Error())
		return 0, errors.New("something_went_wrong")
	}

	return id, nil
}
