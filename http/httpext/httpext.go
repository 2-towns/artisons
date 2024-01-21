package httpext

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/cookies"
	"gifthub/http/httperrors"
	"gifthub/templates"
	"html/template"
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
	"golang.org/x/text/language"
)

type form[T any] struct {
	Lang     language.Tag
	Page     string
	ID       interface{}
	Data     T
	Currency string
	Extra    interface{}
}

func Redirect(w http.ResponseWriter, r *http.Request, url string, status int) {
	isHX, _ := r.Context().Value(contexts.HX).(bool)

	if isHX {
		w.Header().Set("HX-Redirect", url)
	} else {
		http.Redirect(w, r, url, http.StatusFound)
	}
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

func rollbackUpload(ctx context.Context, images []string) {
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

type Entity interface {
	Save(ctx context.Context) (string, error)
	Validate(ctx context.Context) error
}

type ListFeature[T any] interface {
	Search(ctx context.Context, q string, offset, num int) (SearchResults[T], error)
	ListTemplate(ctx context.Context) *template.Template
}

type FormFeature[T any] interface {
	ID(ctx context.Context, id string) (interface{}, error)
	Find(ctx context.Context, id interface{}) (T, error)
}

type List[T any] struct {
	Name    string
	URL     string
	Feature ListFeature[T]
}

type Delete[T Entity] struct {
	List[T]
	Feature DeleteFeature[T]
}

type Form[T any] struct {
	Name    string
	Feature FormFeature[T]
}

type DeleteFeature[T Entity] interface {
	ID(ctx context.Context, id string) (interface{}, error)
	Delete(ctx context.Context, id interface{}) error
}

type DigestFeature[T Entity] interface {
	IsImageRequired(e T, key string) bool
	Digest(ctx context.Context, r *http.Request) (T, error)
	UpdateImage(e *T, key, image string)
}

type Save[T Entity] struct {
	Name       string
	URL        string
	Images     []string
	Folder     string
	Form       FormType[T]
	Feature    DigestFeature[T]
	NoRedirect bool
}

type FormType[T Entity] interface {
	Parse(ctx context.Context, r *http.Request) error
}

type MultipartForm struct {
}

type UrlEncodedForm struct {
}

func (f MultipartForm) Parse(ctx context.Context, r *http.Request) error {
	if err := r.ParseMultipartForm(conf.MaxUploadSize); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	return nil
}

func (f UrlEncodedForm) Parse(ctx context.Context, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	return nil
}

type SearchResults[T any] struct {
	Total int
	Items []T
}

func DigestList[T any](w http.ResponseWriter, r *http.Request, l List[T]) {
	var page int = 1

	ppage := r.URL.Query().Get("page")
	if ppage != "" {
		if d, err := strconv.ParseInt(ppage, 10, 32); err == nil && d > 0 {
			page = int(d)
		}
	}

	query := r.URL.Query().Get("q")
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	offset := (page - 1) * conf.ItemsPerPage
	num := offset + conf.ItemsPerPage

	res, err := l.Feature.Search(ctx, query, offset, num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	pag := templates.Paginate(page, len(res.Items), int(res.Total))
	pag.URL = l.URL
	pag.Lang = lang

	flash := ""
	c, err := r.Cookie(cookies.FlashMessage)
	if err == nil && c != nil {
		flash = c.Value

		cookie := &http.Cookie{
			Name:     cookies.FlashMessage,
			Value:    flash,
			MaxAge:   -1,
			Path:     "/",
			HttpOnly: true,
			Secure:   conf.Cookie.Secure,
			Domain:   conf.Cookie.Domain,
		}

		http.SetCookie(w, cookie)
	}

	d := struct {
		Lang       language.Tag
		Page       string
		Items      []T
		Empty      bool
		Currency   string
		Pagination templates.Pagination
		Flash      string
	}{
		lang,
		l.Name,
		res.Items,
		len(res.Items) == 0,
		conf.Currency,
		pag,
		flash,
	}

	tpl := l.Feature.ListTemplate(ctx)

	if err = tpl.Execute(w, &d); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func DigestSave[T Entity](w http.ResponseWriter, r *http.Request, f Save[T]) {
	ctx := r.Context()

	if err := f.Form.Parse(ctx, r); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	data, err := f.Feature.Digest(ctx, r)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	err = data.Validate(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	filepaths := []string{}

	for _, filename := range f.Images {
		file, header, err := r.FormFile(filename)

		if err != nil {
			if err == http.ErrMissingFile && f.Feature.IsImageRequired(data, filename) {
				slog.LogAttrs(ctx, slog.LevelError, "cannot process the image", slog.String("error", err.Error()), slog.String("image", filename))
				httperrors.HXCatch(w, ctx, fmt.Sprintf("input:%s", filename))
				return
			}

			continue
		}

		defer file.Close()

		filepath, err := upload(ctx, file, header, filename, f.Folder)
		if err != nil {
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}

		filepaths = append(filepaths, filepath)
		f.Feature.UpdateImage(&data, filename, filepath)
	}

	_, err = data.Save(ctx)
	if err != nil {
		rollbackUpload(ctx, filepaths)
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if f.NoRedirect {
		return
	}

	cookie := &http.Cookie{
		Name:     cookies.FlashMessage,
		Value:    "The data has been saved successfully.",
		MaxAge:   int(time.Minute.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   conf.Cookie.Secure,
		Domain:   conf.Cookie.Domain,
	}

	http.SetCookie(w, cookie)
	w.Header().Set("HX-Redirect", f.URL)
	w.Write([]byte(""))
}

func DigestForm[T any](w http.ResponseWriter, r *http.Request, f Form[T]) form[T] {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	id := chi.URLParam(r, "id")

	var item T

	fid, err := f.Feature.ID(ctx, id)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.Any("id", id), slog.String("error", err.Error()))
		httperrors.Page(w, ctx, "oops the data is not found", 404)
		return form[T]{}
	}

	if fid != "" && fid != 0 {
		item, err = f.Feature.Find(ctx, fid)
		if err != nil {
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return form[T]{}
		}
	}

	data := form[T]{
		lang,
		f.Name,
		id,
		item,
		conf.Currency,
		"",
	}

	return data
}

func DigestDelete[T Entity](w http.ResponseWriter, r *http.Request, f Delete[T]) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	id := chi.URLParam(r, "id")

	fid, err := f.Feature.ID(ctx, id)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.Any("id", id), slog.String("error", err.Error()))
		httperrors.Page(w, ctx, "oops the data is not found", 404)
	}

	err = f.Feature.Delete(ctx, fid)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	page := r.FormValue("page")
	r.URL.RawQuery = fmt.Sprintf("page%s", page)

	query := r.FormValue("query")
	if query != "" {
		r.URL.RawQuery += fmt.Sprintf("&query=%s", query)
	}

	DigestList(w, r, f.List)
}
