package orders

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/http/httphelpers"
	"artisons/shops"
	"artisons/tags/tree"
	"artisons/templates"
	"artisons/users"
	"html/template"
	"log"
	"log/slog"
	"net/http"

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
		append(files, append(templates.AdminListHandler,
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

func OrderListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := httphelpers.BuildPaginator(r)

	qry := Query{}
	if p.Query != "" {
		qry.Keywords = db.Escape(p.Query)
	}

	res, err := Search(ctx, qry, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	t := ordersTpl
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		t = ordersHxTpl
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := httphelpers.List[Order]{
		Lang:       lang,
		Items:      res.Orders,
		Empty:      len(res.Orders) == 0,
		Currency:   conf.Currency,
		Pagination: p.Build(ctx, res.Total, len(res.Orders)),
		Page:       "Orders",
		Flash:      httphelpers.Flash(w, r),
	}

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func OrderFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	var order Order

	if id != "" {
		var err error
		order, err = Find(ctx, id)

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot find the order", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := httphelpers.Form[Order]{
		Data:     order,
		Lang:     lang,
		Currency: conf.Currency,
		Page:     "Orders",
	}

	if err := ordersFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func OrderUpdateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	oid := r.PathValue("id")
	status := r.FormValue("status")

	err := UpdateStatus(ctx, oid, status)
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

func OrderAddNoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	oid := r.PathValue("id")
	note := r.FormValue("note")

	err := AddNote(ctx, oid, note)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	o, err := Find(ctx, oid)
	if err != nil {
		httperrors.Page(w, ctx, err.Error(), 400)
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)

	data := struct {
		Flash string
		Lang  language.Tag
		Data  Order
	}{
		"The data has been saved successfully.",
		lang,
		o,
	}

	if err := ordersNoteAddStatusTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	id := r.PathValue("id")

	order, err := Find(ctx, id)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	data := struct {
		Lang  language.Tag
		Shop  shops.Settings
		Tags  []tree.Leaf
		Order Order
	}{
		lang,
		shops.Data,
		tree.Tree,
		order,
	}

	if err := templates.Pages["order"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	user := ctx.Value(contexts.User).(users.User)
	query := Query{UID: user.ID}
	p := httphelpers.BuildPaginator(r)

	res, err := Search(ctx, query, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	pag := p.Build(ctx, res.Total, len(res.Orders))

	data := struct {
		Lang       language.Tag
		Shop       shops.Settings
		Tags       []tree.Leaf
		Orders     []Order
		Empty      bool
		Pagination httphelpers.Pagination
	}{
		lang,
		shops.Data,
		tree.Tree,
		res.Orders,
		len(res.Orders) == 0,
		pag,
	}

	var t *template.Template
	isHX, _ := ctx.Value(contexts.HX).(bool)

	if isHX {
		t = templates.Pages["hx-orders"]
	} else {
		t = templates.Pages["orders"]
	}

	w.Header().Set("Content-Type", "text/html")

	if err := t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
