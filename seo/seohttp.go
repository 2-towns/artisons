package seo

import (
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/http/httphelpers"
	"artisons/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"

	"golang.org/x/text/language"
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
		append(files, append(templates.AdminListHandler,
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

func AdminListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := httphelpers.BuildPaginator(r)

	res, err := List(ctx, p.Offset, p.Num)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the seo items", slog.String("error", err.Error()))
		httperrors.Page(w, ctx, "something went wrong", 500)
		return
	}

	t := seoTpl
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		t = seoHxTpl
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := httphelpers.List[Content]{
		Lang:       lang,
		Items:      res.Content,
		Empty:      len(res.Content) == 0,
		Currency:   conf.Currency,
		Pagination: p.Build(ctx, res.Total, len(res.Content)),
		Page:       "SEO",
		Flash:      httphelpers.Flash(w, r),
	}

	if err := t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AdminFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	var content Content

	if id != "" {
		var err error
		content, err = Find(ctx, id)

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := httphelpers.Form[Content]{
		Data:     content,
		Lang:     lang,
		Currency: conf.Currency,
		Page:     "SEO",
	}

	if err := seoFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AdminSaveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	c := Content{
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

	httphelpers.Success(w, "/admin/seo")
}
