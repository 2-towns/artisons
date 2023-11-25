package admin

import (
	"gifthub/http/contexts"
	"gifthub/locales"
	"html/template"
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	tpl, err := template.ParseFiles("web/views/admin/base.html", "web/views/admin/dashboard.html")
	if err != nil {
		http.Error(w, locales.TranslateError(err, lang), http.StatusInternalServerError)
		return
	}

	p := message.NewPrinter(lang)

	t := map[string]string{
		"title":       p.Sprintf("dashboard_title"),
		"description": p.Sprintf("dashboard_description"),
	}

	data := struct {
		T map[string]string
	}{
		t,
	}

	tpl.Execute(w, &data)

}
