// Package pages provides the application pages
package pages

import (
	"artisons/blog"
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/shops"
	"artisons/tags"
	"artisons/templates"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
)

func Blog(w http.ResponseWriter, r *http.Request) {
	var page int = 1
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	uquery := r.URL.Query()

	p := uquery.Get("page")
	if p != "" {
		if d, err := strconv.ParseInt(p, 10, 32); err == nil && d > 0 {
			page = int(d)
		}
	}

	q := uquery.Get("q")
	offset := 0
	if page > 0 {
		offset = (page - 1) * conf.ItemsPerPage
	}
	num := offset + conf.ItemsPerPage
	query := blog.Query{
		Keywords: q,
		Type:     "blog",
	}

	res, err := blog.Search(ctx, query, offset, num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	pag := templates.Paginate(page, len(res.Articles), int(res.Total))
	pag.URL = "/blog.html"
	pag.Lang = lang

	data := struct {
		Lang       language.Tag
		Shop       shops.Settings
		Articles   []blog.Article
		Empty      bool
		Pagination templates.Pagination
		Tags       []tags.Leaf
	}{
		lang,
		shops.Data,
		res.Articles,
		len(res.Articles) == 0,
		pag,
		tags.Tree,
	}

	var t *template.Template
	isHX, _ := ctx.Value(contexts.HX).(bool)

	if isHX {
		t = templates.Pages["hx-blog"]
	} else {
		t = templates.Pages["blog"]
	}

	w.Header().Set("Content-Type", "text/html")

	if err := t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
