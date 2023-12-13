package httperrors

import (
	"context"
	"fmt"
	"gifthub/admin/urls"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/locales"
	"gifthub/templates"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/text/language"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	Page(w, r.Context(), "error_http_notfound", 404)
}

func InputMessage(w http.ResponseWriter, ctx context.Context, msg string) {
	lang := ctx.Value(contexts.Locale).(language.Tag)

	tpl, err := templates.Build(lang, true).ParseFiles(
		"web/views/admin/htmx.html",
		"web/views/admin/input-error.html",
	)
	if err != nil {
		http.Error(w, locales.TranslateError(err, lang), http.StatusInternalServerError)
		return
	}

	data := struct{ Message string }{msg}

	key := strings.Split(msg, "_")[1]

	w.Header().Set("HX-Retarget", fmt.Sprintf("#%s-error", key))
	w.Header().Set("HX-Reswap", "innerHTML")

	if err = tpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

// Catch an error which can come from ajax (HTMX) or from a standard
// http request.
func Catch(w http.ResponseWriter, ctx context.Context, msg string, code int) {
	isHX, _ := ctx.Value(contexts.HX).(bool)

	if isHX {
		HXCatch(w, ctx, msg)
	} else {
		Page(w, ctx, msg, code)
	}
}

// Catch an ajax error. If the error starts with "input_", the error
// is related to an wrong input value, so this input will be updated.
// Otherwise an alert will be displayed.
func HXCatch(w http.ResponseWriter, ctx context.Context, msg string) {
	if strings.HasPrefix("input_", msg) {
		InputMessage(w, ctx, msg)
	} else {
		Alert(w, ctx, msg)
	}
}

// Alert display an error message through a banner
func Alert(w http.ResponseWriter, ctx context.Context, msg string) {
	lang := ctx.Value(contexts.Locale).(language.Tag)

	tpl, err := templates.Build(lang, true).ParseFiles(
		"web/views/admin/htmx.html",
		"web/views/admin/alert.html",
		"web/views/admin/icons/error.svg",
	)
	if err != nil {
		http.Error(w, locales.TranslateError(err, lang), http.StatusInternalServerError)
		return
	}

	data := struct{ Message string }{msg}

	w.Header().Set("HX-Replace-Url", "false")
	w.Header().Set("HX-Retarget", "#alert")
	w.Header().Set("HX-Reswap", "innerHTML")

	if err = tpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

// Page display a full page error for standard http request error
func Page(w http.ResponseWriter, ctx context.Context, msg string, code int) {
	lang := ctx.Value(contexts.Locale).(language.Tag)

	tpl, err := templates.Build(lang, false).ParseFiles(
		conf.WorkingSpace+"web/views/admin/base.html",
		conf.WorkingSpace+"web/views/admin/error.html",
		conf.WorkingSpace+"web/views/admin/icons/back.svg",
	)
	if err != nil {
		log.Println(err)
		http.Error(w, locales.TranslateError(err, lang), http.StatusInternalServerError)
		return
	}

	rid := ctx.Value(middleware.RequestIDKey).(string)

	url := urls.AdminPrefix

	data := struct {
		Code    int
		Link    string
		Message string
		RID     string
	}{
		code,
		url,
		msg,
		rid,
	}

	w.WriteHeader(code)

	tpl.Execute(w, &data)
}
