package admin

import (
	"context"
	"errors"
	"gifthub/blogs"
	"gifthub/http/httpext"
	"log/slog"
	"mime/multipart"
	"strconv"
)

func processBlogFrom(ctx context.Context, form multipart.Form, id string) (blogs.Article, error) {
	exists := id != ""

	title := ""
	if len(form.Value["title"]) > 0 {
		title = form.Value["title"][0]
	}

	description := ""
	if len(form.Value["description"]) > 0 {
		description = form.Value["description"][0]
	}

	status := ""
	if len(form.Value["status"]) > 0 {
		status = form.Value["status"][0]
	}

	a := blogs.Article{
		Title:       title,
		Description: description,
		Status:      status,
	}

	if exists {
		id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.String("id", form.Value["id"][0]), slog.String("error", err.Error()))
			return blogs.Article{}, errors.New("input:id")
		}
		a.ID = int(id)
	} else {
		id, err := blogs.NextID(ctx)
		if err != nil {
			return blogs.Article{}, err
		}
		a.ID = id
	}

	err := a.Validate(ctx)
	if err != nil {
		return blogs.Article{}, err
	}

	files, err := httpext.ProcessFiles(ctx, form.File, []string{"image"})
	if err != nil {
		return blogs.Article{}, err
	}

	if files["image"] == nil && !exists {
		slog.LogAttrs(ctx, slog.LevelInfo, "the image is required")
		return blogs.Article{}, errors.New("input:image")
	}

	images, err := httpext.Upload(ctx, files)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot update the files", slog.String("error", err.Error()))
		return blogs.Article{}, errors.New("something went wrong")
	}

	if images["image"] != "" {
		a.Image = images["image"]
	}

	return a, nil

}
