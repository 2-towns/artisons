package admin

import (
	"context"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/contexts"
	"gifthub/http/httperrors"
	"gifthub/http/httpext"
	"gifthub/orders"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

const ordersName = "Orders"
const ordersURL = "/admin/orders.html"

var ordersTpl *template.Template
var ordersHxTpl *template.Template
var ordersFormTpl *template.Template
var ordersUpdateStatusTpl *template.Template
var ordersNoteAddStatusTpl *template.Template

type ordersFeature struct{}

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

func (f ordersFeature) ListTemplate(ctx context.Context) *template.Template {
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		return ordersHxTpl
	}

	return ordersTpl
}

func (f ordersFeature) Search(ctx context.Context, q string, offset, num int) (httpext.SearchResults[orders.Order], error) {
	query := orders.Query{}
	if q != "" {
		query.Keyword = db.Escape(q)
	}

	res, err := orders.Search(ctx, query, offset, num)

	return httpext.SearchResults[orders.Order]{
		Total: res.Total,
		Items: res.Orders,
	}, err
}

func (f ordersFeature) Find(ctx context.Context, id interface{}) (orders.Order, error) {
	return orders.Find(ctx, id.(string))
}

func (f ordersFeature) FormTemplate(ctx context.Context, w http.ResponseWriter) *template.Template {
	return ordersFormTpl
}

func (f ordersFeature) ID(ctx context.Context, id string) (interface{}, error) {
	return id, nil
}

func OrdersList(w http.ResponseWriter, r *http.Request) {
	httpext.DigestList[orders.Order](w, r, httpext.List[orders.Order]{
		Name:    ordersName,
		URL:     ordersURL,
		Feature: ordersFeature{},
	})
}

func OrdersForm(w http.ResponseWriter, r *http.Request) {
	httpext.DigestForm[orders.Order](w, r, httpext.Form[orders.Order]{
		Name:    ordersName,
		Feature: ordersFeature{},
	})
}

func OrdersUpdateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
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
		"The data has been saved successfully.",
		lang,
	}

	if err := ordersUpdateStatusTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func OrdersAddNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
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
		"The data has been saved successfully.",
		lang,
		o,
	}

	if err := ordersNoteAddStatusTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
