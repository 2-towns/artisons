// tests gather test utilites
package tests

import (
	"artisons/db"
	"artisons/http/contexts"
	"artisons/string/stringutil"
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/text/language"
)

func Del(ctx context.Context, prefix string) {
	keys, err := db.Redis.Keys(ctx, fmt.Sprintf("%s*", prefix)).Result()
	if err != nil {
		log.Fatal(err)
	}

	pipe := db.Redis.Pipeline()

	for _, key := range keys {
		pipe.Del(ctx, key)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func ImportData(ctx context.Context, file string) {
	lines := db.ParseData(ctx, file)
	pipe := db.Redis.Pipeline()

	for _, line := range lines {
		pipe.Do(ctx, line...)
	}

	cmds, err := pipe.Exec(ctx)

	for _, cmd := range cmds {
		if cmd.Err() != nil {
			log.Println(cmd.String(), cmd.Err())
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}

func Context() context.Context {
	var ctx context.Context = context.WithValue(context.Background(), contexts.RequestID, fmt.Sprintf("%d", time.Now().UnixMilli()))
	ctx = context.WithValue(ctx, contexts.Locale, language.English)

	rid, _ := stringutil.Random()
	ctx = context.WithValue(ctx, contexts.RequestID, rid)

	return context.WithValue(ctx, contexts.Device, fmt.Sprintf("%d", time.Now().UnixMilli()))
}
