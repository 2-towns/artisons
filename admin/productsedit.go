package admin

import (
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/cookies"
	"gifthub/http/httperrors"
	"gifthub/http/httpext"
	"gifthub/products"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

var productsEditTpl *template.Template

func init() {
	var err error

	productsEditTpl, err = templates.Build("base.html").ParseFiles([]string{
		conf.WorkingSpace + "web/views/admin/base.html",
		conf.WorkingSpace + "web/views/admin/ui.html",
		conf.WorkingSpace + "web/views/admin/icons/home.svg",
		conf.WorkingSpace + "web/views/admin/icons/close.svg",
		conf.WorkingSpace + "web/views/admin/icons/building-store.svg",
		conf.WorkingSpace + "web/views/admin/products/products-add.html",
		conf.WorkingSpace + "web/views/admin/products/products-form.html",
	}...)

	if err != nil {
		log.Panicln(err)
	}
}

func EditProductForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	id := chi.URLParam(r, "id")

	p, err := products.Find(ctx, id)
	if err != nil {
		httperrors.Page(w, ctx, "error_http_productnotfound", 404)
		return
	}

	data := struct {
		Lang language.Tag
		Page string
		ID   string
		Data products.Product
	}{
		lang,
		"products",
		id,
		p,
	}

	if err := productsEditTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func EditProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(conf.MaxUploadSize); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "error_http_general")
		return
	}

	pid := chi.URLParam(r, "id")
	p, err := processProductFrom(ctx, *r.MultipartForm, pid)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	err = p.Save(ctx)
	if err != nil {
		httpext.RollbackUpload(ctx, []string{p.Image1, p.Image2, p.Image3, p.Image4})
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     cookies.FlashMessage,
		Value:    "label_dashboard_producteditsuccess",
		MaxAge:   int(time.Minute.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   conf.Cookie.Secure,
		Domain:   conf.Cookie.Domain,
	}

	http.SetCookie(w, cookie)
	w.Header().Set("HX-Redirect", "/admin/products.html")
	w.Write([]byte(""))
}
