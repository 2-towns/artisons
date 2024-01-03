package httpext

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/http/contexts"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slices"
	"time"

	"golang.org/x/exp/maps"
)

func Redirect(w http.ResponseWriter, r *http.Request, url string, status int) {
	isHX, _ := r.Context().Value(contexts.HX).(bool)

	if isHX {
		w.Header().Set("HX-Redirect", url)
	} else {
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func ProcessFiles(ctx context.Context, files map[string][]*multipart.FileHeader, names []string) (map[string]*multipart.FileHeader, error) {
	headers := map[string]*multipart.FileHeader{}
	for _, name := range names {
		if files[name] == nil {
			continue
		}

		img := files[name][0]
		slog.LogAttrs(ctx, slog.LevelInfo, "image info", slog.String("image", img.Filename), slog.Int64("size", img.Size), slog.Any("headers", img.Header))

		ct := img.Header["Content-Type"][0]
		if !slices.Contains(conf.ImagesAllowed, ct) {
			slog.LogAttrs(ctx, slog.LevelError, "cannot use an image in a unknown extension", slog.String("contentType", ct))
			return map[string]*multipart.FileHeader{}, fmt.Errorf("input_%s_invalid", name)
		}

		headers[name] = img
	}

	return headers, nil
}

func Upload(ctx context.Context, headers map[string]*multipart.FileHeader) (map[string]string, error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "uploading images")

	fns := map[string]string{}

	for key, h := range headers {

		file, err := h.Open()
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot open the image")
			return map[string]string{}, errors.New("input_images_invalid")
		}

		defer file.Close()

		ext := filepath.Ext(h.Filename)
		fn := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

		dst, err := os.Create(path.Join(conf.ImgProxy.Path, "products", fn))
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot create the file", slog.String("error", err.Error()), slog.String("filename", h.Filename))

			RollbackUpload(ctx, maps.Values(fns))

			return map[string]string{}, err
		}

		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot copy the file", slog.String("error", err.Error()))

			RollbackUpload(ctx, maps.Values(fns))

			return map[string]string{}, err
		}

		fns[key] = fn
	}

	return fns, nil
}

func RollbackUpload(ctx context.Context, images []string) {
	for _, value := range images {
		if value == "" {
			continue
		}

		err := os.Remove(path.Join(conf.ImgProxy.Path, "products", value))
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot remove the image", slog.String("error", err.Error()), slog.String("image", value))
			continue
		}
	}
}
