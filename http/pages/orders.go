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
	pag.URL = "/blog.html"
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

	if err := templates.Pages["orders"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
