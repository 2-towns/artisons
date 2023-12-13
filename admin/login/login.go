package login

import (
	"fmt"
	"gifthub/admin/urls"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/cookies"
	"gifthub/http/httperrors"
	"gifthub/http/httpext"
	"gifthub/http/security"
	"gifthub/locales"
	"gifthub/templates"
	"gifthub/users"
	"log/slog"
	"net/http"
	"strings"

	"golang.org/x/text/language"
)

func Form(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	isHX, _ := ctx.Value(contexts.HX).(bool)

	_, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		slog.LogAttrs(ctx, slog.LevelInfo, "the user is already connected")
		httpext.Redirect(w, r, urls.AdminPrefix, http.StatusFound)
		return
	}

	files := []string{}

	if isHX {
		files = append(files,
			"web/views/admin/htmx.html",
			"web/views/admin/login/login.html",
		)
	} else {
		files = append(files,
			"web/views/admin/base.html",
			"web/views/admin/simple.html",
			"web/views/admin/login/login.html",
			"web/views/admin/icons/logo.svg",
		)
	}

	tpl, err := templates.Build(lang, isHX).ParseFiles(files...)
	if err != nil {
		http.Error(w, locales.TranslateError(err, lang), http.StatusInternalServerError)
		return
	}

	if !isHX {
		policy := fmt.Sprintf("default-src 'self'; script-src-elem 'self' %s;", security.CSP["otp"])
		w.Header().Set("Content-Security-Policy", policy)
	}

	data := struct{}{}
	if err = tpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func Otp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		slog.LogAttrs(ctx, slog.LevelInfo, "the user is already connected")
		httpext.Redirect(w, r, urls.AdminPrefix, http.StatusFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		httperrors.Alert(w, ctx, err.Error())
		return
	}

	email := r.FormValue("email")
	if !users.IsAdmin(ctx, email) {
		httperrors.InputMessage(w, ctx, "input_email_notadmin")
		return
	}

	glue, err := users.Otp(ctx, email)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	tpl, err := templates.Build(lang, true).ParseFiles(
		"web/views/admin/htmx.html",
		"web/views/admin/login/otp.html",
		"web/views/admin/js/otp.js.html",
	)

	if err != nil {
		httperrors.Alert(w, ctx, locales.TranslateError(err, lang))
		return
	}

	data := struct {
		Glue  string
		Email string
	}{glue, email}

	if err = tpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		slog.LogAttrs(ctx, slog.LevelInfo, "the user is already connected")
		httpext.Redirect(w, r, urls.AdminPrefix, http.StatusFound)
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
	// Todo check user role to redirect on correct page
	w.Header().Set("HX-Redirect", urls.Map["admin_dashboard"])
	w.WriteHeader(http.StatusOK)
}
