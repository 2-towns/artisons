package populate

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func tag(ctx context.Context, pipe redis.Pipeliner) {
	pipe.HSet(ctx, "tag",
		"mens", "Pour les hommes",
		"womens", "Pour les femmes",
		"shoes", "Chaussures",
		"tshirts", "Tshirts",
		"clothes", "Vêtements",
		"socks", "Chaussettes",
		"arabic", "Arabe",
		"books", "Livres",
		"en", "Root",
		"games", "Jeux",
	)

	pipe.ZAdd(ctx, "tag:womens", redis.Z{
		Score:  1,
		Member: "tshirts",
	})

	pipe.ZAdd(ctx, "tag:womens", redis.Z{
		Score:  2,
		Member: "clothes",
	})

	pipe.ZAdd(ctx, "tag:mens", redis.Z{
		Score:  1,
		Member: "tshirts",
	})

	pipe.ZAdd(ctx, "tag:mens", redis.Z{
		Score:  2,
		Member: "books",
	})

	pipe.ZAdd(ctx, "tag:mens", redis.Z{
		Score:  3,
		Member: "clothes",
	})

	pipe.ZAdd(ctx, "tag:books", redis.Z{
		Score:  1,
		Member: "arabic",
	})

	pipe.ZAdd(ctx, "tag:en", redis.Z{
		Score:  1,
		Member: "mens",
	})

	pipe.ZAdd(ctx, "tag:en", redis.Z{
		Score:  2,
		Member: "womens",
	})

	pipe.ZAdd(ctx, "tag:games", redis.Z{
		Score:  1,
		Member: "kids",
	})
}
