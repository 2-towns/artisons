package admin

import (
	"gifthub/http/contexts"
	"gifthub/http/httperrors"
	"gifthub/templates"
	"log/slog"
	"net/http"

	"golang.org/x/text/language"
)

func Products(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"web/views/admin/base.html",
		"web/views/admin/ui.html",
		"web/views/admin/icons/home.svg",
		"web/views/admin/icons/building-store.svg",
		"web/views/admin/products/products.html",
	}

	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	tpl, err := templates.Build(lang, false).ParseFiles(files...)

	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	data := struct {
		Page string
	}{
		"products",
	}
	if err = tpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
