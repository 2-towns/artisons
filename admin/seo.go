package admin

import (
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/cookies"
	"gifthub/http/seo"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
)

var seoTpl *template.Template
var seoHxTpl *template.Template

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
}

func Seo(w http.ResponseWriter, r *http.Request) {
	var page int = 1

	ppage := r.URL.Query().Get("page")
	if ppage != "" {
		if d, err := strconv.ParseInt(ppage, 10, 32); err == nil && d > 0 {
			page = int(d)
		}
	}

	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	offset := (page - 1) * conf.ItemsPerPage
	num := offset + conf.ItemsPerPage

	res := seo.List(ctx, offset, num)
	pag := templates.Paginate(page, len(res.Content), int(res.Total))
	pag.URL = "/admin/seo.html"
	pag.Lang = lang

	flash := ""
	c, err := r.Cookie(cookies.FlashMessage)
	if err == nil && c != nil {
		flash = c.Value
	}

	data := struct {
		Lang       language.Tag
		Page       string
		Content    []seo.Content
		Empty      bool
		Pagination templates.Pagination
		Flash      string
	}{
		lang,
		"SEO",
		res.Content,
		len(res.Content) == 0,
		pag,
		flash,
	}

	isHX, _ := ctx.Value(contexts.HX).(bool)
	var t *template.Template = seoTpl
	if isHX {
		t = seoHxTpl
	}

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
