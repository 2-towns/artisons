package httperrors

import (
	"context"
	"fmt"
	"gifthub/admin/urls"
	"gifthub/http/contexts"
	"gifthub/locales"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	Page(w, r.Context(), "error_http_notfound", 404)
}

func InputMessage(w http.ResponseWriter, ctx context.Context, msg string) {
	lang := ctx.Value(contexts.Locale).(language.Tag)

	tpl, err := template.ParseFiles("web/views/admin/input-error.html")
	if err != nil {
		http.Error(w, locales.TranslateError(err, lang), http.StatusInternalServerError)
		return
	}

	p := message.NewPrinter(lang)
	data := struct {
		Message string
	}{
		p.Sprintf(msg),
	}

	key := strings.Split(msg, "_")[1]

	w.Header().Set("HX-Retarget", fmt.Sprintf("#%s-error", key))
	w.Header().Set("HX-Reswap", "innerHTML")

	tpl.Execute(w, &data)
}

func Catch(w http.ResponseWriter, ctx context.Context, msg string) {
	if strings.HasPrefix("input_", msg) {
		InputMessage(w, ctx, msg)
	} else {
		Alert(w, ctx, msg)
	}
}

func Alert(w http.ResponseWriter, ctx context.Context, msg string) {
	lang := ctx.Value(contexts.Locale).(language.Tag)

	tpl, err := template.ParseFiles("web/views/admin/alert.html", "web/views/admin/icons/error.svg")
	if err != nil {
		http.Error(w, locales.TranslateError(err, lang), http.StatusInternalServerError)
		return
	}

	p := message.NewPrinter(lang)
	t := map[string]string{
		"title_error":   p.Sprintf("title_error_common"),
		"message_error": p.Sprintf(msg),
	}

	data := struct {
		T map[string]string
	}{
		t,
	}

	w.Header().Set("HX-Replace-Url", "false")
	w.Header().Set("HX-Retarget", "#alert")
	w.Header().Set("HX-Reswap", "innerHTML")

	tpl.Execute(w, &data)
}

func Page(w http.ResponseWriter, ctx context.Context, msg string, code int) {
	lang := ctx.Value(contexts.Locale).(language.Tag)

	tpl, err := template.ParseFiles("web/views/admin/base.html", "web/views/admin/error.html", "web/views/admin/icons/back.html")
	if err != nil {
		http.Error(w, locales.TranslateError(err, lang), http.StatusInternalServerError)
		return
	}

	rid := ctx.Value(middleware.RequestIDKey).(string)

	p := message.NewPrinter(lang)
	t := map[string]string{
		"title":       p.Sprintf("error_http_page"),
		"request":     p.Sprintf("error_http_requestid", rid),
		"description": p.Sprintf(msg),
		"message":     p.Sprintf(msg, rid),
		"home_button": p.Sprintf("label_button_error"),
	}

	url := urls.AdminPrefix

	data := struct {
		Code int
		Link string
		T    map[string]string
	}{
		code,
		url,
		t,
	}

	w.WriteHeader(code)

	tpl.Execute(w, &data)
}
