// Package pages provides the application pages
package pages

import (
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/products"
	"artisons/shops"
	"artisons/tags"
	"artisons/templates"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
)

func Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	p := ctx.Value(contexts.Pagination).(Paginator)

	q := r.URL.Query()
	var min float32 = 0
	if q.Has("min") {
		if val, err := strconv.ParseFloat(q.Get("min"), 32); err == nil {
			min = float32(val)
		}
	}

	var max float32 = 0
	if q.Has("max") {
		if val, err := strconv.ParseFloat(q.Get("max"), 32); err == nil {
			max = float32(val)
		}
	}

	meta := map[string][]string{}
	for key, val := range q {
		if key == "min" || key == "max" || key == "q" || key == "tags" {
			continue
		}

		meta[key] = val
	}

	query := products.Query{
		PriceMin: min,
		PriceMax: max,
		Keywords: q.Get("q"),
		Tags:     q["tags"],
		Meta:     meta,
	}

	res, err := products.Search(ctx, query, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	pag := p.Build(ctx, res.Total, len(res.Products))

	data := struct {
		Lang       language.Tag
		Shop       shops.Settings
		Tags       []tags.Leaf
		Products   []products.Product
		Empty      bool
		Pagination Pagination
	}{
		lang,
		shops.Data,
		tags.Tree,
		res.Products,
		len(res.Products) == 0,
		pag,
	}

	var t *template.Template
	isHX, _ := ctx.Value(contexts.HX).(bool)

	if isHX {
		t = templates.Pages["hx-search"]
	} else {
		t = templates.Pages["search"]
	}

	w.Header().Set("Content-Type", "text/html")

	if err := t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
