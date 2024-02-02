package forms

import (
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/httperrors"
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
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Form struct {
}

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

func ParseForm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if err := r.ParseForm(); err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
			httperrors.HXCatch(w, ctx, "something went wrong")
			return
		}

		ctx = context.WithValue(ctx, contexts.Form, Form{})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ParseMultipartForm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := r.ParseMultipartForm(conf.MaxUploadSize); err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
			httperrors.HXCatch(w, ctx, "something went wrong")
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ParseOptionalID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := chi.URLParam(r, "id")
		if id != "" {
			val, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.Any("id", id), slog.String("error", err.Error()))
				httperrors.Page(w, ctx, "oops the data is not found", 404)
				return
			}

			ctx = context.WithValue(ctx, contexts.ID, val)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ParseID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")

		val, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}

		ctx = context.WithValue(ctx, contexts.ID, int(val))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
