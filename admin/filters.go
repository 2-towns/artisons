package admin

import (
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/http/pages"
	"artisons/products/filters"
	"artisons/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var filtersTpl *template.Template
var filtersHxTpl *template.Template
var filtersFormTpl *template.Template

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

func FilterSave(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var score int = 0
	if r.FormValue("score") != "" {
		val, err := strconv.ParseInt(r.FormValue("score"), 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the score", slog.String("score", r.FormValue("score")), slog.String("error", err.Error()))
			httperrors.HXCatch(w, ctx, "input:score")
			return
		}
		score = int(val)
	}

	id := chi.URLParam(r, "id")
	key := id
	if key == "" {
		key = r.FormValue("key")
	}

	f := filters.Filter{
		Key:    key,
		Label:  r.FormValue("label"),
		Score:  score,
		Active: r.FormValue("active") == "on",
		Values: r.Form["values"],
	}

	err := f.Validate(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if id == "" {
		exists, err := filters.Exists(ctx, key)
		if err != nil {
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}

		if exists {
			httperrors.HXCatch(w, ctx, "the filter exists already")
			return
		}
	}

	_, err = f.Save(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	Success(w, "/admin/filters.html")
}

func FilterList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := ctx.Value(contexts.Pagination).(pages.Paginator)

	res, err := filters.List(ctx, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	t := filtersTpl
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		t = filtersHxTpl
	}

	data := pages.Datalist(ctx, res.Filters)
	data.Pagination = p.Build(ctx, res.Total, len(res.Filters))
	data.Page = "Filters"

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func FilterForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var filter filters.Filter

	if id != "" {
		var err error
		filter, err = filters.Find(ctx, id)

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}
	}

	data := pages.Dataform[filters.Filter](ctx, filter)
	data.Page = "Filtres"

	if err := filtersFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func FilterDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	editable, err := filters.Editable(ctx, id)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if !editable {
		httperrors.HXCatch(w, ctx, "the filter cannot be editable")
		return
	}

	err = filters.Delete(ctx, id)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	pages.UpdateQuery(r)

	FilterList(w, r)
}
