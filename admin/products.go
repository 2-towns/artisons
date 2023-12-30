package admin

import (
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/cookies"
	"gifthub/http/httperrors"
	"gifthub/products"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
)

var productsTpl *template.Template
var productsHxTpl *template.Template

func init() {
	var err error

	files := []string{
		conf.WorkingSpace + "web/views/admin/icons/arrow-right.svg",
		conf.WorkingSpace + "web/views/admin/icons/arrow-left.svg",
		conf.WorkingSpace + "web/views/admin/icons/trash.svg",
		conf.WorkingSpace + "web/views/admin/icons/edit.svg",
		conf.WorkingSpace + "web/views/admin/icons/question-mark.svg",
		conf.WorkingSpace + "web/views/admin/icons/success.svg",
		conf.WorkingSpace + "web/views/admin/products/products-table.html",
		conf.WorkingSpace + "web/views/admin/pagination.html",
	}

	productsTpl, err = templates.Build("base.html").ParseFiles(append(files, []string{
		conf.WorkingSpace + "web/views/admin/base.html",
		conf.WorkingSpace + "web/views/admin/ui.html",
		conf.WorkingSpace + "web/views/admin/products/products-actions.html",
		conf.WorkingSpace + "web/views/admin/icons/home.svg",
		conf.WorkingSpace + "web/views/admin/icons/building-store.svg",
		conf.WorkingSpace + "web/views/admin/products/products.html",
	}...)...)

	if err != nil {
		log.Panicln(err)
	}

	productsHxTpl, err = templates.Build("products-table.html").ParseFiles(files...)

	if err != nil {
		log.Panicln(err)
	}
}

func Products(w http.ResponseWriter, r *http.Request) {
	var page int = 1

	ppage := r.URL.Query().Get("page")
	if ppage != "" {
		if d, err := strconv.ParseInt(ppage, 10, 32); err == nil && d > 0 {
			page = int(d)
		}
	}

	q := r.URL.Query().Get("q")
	query := products.Query{}
	if q != "" {
		query.Keywords = q
	}

	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	offset := (page - 1) * conf.ItemsPerPage
	num := offset + conf.ItemsPerPage

	res, err := products.Search(ctx, query, offset, num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	pag := templates.Paginate(page, len(res.Products), int(res.Total))
	pag.URL = "/admin/products.html"
	pag.Lang = lang

	flash := ""
	c, err := r.Cookie(cookies.FlashMessage)
	if err != nil && c != nil {
		flash = c.Value
	}

	data := struct {
		Lang       language.Tag
		Page       string
		Products   []products.Product
		Empty      bool
		Currency   string
		Pagination templates.Pagination
		Flash      string
	}{
		lang,
		"products",
		res.Products,
		len(res.Products) == 0,
		conf.Currency,
		pag,
		flash,
	}

	isHX, _ := ctx.Value(contexts.HX).(bool)
	var t *template.Template = productsTpl
	if isHX {
		t = productsHxTpl
	}

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
