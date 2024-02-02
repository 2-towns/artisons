package admin

import (
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/forms"
	"artisons/http/httperrors"
	"artisons/http/pages"
	"artisons/tags"
	"artisons/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var tagsTpl *template.Template
var tagsHxTpl *template.Template
var tagsFormTpl *template.Template

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

func TagSave(w http.ResponseWriter, r *http.Request) {
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

	t := tags.Tag{
		Key:      key,
		Label:    r.FormValue("label"),
		Children: r.MultipartForm.Value["children"],
		Root:     r.FormValue("root") == "on",
		Score:    score,
	}

	err := t.Validate(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if id == "" {
		exists, err := tags.Exists(ctx, key)
		if err != nil {
			httperrors.HXCatch(w, ctx, "input:slug")
			return
		}

		if exists {
			httperrors.HXCatch(w, ctx, "the tag exists already")
			return
		}
	}

	eligible, err := tags.AreEligible(ctx, t.Children)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if !eligible {
		httperrors.HXCatch(w, ctx, "the children cannot be root tags")
		return
	}

	images := []string{"image"}
	files, err := forms.Upload(r, "tags", images)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	t.Image = files[0]

	_, err = t.Save(ctx)
	if err != nil {
		forms.RollbackUpload(ctx, files)
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	Success(w, "/admin/tags.html")
}

func TagList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := ctx.Value(contexts.Pagination).(pages.Paginator)

	res, err := tags.List(ctx, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	t := tagsTpl
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		t = tagsHxTpl
	}

	data := pages.Datalist(ctx, res.Tags)
	data.Pagination = p.Build(ctx, res.Total, len(res.Tags))
	data.Page = "Tags"

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func TagForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var tag tags.Tag

	if id != "" {
		var err error
		tag, err = tags.Find(ctx, id)

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}
	}

	data := pages.Dataform[tags.Tag](ctx, tag)
	data.Page = "Tags"

	t, err := tags.List(ctx, 0, 9999)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
	}

	data.Extra = t

	if err := tagsFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func TagDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	err := tags.Delete(ctx, id)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	pages.UpdateQuery(r)

	TagList(w, r)
}
