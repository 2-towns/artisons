// Package pages provides the application pages
package pages

import (
	"gifthub/http/contexts"
	"gifthub/shops"
	"gifthub/tags"
	"gifthub/templates"
	"log/slog"
	"net/http"

	"golang.org/x/text/language"
)

func Account(w http.ResponseWriter, r *http.Request) {
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

	if err := templates.Pages["home"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
