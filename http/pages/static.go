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
	"strings"

	"golang.org/x/text/language"
)

func Static(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	slug := strings.Replace(r.URL.Path, ".html", "", 1)
	slug = strings.Replace(slug, "/", "", 1)

	s, err := blog.FindBySlug(ctx, slug)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the static page", slog.String("slug", slug), slog.String("error", err.Error()))
		httperrors.Page(w, r.Context(), err.Error(), 404)
		return
	}

	data := struct {
		Lang    language.Tag
		Shop    shops.Settings
		Article blog.Article
		Tags    []tags.Leaf
	}{
		lang,
		shops.Data,
		s,
		tags.Tree,
	}

	if err := templates.Pages["static"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
