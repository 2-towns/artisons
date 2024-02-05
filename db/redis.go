// Package db provides redis storage
package db

import (
	"artisons/conf"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

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
var UserIdx = "user-idx"
var SessionIdx = "session-idx"
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
	s := strings.ReplaceAll(value, "-", "\\-")
	s = strings.ReplaceAll(s, "@", "\\@")
	s = strings.ReplaceAll(s, ".", "\\.")
	return strings.ReplaceAll(s, "'", "\\'")
}

// SearchValue replaces the space by the caracter |.
// There is no need to escape other characters here is the quote from Redis:
// The Redis protocol has no concept of string escaping, so injection
// is impossible under normal circumstances using a normal client library.
// The protocol uses prefixed-length strings and is completely binary safe.
// https://github.com/RediSearch/RediSearch/issues/259
// https://redis.io/docs/management/security/
func SearchValue(value string) string {
	esc := Escape(value)
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(esc, "|")
}

func Unescape(s string) string {
	return strings.ReplaceAll(s, "\\", "")
}

func Run(ctx context.Context, args []interface{}) error {
	_, err := Redis.Do(ctx, args...).Result()

	return err
}

func SplitQuery(ctx context.Context, s string) ([]interface{}, error) {
	args := []interface{}{}
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = ' '
	fields, err := r.Read()

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the string", slog.String("string", s), slog.String("error", err.Error()))
		return []interface{}{}, err
	}

	for _, val := range fields {
		if val != "" {
			args = append(args, val)
		}
	}

	return args, nil
}

func ParseData(ctx context.Context, file string) [][]interface{} {
	f, err := os.ReadFile(file)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot open the file", slog.String("error", err.Error()))
		log.Fatal(err)
	}

	cmds := strings.Split(string(f), "\n")
	lines := [][]interface{}{}

	for _, line := range cmds {
		if line == "" {
			continue
		}

		args := []interface{}{}

		l := strings.Replace(line, "20060102", time.Now().Format("20060102"), -1)
		l = strings.Replace(l, "1136160000", fmt.Sprintf("%d", time.Now().Unix()), -1)

		r := csv.NewReader(strings.NewReader(l))
		r.Comma = ' '
		fields, err := r.Read()

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "error when parsing line", slog.String("line", line), slog.String("error", err.Error()))
			log.Fatalln(err)
		}

		for _, val := range fields {
			if val != "" {
				args = append(args, val)
			}
		}

		lines = append(lines, args)
	}

	return lines
}

/*func SubscribeToExpireKeys() {
	if _, err := Redis.ConfigSet(ctx, "notify-keyspace-events", "KEA").Result(); err != nil {
		log.Panicln(err)
	}

	pubsub := Redis.PSubscribe(ctx, "__key*__:session:*")
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
		}
		parts := strings.Split(msg.Channel, "session:")
		fmt.Println(parts[1])
	}
}
*/
