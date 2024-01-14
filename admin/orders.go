package admin

import (
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/cookies"
	"gifthub/http/httperrors"
	"gifthub/orders"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
)

var ordersTpl *template.Template
var ordersHxTpl *template.Template

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
}

func Orders(w http.ResponseWriter, r *http.Request) {
	var page int = 1

	ppage := r.URL.Query().Get("page")
	if ppage != "" {
		if d, err := strconv.ParseInt(ppage, 10, 32); err == nil && d > 0 {
			page = int(d)
		}
	}

	q := r.URL.Query().Get("q")
	query := orders.Query{}
	if q != "" {
		query.Keyword = q
	}

	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	offset := (page - 1) * conf.ItemsPerPage
	num := offset + conf.ItemsPerPage

	res, err := orders.Search(ctx, query, offset, num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	pag := templates.Paginate(page, len(res.Orders), int(res.Total))
	pag.URL = "/admin/orders.html"
	pag.Lang = lang

	flash := ""
	c, err := r.Cookie(cookies.FlashMessage)
	if err == nil && c != nil {
		flash = c.Value
	}

	data := struct {
		Lang       language.Tag
		Page       string
		Orders     []orders.Order
		Empty      bool
		Currency   string
		Pagination templates.Pagination
		Flash      string
	}{
		lang,
		"Orders",
		res.Orders,
		len(res.Orders) == 0,
		conf.Currency,
		pag,
		flash,
	}

	isHX, _ := ctx.Value(contexts.HX).(bool)
	var t *template.Template = ordersTpl
	if isHX {
		t = ordersHxTpl
	}

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
