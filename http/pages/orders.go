// Package pages provides the application pages
package pages

import (
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/http/httpext"
	"artisons/orders"
	"artisons/shops"
	"artisons/tags"
	"artisons/templates"
	"artisons/users"
	"html/template"
	"log/slog"
	"net/http"

	"golang.org/x/text/language"
)

func Orders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	user := ctx.Value(contexts.User).(users.User)
	query := orders.Query{UID: user.ID}
	p := httpext.Pagination(r)

	res, err := orders.Search(ctx, query, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	pag := templates.Paginate(p.Page, len(res.Orders), int(res.Total))
	pag.URL = "/account/orders.html"
	pag.Lang = lang

	data := struct {
		Lang       language.Tag
		Shop       shops.Settings
		Tags       []tags.Leaf
		Orders     []orders.Order
		Empty      bool
		Pagination templates.Pagination
	}{
		lang,
		shops.Data,
		tags.Tree,
		res.Orders,
		len(res.Orders) == 0,
		pag,
	}

	var t *template.Template
	isHX, _ := ctx.Value(contexts.HX).(bool)

	if isHX {
		t = templates.Pages["hx-orders"]
	} else {
		t = templates.Pages["orders"]
	}

	w.Header().Set("Content-Type", "text/html")

	if err := t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
