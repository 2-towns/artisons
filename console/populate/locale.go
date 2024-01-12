package populate

import (
	"context"

	"github.com/redis/go-redis/v9"
	"golang.org/x/text/language"
)

func locale(ctx context.Context, pipe redis.Pipeliner, tag language.Tag, key, value string) {
	pipe.HSet(ctx, "locale:"+tag.String(), key, value)
}
