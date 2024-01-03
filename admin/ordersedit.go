package admin

import (
	"fmt"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/httperrors"
	"gifthub/orders"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

var ordersEditTpl *template.Template

func init() {
	var err error

	ordersEditTpl, err = templates.Build("base.html").ParseFiles([]string{
		conf.WorkingSpace + "web/views/admin/base.html",
		conf.WorkingSpace + "web/views/admin/ui.html",
		conf.WorkingSpace + "web/views/admin/icons/home.svg",
		conf.WorkingSpace + "web/views/admin/icons/close.svg",
		conf.WorkingSpace + "web/views/admin/icons/anchor.svg",
		conf.WorkingSpace + "web/views/admin/icons/building-store.svg",
		conf.WorkingSpace + "web/views/admin/icons/receipt.svg",
		conf.WorkingSpace + "web/views/admin/icons/settings.svg",
		conf.WorkingSpace + "web/views/admin/icons/article.svg",
		conf.WorkingSpace + "web/views/admin/icons/close.svg",
		conf.WorkingSpace + "web/views/admin/alert-success.html",
		conf.WorkingSpace + "web/views/admin/orders/orders-edit.html",
		conf.WorkingSpace + "web/views/admin/orders/orders-notes.html",
	}...)

	if err != nil {
		log.Panicln(err)
	}
}

func EditOrderForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	id := chi.URLParam(r, "id")

	o, err := orders.Find(ctx, id)
	if err != nil {
		httperrors.Page(w, ctx, err.Error(), 400)
		return
	}

	o, err = o.WithProducts(ctx)
	if err != nil {
		httperrors.Page(w, ctx, err.Error(), 400)
		return
	}

	data := struct {
		Lang     language.Tag
		Page     string
		ID       string
		Data     orders.Order
		Currency string
	}{
		lang,
		"orders",
		id,
		o,
		conf.Currency,
	}

	policy := fmt.Sprintf("default-src 'self'; img-src 'self' %s ;", conf.ImgProxy.URL)
	w.Header().Set("Content-Security-Policy", policy)

	if err := ordersEditTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
