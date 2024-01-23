package admin

import (
	"context"
	"errors"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/httperrors"
	"gifthub/http/httpext"
	"gifthub/tags"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const tagsName = "Tags"
const tagsURL = "/admin/tags.html"
const tagsFolder = "tags"

var tagsTpl *template.Template
var tagsHxTpl *template.Template
var tagsFormTpl *template.Template

type tagsFeature struct{}

func init() {
	var err error

	files := append(templates.AdminTable,
		conf.WorkingSpace+"web/views/admin/tags/tags-table.html",
	)

	tagsTpl, err = templates.Build("base.html").ParseFiles(
		append(files, append(templates.AdminList,
			conf.WorkingSpace+"web/views/admin/tags/tags-actions.html",
			conf.WorkingSpace+"web/views/admin/tags/tags.html")...,
		)...)

	if err != nil {
		log.Panicln(err)
	}

	tagsHxTpl, err = templates.Build("tags-table.html").ParseFiles(files...)

	if err != nil {
		log.Panicln(err)
	}

	tagsFormTpl, err = templates.Build("base.html").ParseFiles(
		append(templates.AdminUI,
			conf.WorkingSpace+"web/views/admin/tags/tags-scripts.html",
			conf.WorkingSpace+"web/views/admin/tags/tags-head.html",
			conf.WorkingSpace+"web/views/admin/tags/tags-form.html",
		)...)

	if err != nil {
		log.Panicln(err)
	}
}

func (f tagsFeature) ListTemplate(ctx context.Context) *template.Template {
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		return tagsHxTpl
	}

	return tagsTpl
}

func (f tagsFeature) Search(ctx context.Context, q string, offset, num int) (httpext.SearchResults[tags.Tag], error) {
	res, err := tags.List(ctx, offset, num)

	return httpext.SearchResults[tags.Tag]{
		Total: res.Total,
		Items: res.Tags,
	}, err
}

func (data tagsFeature) Digest(ctx context.Context, r *http.Request) (tags.Tag, error) {
	key := chi.URLParam(r, "id")
	if key == "" {
		key = r.FormValue("key")

		exists, err := tags.Exists(ctx, key)
		if err != nil {
			return tags.Tag{}, err
		}

		if exists {
			return tags.Tag{}, errors.New("the tag exists already")
		}
	}

	var score int = 0
	if r.FormValue("score") != "" {
		val, err := strconv.ParseInt(r.FormValue("score"), 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the score", slog.String("score", r.FormValue("score")), slog.String("error", err.Error()))
			return tags.Tag{}, errors.New("input:score")
		}
		score = int(val)
	}

	t := tags.Tag{
		Key:      key,
		Label:    r.FormValue("label"),
		Children: r.MultipartForm.Value["children"],
		Root:     r.FormValue("root") == "on",
		Score:    score,
	}

	if err := tags.AreEligible(ctx, t.Children); err != nil {
		return tags.Tag{}, err
	}

	return t, nil
}

func (f tagsFeature) ID(ctx context.Context, id string) (interface{}, error) {
	return id, nil
}

func (f tagsFeature) Find(ctx context.Context, id interface{}) (tags.Tag, error) {
	return tags.Find(ctx, id.(string))
}

func (f tagsFeature) Delete(ctx context.Context, id interface{}) error {
	return tags.Delete(ctx, id.(string))
}

func (f tagsFeature) IsImageRequired(a tags.Tag, key string) bool {
	return false
}

func (f tagsFeature) UpdateImage(a *tags.Tag, key, image string) {
}

func TagsSave(w http.ResponseWriter, r *http.Request) {
	httpext.DigestSave[tags.Tag](w, r, httpext.Save[tags.Tag]{
		Name:    tagsName,
		URL:     tagsURL,
		Feature: tagsFeature{},
		Form:    httpext.MultipartForm{},
		Images:  []string{"image"},
		Folder:  tagsFolder,
	})
}

func TagsList(w http.ResponseWriter, r *http.Request) {
	httpext.DigestList[tags.Tag](w, r, httpext.List[tags.Tag]{
		Name:    tagsName,
		URL:     tagsURL,
		Feature: tagsFeature{},
	})
}

func TagsForm(w http.ResponseWriter, r *http.Request) {
	data := httpext.DigestForm[tags.Tag](w, r, httpext.Form[tags.Tag]{
		Name:    tagsName,
		Feature: tagsFeature{},
	})

	if data.Page == "" {
		return
	}

	ctx := r.Context()
	t, err := tags.List(ctx, 0, 9999)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
	}

	data.Extra = t

	if err := tagsFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func TagsDelete(w http.ResponseWriter, r *http.Request) {
	httpext.DigestDelete[tags.Tag](w, r, httpext.Delete[tags.Tag]{
		List: httpext.List[tags.Tag]{
			Name:    tagsName,
			URL:     tagsURL,
			Feature: tagsFeature{},
		},
		Feature: tagsFeature{},
	})
}
