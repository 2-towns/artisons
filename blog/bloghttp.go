// Package pages provides the application pages
package blog

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/http/forms"
	"artisons/http/httperrors"
	"artisons/http/httphelpers"
	"artisons/shops"
	"artisons/string/stringutil"
	"artisons/tags/tree"
	"artisons/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/text/language"
)

var blogTpl *template.Template
var blogHxTpl *template.Template
var blogFormTpl *template.Template
var blogCspPolicy = ""

func init() {
	var err error

	files := append(templates.AdminTable,
		conf.WorkingSpace+"web/views/admin/blog/blog-table.html",
	)

	blogTpl, err = templates.Build("base.html").ParseFiles(
		append(files, append(templates.AdminListHandler,
			conf.WorkingSpace+"web/views/admin/blog/blog-actions.html",
			conf.WorkingSpace+"web/views/admin/blog/blog.html",
		)...)...)

	if err != nil {
		log.Panicln(err)
	}

	blogHxTpl, err = templates.Build("blog-table.html").ParseFiles(files...)

	if err != nil {
		log.Panicln(err)
	}

	blogFormTpl, err = templates.Build("base.html").ParseFiles(
		append(templates.AdminUI,
			conf.WorkingSpace+"web/views/admin/blog/blog-head.html",
			conf.WorkingSpace+"web/views/admin/blog/blog-scripts.html",
			conf.WorkingSpace+"web/views/admin/blog/blog-form.html",
		)...)

	if err != nil {
		log.Panicln(err)
	}

	blogCspPolicy = "default-src 'self'"
	blogCspPolicy += " https://maxcdn.bootstrapcdn.com/font-awesome/latest/css/font-awesome.min.css"
	blogCspPolicy += " https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.eot"
	blogCspPolicy += " https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.woff2"
	blogCspPolicy += " https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.ttf"
	blogCspPolicy += " https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.svg"
	blogCspPolicy += " https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.woff"
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	p := httphelpers.BuildPaginator(r)

	query := Query{
		Keywords: p.Query,
		Type:     "blog",
	}

	res, err := Search(ctx, query, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	pag := p.Build(ctx, res.Total, len(res.Articles))

	data := struct {
		Lang       language.Tag
		Shop       shops.Settings
		Articles   []Article
		Empty      bool
		Pagination httphelpers.Pagination
		Tags       []tree.Leaf
	}{
		lang,
		shops.Data,
		res.Articles,
		len(res.Articles) == 0,
		pag,
		tree.Tree,
	}

	var t *template.Template
	isHX, _ := ctx.Value(contexts.HX).(bool)

	if isHX {
		t = templates.Pages["hx-blog"]
	} else {
		t = templates.Pages["blog"]
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	slug := r.PathValue("slug")

	query := Query{Slug: slug}
	res, err := Search(ctx, query, 0, 1)

	log.Println("slug", slug, res.Articles)

	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the article", slog.String("slug", slug), slog.String("error", err.Error()))
		httperrors.Page(w, r.Context(), err.Error(), 500)
		return
	}

	if res.Total == 0 {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot find the article", slog.String("slug", slug))
		httperrors.Page(w, r.Context(), "oops the data is not found", 404)
		return
	}

	a := res.Articles[0]

	data := struct {
		Lang    language.Tag
		Shop    shops.Settings
		Article Article
		Tags    []tree.Leaf
	}{
		lang,
		shops.Data,
		a,
		tree.Tree,
	}

	if err := templates.Pages["static"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AdminSaveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(conf.MaxUploadSize); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	a := Article{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Type:        "blog",
		Status:      "online",
	}

	id := r.PathValue("id")
	if id != "" {
		val, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}

		a.ID = int(val)
	}

	if r.FormValue("status") != "on" {
		a.Status = "offline"
	}

	if r.FormValue("slug") != "" {
		a.Slug = r.FormValue("slug")
	} else {
		a.Slug = stringutil.Slugify(a.Title)
	}

	err := a.Validate(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	query := Query{Slug: a.Slug}
	res, err := Search(ctx, query, 0, 1)
	if err != nil || (res.Total > 0 && (res.Articles[0].ID != a.ID)) {
		httperrors.HXCatch(w, ctx, "input:slug")
		return
	}

	images := []string{"image"}
	files, err := forms.Upload(r, "blog", images)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	a.Image = files[0]

	_, err = a.Save(ctx)
	if err != nil {
		forms.RollbackUpload(ctx, files)
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	httphelpers.Success(w, "/admin/blog")
}

func AdminListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := httphelpers.BuildPaginator(r)

	qry := Query{}
	if p.Query != "" {
		qry.Keywords = db.Escape(p.Query)
	}

	res, err := Search(ctx, qry, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	t := blogTpl
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		t = blogHxTpl
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := httphelpers.List[Article]{
		Lang:       lang,
		Items:      res.Articles,
		Empty:      len(res.Articles) == 0,
		Currency:   conf.Currency,
		Pagination: p.Build(ctx, res.Total, len(res.Articles)),
		Page:       "CMS",
		Flash:      httphelpers.Flash(w, r),
	}

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AdminFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	var article Article

	if id != "" {
		var err error
		bid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}

		article, err = Find(ctx, int(bid))

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot find the blog", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := httphelpers.Form[Article]{
		Data:     article,
		Lang:     lang,
		Currency: conf.Currency,
		Page:     "CMS",
	}

	w.Header().Set("Content-Security-Policy", blogCspPolicy)

	if err := blogFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AdminDeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	bid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.Any("id", id), slog.String("error", err.Error()))
		httperrors.Page(w, ctx, "oops the data is not found", 404)
		return
	}

	b, err := Deletable(ctx, int(bid))
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if !b {
		httperrors.HXCatch(w, ctx, "you are not authorized to process this request")
		return
	}

	err = Delete(ctx, int(bid))
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	p, _ := url.Parse(r.Header.Get("HX-Current-Url"))
	r.URL.Path = p.Path

	AdminListHandler(w, r)
}
