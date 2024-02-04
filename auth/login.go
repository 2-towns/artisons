package auth

import (
	"artisons/carts"
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/cookies"
	"artisons/http/httperrors"
	"artisons/http/security"
	"artisons/templates"
	"artisons/users"
	"context"
	"fmt"
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
		"web/views/login.html",
		"web/views/admin/icons/logo.svg",
	}...)

	if err != nil {
		log.Panicln(err)
	}

	hxtpl, err = templates.Build("login.html").ParseFiles([]string{
		"web/views/login.html",
	}...)

	if err != nil {
		log.Panicln(err)
	}

	otptpl, err = templates.Build("otp.html").ParseFiles(
		"web/views/otp.html",
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
			http.Redirect(w, r, "/admin/index.html", http.StatusFound)
		} else {
			http.Redirect(w, r, "/account/index.html", http.StatusFound)
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
		w.Header().Set("HX-Redirect", "/admin/index.html")
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

	err = users.Otp(ctx, email)
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
		Email     string
		CancelURL string
	}{lang, email, cancelURL}

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
			w.Header().Set("HX-Redirect", "/admin/index.html")
		} else {
			w.Header().Set("HX-Redirect", "/account/index.html")
		}

		return
	}

	err := r.ParseForm()
	if err != nil {
		httperrors.Alert(w, ctx, err.Error())
		return
	}

	otp := strings.Join(r.Form["otp"], "")
	email := r.FormValue("email")
	device := r.Header.Get("User-Agent")

	sid, uid, err := users.Login(ctx, email, otp, device)
	if err != nil || sid == "" {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	ctx = context.WithValue(ctx, contexts.UserID, uid)

	cookie := &http.Cookie{
		Name:     cookies.SessionID,
		Value:    sid,
		MaxAge:   int(conf.Cookie.MaxAge),
		Path:     "/",
		HttpOnly: true,
		Secure:   conf.Cookie.Secure,
		Domain:   conf.Cookie.Domain,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)

	if strings.HasSuffix(r.Header.Get("HX-Current-URL"), "sso.html") {
		w.Header().Set("HX-Redirect", "/admin/index.html")
	} else {
		cid, ok := ctx.Value(contexts.Cart).(string)
		if !ok || !carts.Exists(ctx, cid) {
			w.Header().Set("HX-Redirect", "/account/index.html")
			return
		}

		if err := carts.Merge(ctx); err != nil {
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}

		w.Header().Set("HX-Redirect", "/account/index.html")
	}
}
