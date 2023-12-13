package templates

import (
	"gifthub/admin/urls"
	"gifthub/cache"
	"gifthub/locales"
	"html/template"

	"golang.org/x/text/language"
)

func Build(lang language.Tag, ixHX bool) *template.Template {
	name := "base.html"
	if ixHX {
		name = "htmx.html"
	}

	return template.New(name).Funcs(template.FuncMap{
		"translate":   locales.Translate(lang),
		"cachebuster": cache.Buster,
		"urls":        urls.Get,
	})
}
