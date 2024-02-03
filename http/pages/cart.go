// Package pages provides the application pages
package pages

import (
	"artisons/carts"
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/cookies"
	"artisons/http/httperrors"
	"artisons/shops"
	"artisons/tags"
	"artisons/templates"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

func Cart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	cart, err := carts.Get(ctx)
	if err != nil {
		httperrors.Page(w, r.Context(), err.Error(), 500)
		return
	}

	data := struct {
		Lang  language.Tag
		Shop  shops.Settings
		Tags  []tags.Leaf
		Cart  carts.Cart
		Empty bool
	}{
		lang,
		shops.Data,
		tags.Tree,
		cart,
		len(cart.Products) == 0,
	}

	var t *template.Template
	isHX, _ := ctx.Value(contexts.HX).(bool)

	if isHX {
		t = templates.Pages["hx-cart"]
	} else {
		t = templates.Pages["cart"]
	}

	w.Header().Set("Content-Type", "text/html")

	if err := t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func CartAdd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	qty, err := strconv.ParseInt(r.FormValue("quantity"), 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the quantity", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "input:quantity")
		return
	}

	cid, err := carts.Add(ctx, id, int(qty))
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if cid != "" {
		cookie := &http.Cookie{
			Name:     cookies.CartID,
			Value:    cid,
			MaxAge:   int(conf.Cookie.MaxAge),
			Path:     "/",
			HttpOnly: true,
			Secure:   conf.Cookie.Secure,
			Domain:   conf.Cookie.Domain,
			SameSite: http.SameSiteStrictMode,
		}

		http.SetCookie(w, cookie)
	}

	if r.Header.Get("HX-Current-URL") == "/cart.html" {
		Cart(w, r)
		return
	}

	if shops.Data.Redirect {
		w.Header().Set("HX-Redirect", "/cart.html")
		w.Write([]byte(""))
		return
	}

	w.Write([]byte(""))
}

func CartDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	qty, err := strconv.ParseInt(r.FormValue("quantity"), 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the quantity", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "input:quantity")
		return
	}

	err = carts.Delete(ctx, id, int(qty))
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if r.Header.Get("HX-Current-URL") == "/cart.html" {
		Cart(w, r)
		return
	}

	if shops.Data.Redirect {
		w.Header().Set("HX-Redirect", "/cart.html")
		w.Write([]byte(""))
		return
	}

	w.Write([]byte(""))
}
