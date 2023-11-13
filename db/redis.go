// Package db provides redis storage
package db

import (
	"context"
	"fmt"
	"gifthub/conf"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

// Redis is the client to use for Redis interactions
var Redis = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",                 // no password set
	DB:       conf.DatabaseIndex, // use default DB
})

var ProductIdx = "product-idx"
var OrderIdx = "order-idx"

func ProductIndex(ctx context.Context) error {
	_, err := Redis.Do(
		ctx,
		"FT.DROPINDEX",
		ProductIdx,
	).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot make remove the previous product index", slog.String("error", err.Error()))
	}

	_, err = Redis.Do(
		ctx,
		"FT.CREATE", ProductIdx,
		"ON", "HASH",
		"PREFIX", "1", "product:",
		"SCHEMA",
		"title", "TEXT",
		"sku", "TAG",
		"description", "TEXT",
		"price", "NUMERIC", "SORTABLE",
		"created_at", "NUMERIC", "SORTABLE",
		"updated_at", "NUMERIC", "SORTABLE",
		"tags", "TAG", "SEPARATOR", ";",
		"status", "TAG",
		"meta", "TAG", "SEPARATOR", ";",
	).Result()

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot make product migration", slog.String("error", err.Error()))
	}

	return err
}

func OrderIndex(ctx context.Context) error {
	_, err := Redis.Do(
		ctx,
		"FT.DROPINDEX",
		OrderIdx,
	).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot make remove the previous order index", slog.String("error", err.Error()))
	}

	_, err = Redis.Do(
		ctx,
		"FT.CREATE", OrderIdx,
		"ON", "HASH",
		"PREFIX", "1", "order:",
		"SCHEMA",
		"status", "TAG",
		"delivery", "TAG",
		"payment", "TAG",
		"created_at", "NUMERIC", "SORTABLE",
		"updated_at", "NUMERIC", "SORTABLE",
	).Result()

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot make order migration", slog.String("error", err.Error()))
	}

	return err
}

// ConvertMap converts the redis search result to an map
func ConvertMap(m map[interface{}]interface{}) map[string]string {
	v := map[string]string{}

	for key, value := range m {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		v[strKey] = strValue
	}

	return v
}

/*func SubscribeToExpireKeys() {
	ctx := context.Background()
	if _, err := Redis.ConfigSet(ctx, "notify-keyspace-events", "KEA").Result(); err != nil {
		log.Panicln(err)
	}

	pubsub := Redis.PSubscribe(ctx, "__key*__:auth:*")
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
		}
		parts := strings.Split(msg.Channel, "auth:")
		fmt.Println(parts[1])
	}
}
*/
