package populate

import (
	"context"
	"gifthub/db"

	"github.com/redis/go-redis/v9"
)

func Tag(ctx context.Context) error {
	if _, err := db.Redis.HSet(ctx, "tag",
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
	).Result(); err != nil {
		return err
	}

	if _, err := db.Redis.ZAdd(ctx, "tag:womens", redis.Z{
		Score:  1,
		Member: "tshirts",
	}).Result(); err != nil {
		return err
	}

	if _, err := db.Redis.ZAdd(ctx, "tag:womens", redis.Z{
		Score:  2,
		Member: "clothes",
	}).Result(); err != nil {
		return err
	}

	if _, err := db.Redis.ZAdd(ctx, "tag:mens", redis.Z{
		Score:  1,
		Member: "tshirts",
	}).Result(); err != nil {
		return err
	}

	if _, err := db.Redis.ZAdd(ctx, "tag:mens", redis.Z{
		Score:  2,
		Member: "books",
	}).Result(); err != nil {
		return err
	}

	if _, err := db.Redis.ZAdd(ctx, "tag:mens", redis.Z{
		Score:  3,
		Member: "clothes",
	}).Result(); err != nil {
		return err
	}

	if _, err := db.Redis.ZAdd(ctx, "tag:books", redis.Z{
		Score:  1,
		Member: "arabic",
	}).Result(); err != nil {
		return err
	}

	if _, err := db.Redis.ZAdd(ctx, "tag:en", redis.Z{
		Score:  1,
		Member: "mens",
	}).Result(); err != nil {
		return err
	}

	if _, err := db.Redis.ZAdd(ctx, "tag:en", redis.Z{
		Score:  2,
		Member: "womens",
	}).Result(); err != nil {
		return err
	}

	if _, err := db.Redis.ZAdd(ctx, "tag:games", redis.Z{
		Score:  1,
		Member: "kids",
	}).Result(); err != nil {
		return err
	}

	return nil
}
