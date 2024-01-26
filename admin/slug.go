package admin

import (
	"gifthub/conf"
	"gifthub/products"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"

	"github.com/gosimple/slug"
)

var slugTpl *template.Template

func init() {
	var err error
	slugTpl, err = templates.Build("products-slug.html").ParseFiles(conf.WorkingSpace + "web/views/admin/products/products-slug.html")

	if err != nil {
		log.Panicln(err)
	}
}

func Slug(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("title")
	s := slug.MakeLang(t, conf.DefaultLocale.String())

	data := struct {
		Data products.Product
	}{
		Data: products.Product{Slug: s},
	}

	if err := slugTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
