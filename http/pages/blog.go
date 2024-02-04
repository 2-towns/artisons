// Package pages provides the application pages
package pages

import (
	"artisons/blog"
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/shops"
	"artisons/tags"
	"artisons/templates"
	"html/template"
	"log/slog"
	"net/http"

	"golang.org/x/text/language"
)

func Blog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	p := ctx.Value(contexts.Pagination).(Paginator)

	query := blog.Query{
		Keywords: p.Query,
		Type:     "blog",
	}

	res, err := blog.Search(ctx, query, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	pag := p.Build(ctx, res.Total, len(res.Articles))

	data := struct {
		Lang       language.Tag
		Shop       shops.Settings
		Articles   []blog.Article
		Empty      bool
		Pagination Pagination
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
