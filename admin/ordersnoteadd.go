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

var ordersNoteAddStatusTpl *template.Template

func init() {
	var err error

	ordersNoteAddStatusTpl, err = templates.Build("orders-add-note-success.html").ParseFiles(
		conf.WorkingSpace+"web/views/admin/icons/success.svg",
		conf.WorkingSpace+"web/views/admin/alert-success.html",
		conf.WorkingSpace+"web/views/admin/orders/orders-add-note-success.html",
		conf.WorkingSpace+"web/views/admin/orders/orders-notes.html",
	)

	if err != nil {
		log.Panicln(err)
	}
}

func AddOrderNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "error_http_general")
		return
	}

	oid := chi.URLParam(r, "id")
	note := r.FormValue("note")

	err := orders.AddNote(ctx, oid, note)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	id := chi.URLParam(r, "id")
	o, err := orders.Find(ctx, id)
	if err != nil {
		httperrors.Page(w, ctx, err.Error(), 400)
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)

	data := struct {
		Flash string
		Lang  language.Tag
		Data  orders.Order
	}{
		"text_general_ordersnoteadded",
		lang,
		o,
	}

	if err := ordersNoteAddStatusTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
