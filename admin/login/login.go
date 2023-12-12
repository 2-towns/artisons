package login

import (
	"fmt"
	"gifthub/admin/urls"
	"gifthub/cache"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/cookies"
	"gifthub/http/httperrors"
	"gifthub/http/httpext"
	"gifthub/http/security"
	"gifthub/locales"
	"gifthub/users"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, locales.TranslateError(err, lang), http.StatusInternalServerError)
		return
	}

	p := message.NewPrinter(lang)
	t := map[string]string{
		"seo_title":       p.Sprintf("seo_login_title"),
		"seo_description": p.Sprintf("seo_login_description"),
		"label_email":     p.Sprintf("input_label_email"),
		"label_signin":    p.Sprintf("label_login_signin"),
		"title_sub":       p.Sprintf("title_login_sub"),
		"message_signin":  p.Sprintf("message_login_signin"),
	}

	data := struct {
		T   map[string]string
		Url string
		CB  map[string]string
	}{
		t,
		urls.Otp,
		cache.Buster,
	}

	if !isHX {
		policy := fmt.Sprintf("default-src 'self'; script-src-elem 'self' %s;", security.CSP["otp"])
		w.Header().Set("Content-Security-Policy", policy)
	}

	tpl.Execute(w, &data)
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
		httperrors.Catch(w, ctx, err.Error())
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	tpl, err := template.ParseFiles(
		"web/views/admin/login/otp.html",
		"web/views/admin/js/otp.js.html",
	)

	if err != nil {
		httperrors.Alert(w, ctx, locales.TranslateError(err, lang))
		return
	}

	p := message.NewPrinter(lang)
	t := map[string]string{
		"label_otp":    p.Sprintf("input_label_otp"),
		"label_verify": p.Sprintf("label_login_verify"),
		"label_cancel": p.Sprintf("label_login_cancel"),
		"title_sub":    p.Sprintf("title_otp_sub"),
		"message_otp":  p.Sprintf("message_otp_login", email),
	}

	data := struct {
		T         map[string]string
		Url       string
		UrlCancel string
		Glue      string
	}{
		t,
		urls.Login,
		urls.AuthPrefix,
		glue,
	}

	tpl.Execute(w, &data)
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
		httperrors.Catch(w, ctx, err.Error())
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
	w.Header().Set("HX-Redirect", urls.Dashboard)
	w.WriteHeader(http.StatusOK)
}
