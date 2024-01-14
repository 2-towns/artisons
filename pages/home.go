// Package pages provides the application pages
package pages

import (
	"context"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/httperrors"
	"gifthub/products"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"

	"golang.org/x/text/language"
)

var tpl *template.Template

func init() {
	var err error

	tpl, err = templates.Build("home.html").ParseFiles([]string{
		"web/views/base.html",
		"web/views/home.html",
	}...)

	if err != nil {
		log.Panicln(err)
	}
}

func getProducts(ctx context.Context) ([]products.Product, error) {
	products := make([]products.Product, 0, conf.ItemsPerPage)

	return products, nil
}

// Home loads the most recent products in order to
// display them on the home page.
func Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	p, err := getProducts(ctx)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the products", slog.String("error", err.Error()))
		httperrors.Page(w, r.Context(), "error_http_general", 400)
		return
	}

	data := struct {
		Lang     language.Tag
		Products []products.Product
	}{
		lang,
		p,
	}
	log.Println("fdfdsd!!!!!!!!!!!!!!!!!!")

	tpl.Execute(w, &data)
}
