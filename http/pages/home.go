// Package pages provides the application pages
package pages

import (
	"gifthub/http/contexts"
	"gifthub/products"
	"gifthub/shops"
	"gifthub/tags"
	"gifthub/templates"
	"gifthub/users"
	"log/slog"
	"net/http"

	"golang.org/x/text/language"
)

// Home loads the most recent products in order to
// display them on the home page.
func Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	p := []products.Product{}
	// if err != nil {
	// 	slog.LogAttrs(ctx, slog.LevelError, "cannot get the products", slog.String("error", err.Error()))
	// 	httperrors.Page(w, r.Context(), "something went wrong", 400)
	// 	return
	// }

	wishes := []string{}
	user, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		var err error
		wishes, err = user.Wishes(ctx)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the wishes", slog.String("error", err.Error()))
		}
	}

	data := struct {
		Lang     language.Tag
		Shop     shops.Settings
		Products []products.Product
		Tags     []tags.Leaf
		Wishes   []string
	}{
		lang,
		shops.Data,
		p,
		tags.Tree,
		wishes,
	}

	if err := templates.Pages["home"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
