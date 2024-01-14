package admin

import (
	"gifthub/blogs"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/cookies"
	"gifthub/http/httperrors"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
)

var blogTpl *template.Template
var blogHxTpl *template.Template

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
}

func Blog(w http.ResponseWriter, r *http.Request) {
	var page int = 1

	ppage := r.URL.Query().Get("page")
	if ppage != "" {
		if d, err := strconv.ParseInt(ppage, 10, 32); err == nil && d > 0 {
			page = int(d)
		}
	}

	q := r.URL.Query().Get("q")
	query := blogs.Query{}
	if q != "" {
		query.Keywords = q
	}

	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	offset := (page - 1) * conf.ItemsPerPage
	num := offset + conf.ItemsPerPage

	res, err := blogs.Search(ctx, query, offset, num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	pag := templates.Paginate(page, len(res.Articles), int(res.Total))
	pag.URL = "/admin/blog.html"
	pag.Lang = lang

	flash := ""
	c, err := r.Cookie(cookies.FlashMessage)
	if err == nil && c != nil {
		flash = c.Value
	}

	data := struct {
		Lang       language.Tag
		Page       string
		Articles   []blogs.Article
		Empty      bool
		Pagination templates.Pagination
		Flash      string
	}{
		lang,
		"Blog",
		res.Articles,
		len(res.Articles) == 0,
		pag,
		flash,
	}

	isHX, _ := ctx.Value(contexts.HX).(bool)
	var t *template.Template = blogTpl
	if isHX {
		t = blogHxTpl
	}

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
