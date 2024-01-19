package admin

import (
	"context"
	"errors"
	"gifthub/http/httpext"
	"gifthub/tags"
	"log/slog"
	"mime/multipart"
	"strconv"
	"strings"
)

func processTagFrom(ctx context.Context, form multipart.Form, id string) (tags.Tag, error) {
	exists := id != ""

	key := id
	if !exists && len(form.Value["key"]) > 0 {
		key = form.Value["key"][0]
	}

	label := ""
	if len(form.Value["label"]) > 0 {
		label = form.Value["label"][0]
	}

	children := []string{}
	if len(form.Value["children"]) > 0 {
		children = strings.Split(form.Value["children"][0], ";")
	}

	var score int = 0
	if len(form.Value["score"]) > 0 && form.Value["score"][0] != "" {
		val, err := strconv.ParseInt(form.Value["score"][0], 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the score", slog.String("score", form.Value["score"][0]), slog.String("error", err.Error()))
			return tags.Tag{}, errors.New("input:score")
		}
		score = int(val)
	}

	root := ""
	if len(form.Value["root"]) > 0 {
		root = form.Value["root"][0]
	}

	t := tags.Tag{
		Key:      key,
		Label:    label,
		Children: children,
		Root:     root == "on",
		Score:    score,
	}

	err := t.Validate(ctx)
	if err != nil {
		return tags.Tag{}, err
	}

	files, err := httpext.ProcessFiles(ctx, form.File, []string{"image"})
	if err != nil {
		return tags.Tag{}, err
	}

	images, err := httpext.Upload(ctx, files)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot update the files", slog.String("error", err.Error()))
		return tags.Tag{}, errors.New("something went wrong")
	}

	if images["image"] != "" {
		t.Image = images["image"]
	}

	return t, nil

}
