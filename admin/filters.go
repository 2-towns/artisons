package admin

import (
	"context"
	"errors"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/httpext"
	"gifthub/products/filters"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const filtersName = "Filters"
const filtersURL = "/admin/filters.html"

var filtersTpl *template.Template
var filtersHxTpl *template.Template
var filtersFormTpl *template.Template

type filtersFeature struct{}

func init() {
	var err error

	files := append(templates.AdminTable,
		conf.WorkingSpace+"web/views/admin/filters/filters-table.html",
	)

	filtersTpl, err = templates.Build("base.html").ParseFiles(
		append(files, append(templates.AdminList,
			conf.WorkingSpace+"web/views/admin/filters/filters-actions.html",
			conf.WorkingSpace+"web/views/admin/filters/filters.html")...,
		)...)

	if err != nil {
		log.Panicln(err)
	}

	filtersHxTpl, err = templates.Build("filters-table.html").ParseFiles(files...)

	if err != nil {
		log.Panicln(err)
	}

	filtersFormTpl, err = templates.Build("base.html").ParseFiles(
		append(templates.AdminUI,
			conf.WorkingSpace+"web/views/admin/filters/filters-scripts.html",
			conf.WorkingSpace+"web/views/admin/filters/filters-head.html",
			conf.WorkingSpace+"web/views/admin/filters/filters-form.html",
		)...)

	if err != nil {
		log.Panicln(err)
	}
}

func (f filtersFeature) ListTemplate(ctx context.Context) *template.Template {
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		return filtersHxTpl
	}

	return filtersTpl
}

func (f filtersFeature) Search(ctx context.Context, q string, offset, num int) (httpext.SearchResults[filters.Filter], error) {
	res, err := filters.List(ctx, offset, num)

	return httpext.SearchResults[filters.Filter]{
		Total: res.Total,
		Items: res.Filters,
	}, err
}

func (data filtersFeature) Digest(ctx context.Context, r *http.Request) (filters.Filter, error) {
	key := chi.URLParam(r, "id")
	if key == "" {
		key = r.FormValue("key")

		exists, err := filters.Exists(ctx, key)
		if err != nil {
			return filters.Filter{}, err
		}

		if exists {
			return filters.Filter{}, errors.New("the filter exists already")
		}
	}

	var score int = 0
	if r.FormValue("score") != "" {
		val, err := strconv.ParseInt(r.FormValue("score"), 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the score", slog.String("score", r.FormValue("score")), slog.String("error", err.Error()))
			return filters.Filter{}, errors.New("input:score")
		}
		score = int(val)
	}

	typ := "list"
	if key == "colors" {
		typ = "color"
	}

	f := filters.Filter{
		Key:    key,
		Label:  r.FormValue("label"),
		Score:  score,
		Active: r.FormValue("active") == "on",
		Type:   typ,
		Values: r.Form["values"],
	}

	log.Println("label!!!!!!!!!!!!!!", r.Form)

	return f, nil
}

func (f filtersFeature) ID(ctx context.Context, id string) (interface{}, error) {
	return id, nil
}

func (f filtersFeature) Find(ctx context.Context, id interface{}) (filters.Filter, error) {
	return filters.Find(ctx, id.(string))
}

func (f filtersFeature) Delete(ctx context.Context, id interface{}) error {
	return filters.Delete(ctx, id.(string))
}

func (f filtersFeature) IsImageRequired(a filters.Filter, key string) bool {
	return false
}

func (f filtersFeature) UpdateImage(a *filters.Filter, key, image string) {
}

func FiltersSave(w http.ResponseWriter, r *http.Request) {
	httpext.DigestSave[filters.Filter](w, r, httpext.Save[filters.Filter]{
		Name:    filtersName,
		URL:     filtersURL,
		Feature: filtersFeature{},
		Form:    httpext.UrlEncodedForm{},
		Images:  []string{},
		Folder:  "",
	})
}

func FiltersList(w http.ResponseWriter, r *http.Request) {
	httpext.DigestList[filters.Filter](w, r, httpext.List[filters.Filter]{
		Name:    filtersName,
		URL:     filtersURL,
		Feature: filtersFeature{},
	})
}

func FiltersForm(w http.ResponseWriter, r *http.Request) {
	data := httpext.DigestForm[filters.Filter](w, r, httpext.Form[filters.Filter]{
		Name:    filtersName,
		Feature: filtersFeature{},
	})

	if data.Page == "" {
		return
	}

	if err := filtersFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func FiltersDelete(w http.ResponseWriter, r *http.Request) {
	httpext.DigestDelete[filters.Filter](w, r, httpext.Delete[filters.Filter]{
		List: httpext.List[filters.Filter]{
			Name:    filtersName,
			URL:     filtersURL,
			Feature: filtersFeature{},
		},
		Feature: filtersFeature{},
	})
}
