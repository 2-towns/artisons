package carts

import (
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/cookies"
	"artisons/http/httperrors"
	"artisons/http/httphelpers"
	"artisons/orders"
	"artisons/shops"
	"artisons/tags/tree"
	"artisons/templates"
	"artisons/tracking"
	"artisons/users"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
)

func get(r *http.Request) Cart {
	ctx := r.Context()

	u, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		return Cart{ID: fmt.Sprintf("%d", u.ID)}
	}

	id, err := r.Cookie(cookies.CartID)
	if err == nil {
		return Cart{ID: id.Value}
	}

	return Cart{ID: ""}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	c := get(r)

	c, err := c.Get(ctx)
	if err != nil {
		httperrors.Page(w, r.Context(), err.Error(), 500)
		return
	}

	data := struct {
		Lang  language.Tag
		Shop  shops.Settings
		Tags  []tree.Leaf
		Cart  Cart
		Empty bool
	}{
		lang,
		shops.Data,
		tree.Tree,
		c,
		len(c.Products) == 0,
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

func AddHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pid := r.PathValue("id")

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

	c := get(r)

	cid, err := c.Add(ctx, pid, int(qty))
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if cid != "" {
		c := httphelpers.NewCookie(cookies.CartID, cid, int(conf.Cookie.MaxAge))
		http.SetCookie(w, &c)
		r.AddCookie(&c)
	}

	tra := map[string]string{"pid": pid, "quantity": fmt.Sprintf("%d", qty)}
	go tracking.Log(ctx, "cart_add", tra)

	if r.Header.Get("HX-Current-Url") == "/cart" {
		Handler(w, r)
		return
	}

	if shops.Data.Redirect {
		w.Header().Set("HX-Redirect", "/cart")
		w.Write([]byte(""))
		return
	}

	w.Write([]byte(""))
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pid := r.PathValue("id")

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

	c := get(r)

	if !Exists(ctx, c.ID) {
		httperrors.HXCatch(w, ctx, "you are not authorized to process this request")
		return
	}

	err = c.Delete(ctx, pid, int(qty))
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	tra := map[string]string{"pid": pid, "quantity": r.FormValue("quantity")}
	go tracking.Log(ctx, "cart_remove", tra)

	// if r.Header.Get("HX-Current-Url") == "/cart" {
	// 	Handler(w, r)
	// 	return
	// }

	// w.Write([]byte(""))
	Handler(w, r)
}

func DeliveryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c := get(r)

	c, err := c.Get(ctx)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 401)
		return
	}

	if len(c.Products) == 0 {
		slog.LogAttrs(ctx, slog.LevelInfo, "the cart is empty")
		httperrors.Catch(w, ctx, "you are not authorized to process this request", 401)
		return
	}

	del, err := shops.Deliveries(ctx)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := struct {
		Lang       language.Tag
		Shop       shops.Settings
		Tags       []tree.Leaf
		Deliveries []string
	}{
		lang,
		shops.Data,
		tree.Tree,
		del,
	}

	w.Header().Set("Content-Type", "text/html")

	if err := templates.Pages["delivery"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func DeliverySetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c := get(r)

	c, err := c.Get(ctx)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 401)
		return
	}

	if len(c.Products) == 0 {
		slog.LogAttrs(ctx, slog.LevelInfo, "the cart is empty")
		httperrors.HXCatch(w, ctx, "you are not authorized to process this request")
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	del := r.FormValue("delivery")
	b := orders.IsValidDelivery(ctx, del)
	if !b {
		httperrors.HXCatch(w, ctx, "the delivery is invalid")
		return
	}

	err = c.UpdateDelivery(ctx, del)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	tra := map[string]string{"delivery": del}
	go tracking.Log(ctx, "cart_delivery", tra)

	w.Header().Set("HX-Redirect", "/payment.html")
	w.Write([]byte(""))

}
