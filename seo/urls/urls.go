package urls

import (
	"artisons/db"
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

var u map[string]string = map[string]string{}

func init() {
	ctx := context.Background()
	keys, err := db.Redis.SMembers(ctx, "seo").Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the seo keys", slog.String("error", err.Error()))
		log.Panicln(err)
	}

	cmds, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for _, key := range keys {
			key := "seo:" + key
			rdb.HGetAll(ctx, key)
		}

		return nil
	})

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the seo", slog.String("error", err.Error()))
		log.Panicln((err))
	}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the seo", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		val := cmd.(*redis.MapStringStringCmd).Val()

		Set(val["key"], "title", db.Unescape(val["title"]))
		Set(val["key"], "description", db.Unescape(val["description"]))
		Set(val["key"], "url", db.Unescape(val["url"]))
	}

}

func Set(key, typ, val string) {
	u[key+":"+typ] = val
}

func Get(key, typ string) string {
	return u[key+":"+typ]
}
