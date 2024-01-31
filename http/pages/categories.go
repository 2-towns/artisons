// Package pages provides the application pages
package pages

import (
	"artisons/http/contexts"
	"artisons/shops"
	"artisons/tags"
	"artisons/templates"
	"log/slog"
	"net/http"

	"golang.org/x/text/language"
)

func Categories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	data := struct {
		Lang language.Tag
		Shop shops.Settings
		Tags []tags.Leaf
	}{
		lang,
		shops.Data,
		tags.Tree,
	}

	w.Header().Set("Content-Type", "text/html")

	if err := templates.Pages["categories"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
