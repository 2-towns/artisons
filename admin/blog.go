package admin

import (
	"artisons/blog"
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/http/forms"
	"artisons/http/httperrors"
	"artisons/http/pages"
	"artisons/string/stringutil"
	"artisons/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
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
		append(files, append(templates.AdminList,
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

func BlogSave(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	a := blog.Article{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Type:        "blog",
		Status:      "online",
	}

	if r.FormValue("status") != "on" {
		a.Status = "offline"
	}

	id, _ := ctx.Value(contexts.ID).(int)
	if id > 0 {
		a.ID = id
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

	query := blog.Query{Slug: a.Slug}
	res, err := blog.Search(ctx, query, 0, 1)
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

	Success(w, "/admin/blog.html")
}

func BlogList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := ctx.Value(contexts.Pagination).(pages.Paginator)

	qry := blog.Query{}
	if p.Query != "" {
		qry.Keywords = db.Escape(p.Query)
	}

	res, err := blog.Search(ctx, qry, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	t := blogTpl
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		t = blogHxTpl
	}

	data := pages.Datalist(ctx, res.Articles)
	data.Pagination = p.Build(ctx, res.Total, len(res.Articles))
	data.Page = "CMS"

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func BlogForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, _ := ctx.Value(contexts.ID).(int)

	var article blog.Article

	if id != 0 {
		var err error
		article, err = blog.Find(ctx, id)

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}
	}

	data := pages.Dataform[blog.Article](ctx, article)
	data.Page = "CMS"

	w.Header().Set("Content-Security-Policy", blogCspPolicy)

	if err := blogFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func BlogDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := ctx.Value(contexts.ID).(int)

	b, err := blog.Deletable(ctx, id)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if !b {
		httperrors.HXCatch(w, ctx, "you are not authorized to process this request")
		return
	}

	err = blog.Delete(ctx, id)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	pages.UpdateQuery(r)

	BlogList(w, r)
}
