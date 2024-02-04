// Package pages provides the application pages
package pages

import (
	"artisons/carts"
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/orders"
	"artisons/shops"
	"artisons/tags"
	"artisons/templates"
	"log/slog"
	"net/http"

	"golang.org/x/text/language"
)

func Delivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	del, err := shops.Deliveries(ctx)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	data := struct {
		Lang       language.Tag
		Shop       shops.Settings
		Tags       []tags.Leaf
		Deliveries []string
	}{
		lang,
		shops.Data,
		tags.Tree,
		del,
	}

	w.Header().Set("Content-Type", "text/html")

	if err := templates.Pages["delivery"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func DeliverySet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	del := r.FormValue("delivery")
	b := orders.IsValidDelivery(ctx, del)
	if !b {
		httperrors.HXCatch(w, ctx, "the delivery is invalid")
		return
	}

	err := carts.UpdateDelivery(ctx, del)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	w.Header().Set("HX-Redirect", "/payment.html")
	w.Write([]byte(""))

}
