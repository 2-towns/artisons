package routes

import (
	"context"
	"gifthub/util"
	"html/template"
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func homeI18n(lang language.Tag) map[string]string {
	p := message.NewPrinter(lang)

	t := util.GetPage(lang, "home")
	t["detail"] = p.Sprintf("detail")
	t["product_url"] = p.Sprintf("product_url")

	return t
}

func getProducts(ctx context.Context) ([]util.Product, error) {
	products := make([]util.Product, 0, util.ItemsPerPage)

	return products, nil
}

func HomeRoute(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("views/base.html", "views/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	lang := ctx.Value(util.ContextLangKey).(language.Tag)
	t := homeI18n(lang)

	products, err := getProducts(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	data := struct {
		T        map[string]string
		Products []util.Product
	}{
		t,
		products,
	}

	tpl.Execute(w, &data)
}
