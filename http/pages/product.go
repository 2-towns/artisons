package pages

import (
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/products"
	"artisons/shops"
	"artisons/tags"
	"artisons/templates"
	"artisons/users"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

func Product(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	slug := chi.URLParam(r, "slug")

	p, err := products.FindBySlug(ctx, slug)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the product", slog.String("slug", slug), slog.String("error", err.Error()))
		httperrors.Page(w, r.Context(), err.Error(), 404)
		return
	}

	wish := false
	user, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		wish = user.HasWish(ctx, p.ID)
	}

	data := struct {
		Lang    language.Tag
		Shop    shops.Settings
		Product products.Product
		Tags    []tags.Leaf
		Wish    bool
	}{
		lang,
		shops.Data,
		p,
		tags.Tree,
		wish,
	}

	if err := templates.Pages["product"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
