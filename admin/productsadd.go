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

	"golang.org/x/text/language"
)

var productsAddTpl *template.Template

func init() {
	var err error

	productsAddTpl, err = templates.Build("base.html").ParseFiles([]string{
		conf.WorkingSpace + "web/views/admin/base.html",
		conf.WorkingSpace + "web/views/admin/ui.html",
		conf.WorkingSpace + "web/views/admin/icons/home.svg",
		conf.WorkingSpace + "web/views/admin/icons/building-store.svg",
		conf.WorkingSpace + "web/views/admin/products/products-add.html",
		conf.WorkingSpace + "web/views/admin/products/products-form.html",
	}...)

	if err != nil {
		log.Panicln(err)
	}
}

func AddProductForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	data := struct {
		Lang    language.Tag
		Page    string
		ID      string
		Data    products.Product
		Picture string
		Images  []string
	}{
		lang,
		"products",
		"",
		products.Product{},
		"",
		[]string{},
	}

	if err := productsAddTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(conf.MaxUploadSize); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "error_http_general")
		return
	}

	pid := ""
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
		Value:    "label_dashboard_productaddsuccess",
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
