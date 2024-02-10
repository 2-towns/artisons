package tree

import (
	"artisons/db"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/redis/go-redis/v9"
)

var Tree = []Leaf{}

type Bud struct {
	Key      string `validate:"required,alphanum"`
	Label    string `validate:"required"`
	Image    string
	Children []string
}

type Leaf struct {
	Bud
	Branches []*Leaf
}

func init() {
	ctx := context.Background()

	t, err := Build(ctx)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the categories", slog.String("error", err.Error()))
		log.Fatalln(err)
	}

	Tree = t
}

func Build(ctx context.Context) ([]Leaf, error) {
	leaves := []Leaf{}

	roots, err := db.Redis.ZRange(ctx, "tags:root", 0, -1).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the roots", slog.String("error", err.Error()))
		return []Leaf{}, errors.New("something went wrong")
	}

	buds, err := db.Redis.ZRange(ctx, "tags", 0, 9999).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the tags", slog.String("error", err.Error()))
		return []Leaf{}, errors.New("something went wrong")
	}

	cmds, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for _, val := range buds {
			rdb.HGetAll(ctx, "tag:"+val)
		}

		return nil
	})

	if err != nil && err.Error() != "redis: nil" {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the tag meta data", slog.String("error", err.Error()))
		return []Leaf{}, errors.New("something went wrong")
	}

	bud := map[string]Bud{}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil && cmd.Err().Error() != "redis: nil" {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the tag meta data", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		val := cmd.(*redis.MapStringStringCmd).Val()

		bud[val["key"]] = Bud{
			Key:      val["key"],
			Label:    val["label"],
			Image:    val["image"],
			Children: strings.Split(val["children"], ";"),
		}
	}

	for _, val := range roots {
		leaf := Leaf{
			Bud: Bud{
				Key:   val,
				Label: bud[val].Label,
				Image: bud[val].Image,
			},
			Branches: []*Leaf{},
		}

		for _, branch := range bud[val].Children {
			l := Leaf{
				Bud: Bud{
					Key:   branch,
					Label: bud[branch].Label,
					Image: bud[branch].Image,
				},
			}

			lb := strings.Split(bud[branch].Key, ";")

			for _, b := range lb {
				l.Branches = append(l.Branches, &Leaf{
					Bud: Bud{
						Key:   b,
						Label: bud[branch].Label,
						Image: bud[branch].Image,
					},
				})
			}

			leaf.Branches = append(leaf.Branches, &l)
		}

		leaves = append(leaves, leaf)
	}

	Tree = leaves

	return leaves, nil
}
