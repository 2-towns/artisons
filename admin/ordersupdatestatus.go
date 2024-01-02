package admin

import (
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

var ordersUpdateStatusTpl *template.Template

func init() {
	var err error

	ordersUpdateStatusTpl, err = templates.Build("alert-success.html").ParseFiles(
		conf.WorkingSpace+"web/views/admin/icons/success.svg",
		conf.WorkingSpace+"web/views/admin/alert-success.html",
	)

	if err != nil {
		log.Panicln(err)
	}
}

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "error_http_general")
		return
	}

	oid := chi.URLParam(r, "id")
	status := r.FormValue("status")

	err := orders.UpdateStatus(ctx, oid, status)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)

	data := struct {
		Flash string
		Lang  language.Tag
	}{
		"label_dashboard_ordersstatusupdated",
		lang,
	}

	if err := ordersUpdateStatusTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
