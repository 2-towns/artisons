package admin

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/http/pages"
	"artisons/orders"
	"artisons/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

var ordersTpl *template.Template
var ordersHxTpl *template.Template
var ordersFormTpl *template.Template
var ordersUpdateStatusTpl *template.Template
var ordersNoteAddStatusTpl *template.Template

func init() {
	var err error

	files := append(templates.AdminTable,
		conf.WorkingSpace+"web/views/admin/orders/orders-table.html",
	)

	ordersTpl, err = templates.Build("base.html").ParseFiles(
		append(files, append(templates.AdminList,
			conf.WorkingSpace+"web/views/admin/orders/orders.html")...,
		)...)

	if err != nil {
		log.Panicln(err)
	}

	ordersHxTpl, err = templates.Build("orders-table.html").ParseFiles(files...)

	if err != nil {
		log.Panicln(err)
	}

	ordersFormTpl, err = templates.Build("base.html").ParseFiles(
		append(templates.AdminUI,
			conf.WorkingSpace+"web/views/admin/icons/anchor.svg",
			conf.WorkingSpace+"web/views/admin/orders/orders-form.html",
			conf.WorkingSpace+"web/views/admin/orders/orders-notes.html",
		)...)

	if err != nil {
		log.Panicln(err)
	}

	ordersUpdateStatusTpl, err = templates.Build("alert-success.html").ParseFiles(templates.AdminSuccess...)

	if err != nil {
		log.Panicln(err)
	}

	ordersNoteAddStatusTpl, err = templates.Build("orders-add-note-success.html").ParseFiles(
		append(templates.AdminSuccess,
			conf.WorkingSpace+"web/views/admin/orders/orders-add-note-success.html",
			conf.WorkingSpace+"web/views/admin/orders/orders-notes.html",
		)...,
	)

	if err != nil {
		log.Panicln(err)
	}
}

func OrderList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := ctx.Value(contexts.Pagination).(pages.Paginator)

	qry := orders.Query{}
	if p.Query != "" {
		qry.Keywords = db.Escape(p.Query)
	}

	res, err := orders.Search(ctx, qry, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	t := ordersTpl
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		t = ordersHxTpl
	}

	data := pages.Datalist(ctx, res.Orders)
	data.Pagination = p.Build(ctx, res.Total, len(res.Orders))
	data.Page = "Orders"

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func OrderForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var order orders.Order

	if id != "" {
		var err error
		order, err = orders.Find(ctx, id)

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot find the order", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}
	}

	data := pages.Dataform[orders.Order](ctx, order)
	data.Page = "Orders"

	if err := ordersFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func OrderUpdateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
		"The data has been saved successfully.",
		lang,
	}

	if err := ordersUpdateStatusTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func OrderAddNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
		"The data has been saved successfully.",
		lang,
		o,
	}

	if err := ordersNoteAddStatusTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
