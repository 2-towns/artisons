package filters

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
	"net/url"
	"strconv"

	"golang.org/x/text/language"
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
		append(files, append(templates.AdminListHandler,
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

func AdminSaveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

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

	id := r.PathValue("id")
	key := id
	if key == "" {
		key = r.FormValue("key")
	}

	f := Filter{
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
		exists, err := Exists(ctx, key)
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

	httphelpers.Success(w, "/admin/filters")
}

func AdminListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := httphelpers.BuildPaginator(r)

	res, err := List(ctx, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	t := filtersTpl
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		t = filtersHxTpl
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := httphelpers.List[Filter]{
		Lang:       lang,
		Items:      res.Filters,
		Empty:      len(res.Filters) == 0,
		Currency:   conf.Currency,
		Pagination: p.Build(ctx, res.Total, len(res.Filters)),
		Page:       "Filters",
		Flash:      httphelpers.Flash(w, r),
	}

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AdminFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	var filter Filter

	if id != "" {
		var err error
		filter, err = Find(ctx, id)

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := httphelpers.Form[Filter]{
		Data:     filter,
		Lang:     lang,
		Currency: conf.Currency,
		Page:     "Filters",
	}

	if err := filtersFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AdminDeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	editable, err := Editable(ctx, id)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if !editable {
		httperrors.HXCatch(w, ctx, "the filter cannot be editable")
		return
	}

	err = Delete(ctx, id)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	p, _ := url.Parse(r.Header.Get("HX-Current-Url"))
	r.URL.Path = p.Path

	AdminListHandler(w, r)
}
