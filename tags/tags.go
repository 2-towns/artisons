package tags

import (
	"context"
	"errors"
	"fmt"
	"gifthub/db"
	"gifthub/validators"
	"log"
	"log/slog"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type Tag struct {
	Key       string `validate:"required,alphanum"`
	Label     string `validate:"required"`
	Image     string
	Children  []string
	Root      bool
	Score     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

var Tree = []Leaf{}

type Leaf struct {
	Tag
	Branches []*Leaf
}

type ListResults struct {
	Total int
	Tags  []Tag
}

func init() {
	ctx := context.Background()

	t, err := tree(ctx)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the categories", slog.String("error", err.Error()))
		log.Fatalln(err)
	}

	Tree = t
}

func (p Tag) Validate(ctx context.Context) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "validating a tag")

	if err := validators.V.Struct(p); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot validate the tag", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input:%s", low)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "tag validated")

	return nil
}

func Exists(ctx context.Context, key string) (bool, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "checking existence", slog.String("key", key))

	exists, err := db.Redis.Exists(ctx, "tag:"+key).Result()

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot check tags existence")
		return false, errors.New("something went wrong")
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "tag existence", slog.String("key", key), slog.Int64("exists", exists))

	return exists > 0, nil
}

func (t Tag) Save(ctx context.Context) (string, error) {
	l := slog.With(slog.String("tag", t.Key))
	l.LogAttrs(ctx, slog.LevelInfo, "adding a new tag")

	children := strings.Join(t.Children, ";")
	now := time.Now()

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, "tag", t.Key, children)
		rdb.HSet(ctx, "tag:"+t.Key,
			"image", t.Image,
			"children", children,
			"label", t.Label,
			"updated_at", now.Unix(),
		)

		rdb.HSetNX(ctx, "tag:"+t.Key, "created_at", now.Unix())

		if t.Root {
			rdb.ZAdd(ctx, "tags", redis.Z{
				Score:  float64(t.Score),
				Member: t.Key,
			})
		}

		t, err := tree(ctx)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot build the categories tree", slog.String("error", err.Error()))
			return err
		}

		Tree = t

		return nil

	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return "", errors.New("something went wrong")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "tag saved successfully")

	return t.Key, nil
}

func parse(ctx context.Context, data map[string]string) (Tag, error) {
	tag := Tag{
		Key:      data["key"],
		Label:    data["label"],
		Image:    data["image"],
		Children: strings.Split(data["children"], ";"),
	}

	return tag, nil
}

func Find(ctx context.Context, key string) (Tag, error) {
	l := slog.With(slog.String("id", key))
	l.LogAttrs(ctx, slog.LevelInfo, "looking for tag")

	if key == "" {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot validate empty tag key")
		return Tag{}, errors.New("input:id")
	}

	if exists, err := db.Redis.Exists(ctx, "tag:"+key).Result(); exists == 0 || err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, "cannot find the tag")
		return Tag{}, errors.New("oops the data is not found")
	}

	data, err := db.Redis.HGetAll(ctx, "tag:"+key).Result()
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot find the tag", slog.String("error", err.Error()))
		return Tag{}, err
	}

	data["key"] = key
	tag, err := parse(ctx, data)
	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot parse the tag", slog.String("error", err.Error()))
		return Tag{}, err
	}

	score, err := db.Redis.ZScore(ctx, "tags", key).Result()
	if err == nil {
		tag.Root = true
		tag.Score = int(score)
	}

	l.LogAttrs(ctx, slog.LevelInfo, "the tag is found")

	return tag, nil
}

func List(ctx context.Context, offset, num int) (ListResults, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "listing tags")

	hashes, err := db.Redis.HGetAll(context.Background(), "tag").Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the tags", slog.String("error", err.Error()))
		return ListResults{}, errors.New("something went wrong")
	}

	tags := map[string]Tag{}

	for key, value := range hashes {
		tags[key] = Tag{
			Key:      key,
			Children: strings.Split(value, ";"),
		}
	}

	cmds, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for key := range hashes {
			rdb.HGet(ctx, "tag:"+key, "updated_at")
		}

		return nil
	})

	if err != nil && err.Error() != "redis: nil" {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the tag meta data", slog.String("error", err.Error()))
		return ListResults{}, errors.New("something went wrong")
	}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil && cmd.Err().Error() != "redis: nil" {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the tag meta data", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		val := cmd.(*redis.StringCmd).Val()

		if val == "" {
			continue
		}

		k := strings.Replace(key, "tag:", "", 1)

		updatedAt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the updated at", slog.String("error", err.Error()), slog.String("updated_at", val))
			return ListResults{}, errors.New("input:updated_at")
		}

		t := tags[k]
		t.UpdatedAt = time.Unix(updatedAt, 0)
		tags[k] = t
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "found tags", slog.Int("length", len(tags)))

	values := maps.Values(tags)
	o := math.Min(float64(offset), float64(len(hashes)))
	n := math.Min(float64(num), float64(len(hashes)))

	sort.Slice(values, func(i, j int) bool {
		return values[i].UpdatedAt.Unix() > values[j].UpdatedAt.Unix()
	})

	return ListResults{
		Tags:  values[int(o):int(n)],
		Total: len(tags),
	}, nil
}

func Delete(ctx context.Context, key string) error {
	l := slog.With(slog.String("tag", key))
	l.LogAttrs(ctx, slog.LevelInfo, "deleting  tag")

	if key == "" {
		slog.LogAttrs(ctx, slog.LevelInfo, "the name cannot be empty")
		return errors.New("input:name")
	}

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HDel(ctx, "tag", key)
		rdb.Del(ctx, "tag:"+key)
		rdb.ZRem(ctx, "tags", key)

		return nil

	}); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot delete the data", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(ctx, slog.LevelInfo, "tag deleted successfully")

	return nil
}

func tree(ctx context.Context) ([]Leaf, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "building tree")

	leaves := []Leaf{}

	roots, err := db.Redis.ZRange(ctx, "tags", 0, -1).Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the roots", slog.String("error", err.Error()))
		return []Leaf{}, errors.New("something went wrong")
	}

	tags, err := db.Redis.HGetAll(ctx, "tag").Result()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the tags", slog.String("error", err.Error()))
		return []Leaf{}, errors.New("something went wrong")
	}

	cmds, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for key := range tags {
			rdb.HGetAll(ctx, "tag:"+key)
		}

		return nil
	})

	if err != nil && err.Error() != "redis: nil" {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the tag meta data", slog.String("error", err.Error()))
		return []Leaf{}, errors.New("something went wrong")
	}

	tag := map[string]Tag{}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil && cmd.Err().Error() != "redis: nil" {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the tag meta data", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		val := cmd.(*redis.MapStringStringCmd).Val()

		tag[val["key"]] = Tag{}
	}

	for _, val := range roots {
		leaf := Leaf{
			Tag: Tag{
				Key:   val,
				Label: tag[val].Label,
				Image: tag[val].Image,
			},
			Branches: []*Leaf{},
		}

		branches := strings.Split(tags[val], ";")

		for _, branch := range branches {
			l := Leaf{
				Tag: Tag{
					Key:   branch,
					Label: tag[branch].Label,
					Image: tag[branch].Image,
				},
			}

			lb := strings.Split(tags[branch], ";")

			for _, b := range lb {
				l.Branches = append(l.Branches, &Leaf{
					Tag: Tag{
						Key:   b,
						Label: tag[branch].Label,
						Image: tag[branch].Image,
					},
				})
			}

			leaf.Branches = append(leaf.Branches, &l)
		}

		leaves = append(leaves, leaf)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "tree built")

	return leaves, nil
}
