package populate

import (
	"context"
	"gifthub/conf"

	"github.com/redis/go-redis/v9"
)

func cart(ctx context.Context, pipe redis.Pipeliner, cid string, uid int64) {
	pipe.HSet(ctx, "cart:"+cid, "cid", "CAR1")
	pipe.Set(ctx, "cart:"+cid+":user", uid, conf.CartDuration)
}
