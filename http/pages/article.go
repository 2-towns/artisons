package pages

import (
	"artisons/blog"
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/shops"
	"artisons/tags"
	"artisons/templates"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

func Article(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	slug := chi.URLParam(r, "slug")

	query := blog.Query{Slug: slug}
	res, err := blog.Search(ctx, query, 0, 1)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the article", slog.String("slug", slug), slog.String("error", err.Error()))
		httperrors.Page(w, r.Context(), err.Error(), 400)
		return
	}

	if res.Total == 0 {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot find the article", slog.String("slug", slug))
		httperrors.Page(w, r.Context(), "oops the data is not found", 404)
		return
	}

	a := res.Articles[0]

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

	if err := templates.Pages["static"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
