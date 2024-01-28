package admin

import (
	"gifthub/conf"
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
	slugTpl, err = templates.Build("slug.html").ParseFiles(conf.WorkingSpace + "web/views/admin/slug.html")

	if err != nil {
		log.Panicln(err)
	}
}

func Slug(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("title")
	s := slug.MakeLang(t, conf.DefaultLocale.String())

	data := struct {
		Data struct {
			Slug string
		}
	}{
		Data: struct{ Slug string }{Slug: s},
	}

	if err := slugTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
