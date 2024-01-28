package auth

import (
	"fmt"
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/cookies"
	"artisons/http/httperrors"
	"artisons/http/httpext"
	"artisons/http/security"
	"artisons/templates"
	"artisons/users"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"golang.org/x/text/language"
)

var tpl *template.Template
var hxtpl *template.Template
var otptpl *template.Template

func init() {
	var err error

	tpl, err = templates.Build("base.html").ParseFiles([]string{
		"web/views/admin/base.html",
		"web/views/admin/simple.html",
		"web/views/login/login.html",
		"web/views/admin/icons/logo.svg",
	}...)

	if err != nil {
		log.Panicln(err)
	}

	hxtpl, err = templates.Build("login.html").ParseFiles([]string{
		"web/views/login/login.html",
	}...)

	if err != nil {
		log.Panicln(err)
	}

	otptpl, err = templates.Build("otp.html").ParseFiles(
		"web/views/login/otp.html",
	)

	if err != nil {
		log.Panicln(err)
	}
}

func Form(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	isHX, _ := ctx.Value(contexts.HX).(bool)

	_, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		slog.LogAttrs(ctx, slog.LevelInfo, "the user is already connected")

		if strings.HasSuffix(r.Header.Get("HX-Current-URL"), "sso.html") {
			httpext.Redirect(w, r, "/admin/index.html", http.StatusFound)
		} else {
			httpext.Redirect(w, r, "/account/index.html", http.StatusFound)
		}

		return
	}

	var t *template.Template

	if isHX {
		t = hxtpl
	} else {
		if r.URL.Path == "/sso.html" {
			t = tpl
		} else {
			t = templates.Pages["login"]
		}
		policy := fmt.Sprintf("default-src 'self'; script-src-elem 'self' %s;", security.CSP["otp"])
		w.Header().Set("Content-Security-Policy", policy)
	}

	data := struct{ Lang language.Tag }{lang}
	if err := t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func Otp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		slog.LogAttrs(ctx, slog.LevelInfo, "the user is already connected")
		httpext.Redirect(w, r, "/admin/index.html", http.StatusFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		httperrors.Alert(w, ctx, err.Error())
		return
	}

	email := r.FormValue("email")
	if strings.HasSuffix(r.Header.Get("HX-Current-URL"), "sso.html") && !users.IsAdmin(ctx, email) {
		httperrors.InputMessage(w, ctx, "input:email")
		return
	}

	glue, err := users.Otp(ctx, email)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)

	cancelURL := "/otp.html"
	if strings.HasSuffix(r.Header.Get("HX-Current-URL"), "sso.html") {
		cancelURL = "/sso.html"
	}

	data := struct {
		Lang      language.Tag
		Glue      string
		Email     string
		CancelURL string
	}{lang, glue, email, cancelURL}

	if err = otptpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		slog.LogAttrs(ctx, slog.LevelInfo, "the user is already connected")
		if strings.HasSuffix(r.Header.Get("HX-Current-URL"), "sso.html") {
			httpext.Redirect(w, r, "/admin/index.html", http.StatusFound)
		} else {
			httpext.Redirect(w, r, "/account/index.html", http.StatusFound)
		}

		return
	}

	err := r.ParseForm()
	if err != nil {
		httperrors.Alert(w, ctx, err.Error())
		return
	}

	otp := strings.Join(r.Form["otp"], "")
	glue := r.FormValue("glue")
	device := r.Header.Get("User-Agent")

	sid, err := users.Login(ctx, otp, glue, device)
	if err != nil || sid == "" {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     cookies.SessionID,
		Value:    sid,
		MaxAge:   int(conf.Cookie.MaxAge),
		Path:     "/",
		HttpOnly: true,
		Secure:   conf.Cookie.Secure,
		Domain:   conf.Cookie.Domain,
	}

	http.SetCookie(w, cookie)

	if strings.HasSuffix(r.Header.Get("HX-Current-URL"), "sso.html") {
		httpext.Redirect(w, r, "/admin/index.html", http.StatusFound)
	} else {
		httpext.Redirect(w, r, "/account/index.html", http.StatusFound)
	}

	w.WriteHeader(http.StatusOK)
}
