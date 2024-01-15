// Package db provides redis storage
package db

import (
	"context"
	"fmt"
	"gifthub/conf"
	"log"
	"strings"

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
var BlogIdx = "blog-idx"
var LocaleIdx = "locale-idx"

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

// Escape escapes the key characters used in Redis Search by
// adding backslashes.
func Escape(value string) string {
	replacements := map[string]string{
		",": "\\,",
		".": "\\.",
		"<": "\\<",
		">": "\\>",
		"{": "\\{",
		"}": "\\}",
		"[": "\\[",
		"]": "\\]",
		`"`: "\\\"",
		":": "\\:",
		";": "\\;",
		"!": "\\!",
		"@": "\\@",
		"#": "\\#",
		"$": "\\$",
		"%": "\\%",
		"^": "\\^",
		"&": "\\&",
		"*": "\\*",
		"(": "\\(",
		")": "\\)",
		"-": "\\-",
		"+": "\\+",
		"=": "\\=",
		"~": "\\~",
	}

	s := strings.Trim(value, " ")

	for key, v := range replacements {
		s = strings.ReplaceAll(s, key, v)
	}

	return s
}

func Unescape(s string) string {
	return strings.ReplaceAll(s, "\\", "")
}

func Run(ctx context.Context, args []interface{}) error {
	log.Println(args)
	r, err := Redis.Do(ctx, args...).Result()

	log.Println(r)

	return err
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
