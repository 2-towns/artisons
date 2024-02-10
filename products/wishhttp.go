package products

import (
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/shops"
	"artisons/tags/tree"
	"artisons/templates"
	"artisons/users"
	"log/slog"
	"net/http"

	"golang.org/x/text/language"
)

func WishesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	wishes := []string{}
	user, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		var err error
		wishes, err = Wishes(ctx, user.ID)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the wishes", slog.String("error", err.Error()))
		}
	}

	pds, err := List(ctx, wishes)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the products", slog.String("error", err.Error()))
		httperrors.Page(w, r.Context(), "something went wrong", 400)
		return
	}

	data := struct {
		Lang     language.Tag
		Shop     shops.Settings
		Products []Product
		Tags     []tree.Leaf
		Empty    bool
	}{
		lang,
		shops.Data,
		pds,
		tree.Tree,
		len(pds) == 0,
	}

	t := templates.Pages["wish"]
	isHX, _ := ctx.Value(contexts.HX).(bool)

	if isHX {
		t = templates.Pages["hx-wish"]
	}

	w.Header().Set("Content-Type", "text/html")

	if err := t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func WishHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")

	user, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		err := Wish(ctx, user.ID, id)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot add to wish", slog.String("error", err.Error()))
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}
	}

	WishesHandler(w, r)
}

func UnWishHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")

	user, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		err := UnWish(ctx, user.ID, id)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot add to wish", slog.String("error", err.Error()))
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}
	}

	WishesHandler(w, r)
}
