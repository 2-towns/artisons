package pages

import (
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/products"
	"artisons/shops"
	"artisons/tags"
	"artisons/templates"
	"artisons/users"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

func Wishes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	wishes := []string{}
	user, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		var err error
		wishes, err = user.Wishes(ctx)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the wishes", slog.String("error", err.Error()))
		}
	}

	pds, err := products.List(ctx, wishes)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the products", slog.String("error", err.Error()))
		httperrors.Page(w, r.Context(), "something went wrong", 400)
		return
	}

	data := struct {
		Lang     language.Tag
		Shop     shops.Settings
		Products []products.Product
		Tags     []tags.Leaf
		Empty    bool
	}{
		lang,
		shops.Data,
		pds,
		tags.Tree,
		len(pds) == 0,
	}

	var t *template.Template
	isHX, _ := ctx.Value(contexts.HX).(bool)

	if isHX {
		t = templates.Pages["hx-wish"]
	} else {
		t = templates.Pages["wish"]
	}

	if err := t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func Wish(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	user, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		err := user.Wish(ctx, id)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot add to wish", slog.String("error", err.Error()))
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}
	}

	Wishes(w, r)
}

func UnWish(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	user, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		err := user.UnWish(ctx, id)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot add to wish", slog.String("error", err.Error()))
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}
	}

	Wishes(w, r)
}
