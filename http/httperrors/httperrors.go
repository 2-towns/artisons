package httperrors

import (
	"context"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/locales"
	"html/template"
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Page(w http.ResponseWriter, ctx context.Context, msg string, code int) {
	lang := ctx.Value(contexts.Locale).(language.Tag)

	tpl, err := template.ParseFiles("web/views/admin/base.html", "web/views/error.html", "web/views/icons/back.html")
	if err != nil {
		http.Error(w, locales.TranslateError(err, lang), http.StatusInternalServerError)
		return
	}

	p := message.NewPrinter(lang)

	t := map[string]string{
		"title":       p.Sprintf("error_title"),
		"description": p.Sprintf(msg),
		"message":     p.Sprintf(msg),
		"home_button": p.Sprintf("error_home_button"),
	}

	url := conf.AdminPrefix

	if code == 401 {
		url = "/"
	}

	data := struct {
		Code int
		Link string
		T    map[string]string
	}{
		code,
		url,
		t,
	}

	tpl.Execute(w, &data)

}
