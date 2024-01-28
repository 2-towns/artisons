package pages

import (
	"artisons/blog"
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/shops"
	"artisons/tags"
	"artisons/templates"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

func Article(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	slug := chi.URLParam(r, "slug")

	a, err := blog.FindBySlug(ctx, slug)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the product", slog.String("slug", slug), slog.String("error", err.Error()))
		httperrors.Page(w, r.Context(), err.Error(), 404)
		return
	}

	log.Println(a)

	data := struct {
		Lang    language.Tag
		Shop    shops.Settings
		Article blog.Article
		Tags    []tags.Leaf
	}{
		lang,
		shops.Data,
		a,
		tags.Tree,
	}

	if err := templates.Pages["article"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
