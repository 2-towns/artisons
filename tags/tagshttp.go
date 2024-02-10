package tags

import (
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/forms"
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

var tagsTpl *template.Template
var tagsHxTpl *template.Template
var tagsFormTpl *template.Template

func init() {
	var err error

	files := append(templates.AdminTable,
		conf.WorkingSpace+"web/views/admin/tags/tags-table.html",
	)

	tagsTpl, err = templates.Build("base.html").ParseFiles(
		append(files, append(templates.AdminListHandler,
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

func AdminSaveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(conf.MaxUploadSize); err != nil {
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

	t := Tag{
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
		exists, err := Exists(ctx, key)
		if err != nil {
			httperrors.HXCatch(w, ctx, "input:slug")
			return
		}

		if exists {
			httperrors.HXCatch(w, ctx, "the tag exists already")
			return
		}
	}

	eligible, err := AreEligible(ctx, t.Children)
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

	httphelpers.Success(w, "/admin/tags")
}

func AdminListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := httphelpers.BuildPaginator(r)

	res, err := List(ctx, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	t := tagsTpl
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		t = tagsHxTpl
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := httphelpers.List[Tag]{
		Lang:       lang,
		Items:      res.Tags,
		Empty:      len(res.Tags) == 0,
		Currency:   conf.Currency,
		Pagination: p.Build(ctx, res.Total, len(res.Tags)),
		Page:       "Tags",
		Flash:      httphelpers.Flash(w, r),
	}

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AdminFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	var tag Tag

	if id != "" {
		var err error
		tag, err = Find(ctx, id)

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := httphelpers.Form[Tag]{
		Data:     tag,
		Lang:     lang,
		Currency: conf.Currency,
		Page:     "Tags",
	}

	t, err := List(ctx, 0, 9999)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
	}

	data.Extra = t

	if err := tagsFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AdminDeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	err := Delete(ctx, id)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	p, _ := url.Parse(r.Header.Get("HX-Current-Url"))
	r.URL.Path = p.Path

	AdminListHandler(w, r)
}
