package forms

import (
	"artisons/conf"
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slices"
	"time"
)

func upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, filename, folder string) (string, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "uploading image", slog.String("image", header.Filename), slog.Int64("size", header.Size), slog.Any("headers", header.Header))

	ct := header.Header["Content-Type"][0]
	if !slices.Contains(conf.ImagesAllowed, ct) {
		slog.LogAttrs(ctx, slog.LevelError, "cannot use an image in a unknown extension", slog.String("contentType", ct))
		return "", fmt.Errorf("input:%s", filename)
	}

	ext := filepath.Ext(header.Filename)
	filepath := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst, err := os.Create(path.Join(conf.ImgProxy.Path, folder, filepath))

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot create the file", slog.String("error", err.Error()), slog.String("filename", header.Filename))
		return "", fmt.Errorf("input:%s", filename)
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot copy the file", slog.String("error", err.Error()))
		return "", fmt.Errorf("input:%s", filename)
	}

	return filepath, nil
}

func Upload(r *http.Request, folder string, images []string) ([]string, error) {
	ctx := r.Context()
	filepaths := []string{}

	for _, filename := range images {
		file, header, err := r.FormFile(filename)

		if err == http.ErrMissingFile {
			filepaths = append(filepaths, "")
			continue
		}

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the form file", slog.String("filename", filename))
			return []string{}, err
		}

		defer file.Close()

		filepath, err := upload(ctx, file, header, filename, folder)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot upload the form file", slog.String("filename", filename))
			return []string{}, err
		}

		filepaths = append(filepaths, filepath)
	}

	return filepaths, nil
}

func RollbackUpload(ctx context.Context, images []string) {
	for _, value := range images {
		if value == "" {
			continue
		}

		err := os.Remove(path.Join(conf.ImgProxy.Path, value))
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot remove the image", slog.String("error", err.Error()), slog.String("image", value))
			continue
		}
	}
}
