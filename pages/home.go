// Package pages provides the application pages
package pages

import (
	"context"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/locales"
	"gifthub/products"
	"gifthub/templates"
	"net/http"

	"golang.org/x/text/language"
)

func getProducts(ctx context.Context) ([]products.Product, error) {
	products := make([]products.Product, 0, conf.ItemsPerPage)

	return products, nil
}

// Home loads the most recent products in order to
// display them on the home page.
func Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	tpl, err := templates.Build(lang, false).ParseFiles("web/views/base.html", "web/views/home.html")
	if err != nil {
		http.Error(w, locales.TranslateError(err, lang), http.StatusInternalServerError)
		return
	}

	p, err := getProducts(ctx)
	if err != nil {
		http.Error(w, locales.TranslateError(err, lang), http.StatusInternalServerError)

		return
	}

	data := struct {
		Products []products.Product
	}{
		p,
	}

	tpl.Execute(w, &data)
}
