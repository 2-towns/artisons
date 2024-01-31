// Package pages provides the application pages
package pages

import (
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/orders"
	"artisons/shops"
	"artisons/tags"
	"artisons/templates"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

func Order(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	id := chi.URLParam(r, "id")

	order, err := orders.Find(ctx, id)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	data := struct {
		Lang  language.Tag
		Shop  shops.Settings
		Tags  []tags.Leaf
		Order orders.Order
	}{
		lang,
		shops.Data,
		tags.Tree,
		order,
	}

	if err := templates.Pages["order"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
