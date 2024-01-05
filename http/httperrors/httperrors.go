package httperrors

import (
	"context"
	"fmt"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/text/language"
)

var itpl *template.Template
var atpl *template.Template
var ptpl *template.Template

func init() {
	var err error

	itpl, err = templates.Build("input-error.html").ParseFiles([]string{
		conf.WorkingSpace + "web/views/admin/input-error.html",
	}...)

	if err != nil {
		log.Panicln(err)
	}

	atpl, err = templates.Build("alert.html").ParseFiles([]string{
		conf.WorkingSpace + "web/views/admin/alert.html",
		conf.WorkingSpace + "web/views/admin/icons/error.svg",
	}...)

	if err != nil {
		log.Panicln(err)
	}

	ptpl, err = templates.Build("base.html").ParseFiles(
		conf.WorkingSpace+"web/views/admin/base.html",
		conf.WorkingSpace+"web/views/admin/error.html",
		conf.WorkingSpace+"web/views/admin/icons/back.svg",
	)

	if err != nil {
		log.Panicln(err)
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	Page(w, r.Context(), "error_http_notfound", 404)
}

func InputMessage(w http.ResponseWriter, ctx context.Context, msg string) {
	lang := ctx.Value(contexts.Locale).(language.Tag)

	data := struct {
		Lang    language.Tag
		Message string
	}{lang, msg}

	key := strings.Split(msg, "_")[1]

	w.Header().Set("HX-Retarget", fmt.Sprintf("#%s-error", key))
	w.Header().Set("HX-Reswap", fmt.Sprintf("innerHTML show:#%s-row:top", key))

	if err := itpl.Execute(w, &data); err != nil {
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
	if strings.HasPrefix(msg, "input_") {
		InputMessage(w, ctx, msg)
	} else {
		Alert(w, ctx, msg)
	}
}

// Alert display an error message through a banner
func Alert(w http.ResponseWriter, ctx context.Context, msg string) {
	lang := ctx.Value(contexts.Locale).(language.Tag)
	rid, _ := ctx.Value(middleware.RequestIDKey).(string)

	data := struct {
		Lang    language.Tag
		Message string
		RID     string
	}{lang, msg, rid}

	w.Header().Set("HX-Replace-Url", "false")
	w.Header().Set("HX-Retarget", "#alert")
	w.Header().Set("HX-Reswap", "innerHTML show:#alert:top")

	if err := atpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

// Page display a full page error for standard http request error
func Page(w http.ResponseWriter, ctx context.Context, msg string, code int) {
	lang := ctx.Value(contexts.Locale).(language.Tag)

	rid := ctx.Value(middleware.RequestIDKey).(string)
	data := struct {
		Lang    language.Tag
		Code    int
		Message string
		RID     string
	}{
		lang,
		code,
		msg,
		rid,
	}

	w.WriteHeader(code)

	ptpl.Execute(w, &data)
}
