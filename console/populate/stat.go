package populate

import (
	"context"
	"fmt"
	"gifthub/db"
	"gifthub/string/stringutil"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/redis/go-redis/v9"
)

func Stats(ctx context.Context) error {
	now := time.Now()

	pipe := db.Redis.Pipeline()

	products := []string{"test", "test-1", "test-2", "test-3"}

	devices := []string{}

	for i := 0; i < 100; i++ {
		cid, err := stringutil.Random()
		if err != nil {
			return err
		}

		devices = append(devices, cid)
	}

	for i := 0; i < 5000; i++ {
		i := rand.Intn(30)
		score := now.AddDate(0, 0, -i)

		idx := rand.Intn(len(devices))
		cid := devices[idx]

		rid, err := stringutil.Random()
		if err != nil {
			return err
		}

		oid, err := stringutil.Random()
		if err != nil {
			return err
		}

		key := ":" + score.Format("20060102")

		pipe.ZAdd(ctx, "statvisits"+key, redis.Z{
			Score:  float64(score.Unix()),
			Member: rid,
		})

		pipe.ZAdd(ctx, "statuniquevisits"+key, redis.Z{
			Score:  float64(score.Unix()),
			Member: cid,
		})

		pipe.ZAdd(ctx, "statnewusers"+key, redis.Z{
			Score:  float64(score.Unix()),
			Member: faker.Email(),
		})

		pipe.ZAdd(ctx, "statactiveusers"+key, redis.Z{
			Score:  float64(score.Unix()),
			Member: faker.Email(),
		})

		pipe.ZAdd(ctx, "statorders"+key, redis.Z{
			Score:  float64(score.Unix()),
			Member: oid,
		})

		pipe.ZAdd(ctx, "statrevenues"+key, redis.Z{
			Score:  float64(score.Unix()),
			Member: fmt.Sprintf("%f:%s", rand.Float32()*100, oid),
		})

		pidx := rand.Intn(len(products))

		pipe.ZAdd(ctx, "statsoldproducts"+key, redis.Z{
			Score:  float64(score.Unix()),
			Member: fmt.Sprintf("%s:%d:%s", products[pidx], 1, oid),
		})

		pipe.ZAdd(ctx, "statvisitproduct"+key, redis.Z{
			Score:  float64(score.Unix()),
			Member: products[pidx] + ":" + rid,
		})

		pipe.ZAdd(ctx, "statuniquevisitproduct:"+key, redis.Z{
			Score:  float64(score.Unix()),
			Member: products[pidx] + cid,
		})
	}

	_, err := pipe.Exec(ctx)

	return err
}
