package admin

import (
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/http/pages"
	"artisons/http/seo"
	"artisons/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var seoTpl *template.Template
var seoHxTpl *template.Template
var seoFormTpl *template.Template

func init() {
	var err error

	files := append(templates.AdminTable,
		conf.WorkingSpace+"web/views/admin/seo/seo-table.html",
	)

	seoTpl, err = templates.Build("base.html").ParseFiles(
		append(files, append(templates.AdminList,
			conf.WorkingSpace+"web/views/admin/seo/seo.html",
		)...)...)

	if err != nil {
		log.Panicln(err)
	}

	seoHxTpl, err = templates.Build("seo-table.html").ParseFiles(files...)

	if err != nil {
		log.Panicln(err)
	}

	seoFormTpl, err = templates.Build("base.html").ParseFiles(
		append(templates.AdminUI,
			conf.WorkingSpace+"web/views/admin/seo/seo-form.html",
		)...)

	if err != nil {
		log.Panicln(err)
	}
}

func SeoList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := ctx.Value(contexts.Pagination).(pages.Paginator)

	res := seo.List(ctx, p.Offset, p.Num)

	t := seoTpl
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		t = seoHxTpl
	}

	data := pages.Datalist(ctx, res.Content)
	data.Pagination = p.Build(ctx, res.Total, len(res.Content))
	data.Page = "SEO"

	if err := t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func SeoForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var content seo.Content

	if id != "" {
		var err error
		content, err = seo.Find(ctx, id)

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}
	}

	data := pages.Dataform[seo.Content](ctx, content)
	data.Page = "SEO"

	if err := seoFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func SeoSave(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := chi.URLParam(r, "id")

	c := seo.Content{
		Key:         key,
		URL:         r.FormValue("url"),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	err := c.Validate(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	_, err = c.Save(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	Success(w, "/admin/seo.html")
}
