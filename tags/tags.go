package tags

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/contexts"
	"gifthub/validators"
	"log/slog"
	"slices"
	"strings"

	"github.com/redis/go-redis/v9"
	"golang.org/x/text/language"
)

type Tag struct {
	Name  string `validate:"required,alpha"`
	Label string `validate:"required"`
	// If the score is positive, the tag will be added in the root tag
	Score int
	Links []Tag
}

func (t Tag) Save(c context.Context) error {
	l := slog.With(slog.String("tag", t.Name))
	l.LogAttrs(c, slog.LevelInfo, "adding a new tag")

	if err := validators.V.Var(t.Name, "alpha"); err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate name", slog.String("name", t.Name), slog.String("error", err.Error()))
		return errors.New("input:tagname")
	}

	if err := validators.V.Var(t.Label, "required"); err != nil {
		l.LogAttrs(c, slog.LevelInfo, "cannot validate label", slog.String("label", t.Label), slog.String("error", err.Error()))
		return errors.New("input:taglabel")
	}

	if slices.Contains(conf.Languages, t.Name) {
		l.LogAttrs(c, slog.LevelInfo, "cannot use a reserved word")
		return errors.New("input_name_reserved")
	}

	ctx := context.Background()

	if _, err := db.Redis.TxPipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, "tag", t.Name, t.Label)

		if t.Score > 0 {
			rdb.ZAdd(ctx, "tag:root", redis.Z{
				Score:  float64(t.Score),
				Member: t.Name,
			})
		} else {
			rdb.ZRem(ctx, "tag:root", t.Name)
		}

		return nil

	}); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot store the data", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "tag saved successfully")

	return nil
}

// Link register a link with another tag with a score used by the Redis sorted set.
// An error is raised if the tag itself exists in the targeted tag links.
func (t Tag) Link(c context.Context, tag string, score float64) error {
	l := slog.With(slog.String("tag", t.Name))
	l.LogAttrs(c, slog.LevelInfo, "linking a tag", slog.String("target", tag))

	if tag == "" {
		slog.LogAttrs(c, slog.LevelInfo, "cannot continue with empty tag")
		return errors.New("input:name")
	}

	if score == 0 {
		slog.LogAttrs(c, slog.LevelInfo, "cannot continue with zero score")
		return errors.New("input:score")
	}

	ctx := context.Background()
	exists, err := db.Redis.HExists(ctx, "tag", tag).Result()
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot retrieve the target tag", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	if !exists {
		return errors.New("input_tag_notfound")
	}

	_, err = db.Redis.ZAdd(ctx, "tag:"+t.Name, redis.Z{
		Score:  score,
		Member: tag,
	}).Result()
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot link the tag", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "tag linked successfully")

	return nil
}

func List(c context.Context) ([]Tag, error) {
	slog.LogAttrs(c, slog.LevelInfo, "listing tags")

	hashes, err := db.Redis.HGetAll(context.Background(), "tag").Result()
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get the tags", slog.String("error", err.Error()))
		return []Tag{}, errors.New("something went wrong")
	}

	tags := []Tag{}

	for key, value := range hashes {
		tags = append(tags, Tag{
			Name:  key,
			Label: value,
		})
	}

	slog.LogAttrs(c, slog.LevelInfo, "found tags", slog.Int("length", len(tags)))

	return tags, nil
}

func (t Tag) WithLinks(c context.Context) (Tag, error) {
	slog.LogAttrs(c, slog.LevelInfo, "listing tags")

	links, err := db.Redis.ZRange(context.Background(), "tag:"+t.Name, 0, -1).Result()
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get the tags", slog.String("error", err.Error()))
		return t, errors.New("something went wrong")
	}

	hashes, err := db.Redis.HGetAll(context.Background(), "tag").Result()
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get the tags", slog.String("error", err.Error()))
		return t, errors.New("something went wrong")
	}

	tag := Tag{
		Name:  t.Name,
		Label: t.Label,
		Links: []Tag{},
	}

	for _, value := range links {
		tag.Links = append(tag.Links, Tag{
			Name:  value,
			Label: hashes[value],
		})
	}

	return tag, nil
}

func (t Tag) RemoveLink(c context.Context, tag string) error {
	l := slog.With(slog.String("tag", t.Name))
	l.LogAttrs(c, slog.LevelInfo, "removing the link", slog.String("target", tag))

	_, err := db.Redis.ZRem(context.Background(), "tag:"+t.Name, tag).Result()
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot link the tag", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	l.LogAttrs(c, slog.LevelInfo, "tag removed successfully")

	return nil
}

type treeDepth struct {
	Depth int
	Limit int
}

func tree(name string, labels map[string]string, links map[string][]string, td treeDepth) Tag {
	if td.Depth >= conf.TagMaxDepth || td.Depth >= td.Limit {
		return Tag{}
	}

	tag := Tag{
		Name:  name,
		Label: labels[name],
		Links: []Tag{},
	}

	for _, link := range links[name] {
		t := tree(link, labels, links, treeDepth{
			Depth: td.Depth + 1,
			Limit: td.Limit,
		})

		if t.Name == "" {
			continue
		}

		tag.Links = append(tag.Links, t)
	}

	return tag
}

// Root the tags available in the root tags.
// The limit is the depth limit used for the
func Root(c context.Context, limit int) ([]Tag, error) {
	slog.LogAttrs(c, slog.LevelInfo, "listing root  tags")

	ctx := context.Background()

	hashes, err := db.Redis.HGetAll(ctx, "tag").Result()
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get the tags", slog.String("error", err.Error()))
		return []Tag{}, errors.New("something went wrong")
	}

	labels := map[string]string{}

	cmds, err := db.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for key, value := range hashes {
			labels[key] = value
			rdb.ZRange(ctx, "tag:"+key, 0, -1)
		}

		return nil
	})

	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get the tag links", slog.String("error", err.Error()))
		return []Tag{}, errors.New("something went wrong")
	}

	links := map[string][]string{}

	for _, cmd := range cmds {
		key := fmt.Sprintf("%s", cmd.Args()[1])

		if cmd.Err() != nil {
			slog.LogAttrs(c, slog.LevelError, "cannot get the tag links", slog.String("key", key), slog.String("error", err.Error()))
			continue
		}

		name := strings.Replace(key, "tag:", "", 1)
		links[name] = cmd.(*redis.StringSliceCmd).Val()
	}

	tags := []Tag{}
	lang := c.Value(contexts.Locale).(language.Tag)

	for _, value := range links[lang.String()] {
		tags = append(tags, tree(value, labels, links, treeDepth{
			Depth: 0,
			Limit: limit,
		}))
	}

	slog.LogAttrs(c, slog.LevelInfo, "root tags loaded", slog.Int("length", len(tags)))

	return tags, nil
}
