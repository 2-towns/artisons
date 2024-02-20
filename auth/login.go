package auth

import (
	"artisons/carts"
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/cookies"
	"artisons/http/httperrors"
	"artisons/http/httphelpers"
	"artisons/templates"
	"artisons/users"
	"context"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/language"
)

var tpl *template.Template
var otptpl *template.Template

func init() {
	var err error

	tpl, err = templates.Build("base.html").ParseFiles([]string{
		conf.WorkingSpace + "web/views/admin/base.html",
		conf.WorkingSpace + "web/views/admin/simple.html",
		conf.WorkingSpace + "web/views/login.html",
		conf.WorkingSpace + "web/views/admin/icons/logo.svg",
	}...)

	if err != nil {
		log.Panicln(err)
	}

	otptpl, err = templates.Build("otp.html").ParseFiles(
		conf.WorkingSpace + "web/views/otp.html",
	)

	if err != nil {
		log.Panicln(err)
	}
}

func Formhandler(w http.ResponseWriter, r *http.Request) {
	ctx := users.Context(r, w)
	_, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		slog.LogAttrs(ctx, slog.LevelInfo, "the user is already connected")

		if r.URL.Path == "/sso" {
			http.Redirect(w, r, "/admin/index", http.StatusFound)
		} else {
			http.Redirect(w, r, "/account/index", http.StatusFound)
		}

		return
	}

	var t *template.Template

	if r.URL.Path == "/sso" {
		t = tpl
	} else {
		t = templates.Pages["login"]
	}

	w.Header().Add("Content-Type", "text/html")

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := struct{ Lang language.Tag }{lang}
	if err := t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func OtpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := users.Context(r, w)
	_, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		slog.LogAttrs(ctx, slog.LevelInfo, "the user is already connected")

		if strings.HasSuffix(r.Header.Get("HX-Current-Url"), "/sso") {
			w.Header().Set("HX-Redirect", "/admin/index")
		} else {
			w.Header().Set("HX-Redirect", "/account/index")
		}

		return
	}

	email := r.FormValue("email")
	if strings.HasSuffix(r.Header.Get("HX-Current-Url"), "sso.html") && !users.IsAdmin(ctx, email) {
		httperrors.InputMessage(w, ctx, "input:email")
		return
	}

	err := users.Otp(ctx, email)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)

	cancelURL := "/otp"
	if strings.HasSuffix(r.Header.Get("HX-Current-Url"), "/sso") {
		cancelURL = "/sso"
	}

	w.Header().Add("HX-Trigger-After-Swap", "otp")

	data := struct {
		Lang      language.Tag
		Email     string
		CancelURL string
	}{lang, email, cancelURL}

	if err := otptpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := users.Context(r, w)
	_, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		slog.LogAttrs(ctx, slog.LevelInfo, "the user is already connected")
		if strings.HasSuffix(r.Header.Get("HX-Current-Url"), "/sso") {
			w.Header().Set("HX-Redirect", "/admin/index")
		} else {
			w.Header().Set("HX-Redirect", "/account/index")
		}

		return
	}

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	otp := strings.Join(r.Form["otp"], "")
	email := r.FormValue("email")
	device := r.Header.Get("User-Agent")

	u, err := users.Login(ctx, email, otp, device)
	if err != nil || u.SID == "" {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if strings.HasSuffix("/sso", r.Header.Get("HX-Current-Url")) && u.Role != "admin" {
		slog.LogAttrs(ctx, slog.LevelInfo, "the user tried to connect to admin")
		httperrors.HXCatch(w, ctx, "you are not authorized to process this request")
		return
	}

	ctx = context.WithValue(ctx, contexts.User, u)

	cookie := httphelpers.NewCookie(cookies.SessionID, u.SID, int(conf.Cookie.MaxAge))
	http.SetCookie(w, &cookie)

	if strings.HasSuffix(r.Header.Get("HX-Current-Url"), "/sso") && u.Role == "admin" {
		w.Header().Set("HX-Redirect", "/admin/index")
	} else {
		coo, err := r.Cookie(cookies.CartID)
		if err == nil {
			cid, err := strconv.ParseInt(coo.Value, 10, 64)
			if err == nil {
				if err := carts.Merge(ctx, int(cid)); err != nil {
					httperrors.HXCatch(w, ctx, err.Error())
					return
				}
			}
		}

		if r.Header.Get("HX-Current-Url") == "/cart" {
			w.Write([]byte(""))
		} else {
			w.Header().Set("HX-Redirect", "/account/index")
			w.Write([]byte(""))
		}
	}
}
