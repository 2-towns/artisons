// Package db provides redis storage
package db

import (
	"gifthub/conf"

	"github.com/redis/go-redis/v9"
)

// Redis is the client to use for Redis interactions
var Redis = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",                 // no password set
	DB:       conf.DatabaseIndex, // use default DB
})

var SearchIdx = "product-idx"

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
