// Package pages provides the application pages
package pages

import (
	"context"
	"gifthub/conf"
	"gifthub/locales"
	"gifthub/products"
	"html/template"
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func homeI18n(lang language.Tag) map[string]string {
	p := message.NewPrinter(lang)

	t := locales.GetPage(lang, "home")
	t["detail"] = p.Sprintf("detail")
	t["product_url"] = p.Sprintf("product_url")

	return t
}

func getProducts(ctx context.Context) ([]products.Product, error) {
	products := make([]products.Product, 0, conf.ItemsPerPage)

	return products, nil
}

// Home loads the most recent products in order to
// display them on the home page.
func Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(locales.ContextKey).(language.Tag)

	tpl, err := template.ParseFiles("web/views/base.html", "web/views/home.html")
	if err != nil {
		http.Error(w, locales.TranslateError(err, lang), http.StatusInternalServerError)
		return
	}

	t := homeI18n(lang)

	p, err := getProducts(ctx)
	if err != nil {
		http.Error(w, locales.TranslateError(err, lang), http.StatusInternalServerError)

		return
	}

	data := struct {
		T        map[string]string
		Products []products.Product
	}{
		t,
		p,
	}

	tpl.Execute(w, &data)
}
