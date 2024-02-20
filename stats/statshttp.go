package stats

import (
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

	"github.com/google/uuid"
	"github.com/mileusna/useragent"
	"golang.org/x/text/language"
)

var dashboardTpl *template.Template
var dashboardHxTpl *template.Template

func init() {
	var err error

	files := []string{
		conf.WorkingSpace + "web/views/admin/dashboard/dashboard.html",
		conf.WorkingSpace + "web/views/admin/dashboard/table-top-values.html",
		conf.WorkingSpace + "web/views/admin/dashboard/table-most-values.html",
		conf.WorkingSpace + "web/views/admin/icons/anchor.svg",
	}

	dashboardTpl, err = templates.Build("base.html").ParseFiles(append(templates.AdminUI,
		append(files,
			conf.WorkingSpace+"web/views/admin/dashboard/dashboard-actions.html",
			conf.WorkingSpace+"web/views/admin/dashboard/dashboard-head.html",
			conf.WorkingSpace+"web/views/admin/dashboard/dashboard-scripts.html")...)...)

	if err != nil {
		log.Panicln(err)
	}

	dashboardHxTpl, err = templates.Build("dashboard-hx.html").ParseFiles(
		append([]string{conf.WorkingSpace + "web/views/admin/dashboard/dashboard-hx.html"}, files...)...,
	)

	if err != nil {
		log.Panicln(err)
	}
}

type table struct {
	Data []MostValue
	Lang language.Tag
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var days int = 7

	pdays := r.URL.Query().Get("days")

	if pdays != "" {
		d, err := strconv.ParseInt(pdays, 10, 32)
		if err == nil {
			days = int(d)
		}

		if d > 30 {
			days = 30
		}
	}

	mvs, err := MostValues(ctx, days)
	if err != nil {
		httperrors.Catch(w, ctx, "something went wrong", 500)
		return
	}

	all, err := GetAll(ctx, days)
	if err != nil {
		httperrors.Catch(w, ctx, "something went wrong", 500)
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	isHX, _ := ctx.Value(contexts.HX).(bool)
	u, ok := ctx.Value(contexts.User).(users.User)
	demo := ok && u.Demo
	data := struct {
		Lang           language.Tag
		Page           string
		Data           Data
		Days           int
		Referers       table
		Browsers       table
		Systems        table
		ProductsShared table
		Visits         table
		Products       table
		Demo           bool
		Currency       string
	}{
		lang,
		"Dashboard",
		all,
		days,
		table{
			Data: mvs[2],
			Lang: lang,
		},
		table{
			Data: mvs[1],
			Lang: lang,
		},
		table{
			Data: mvs[3],
			Lang: lang,
		},
		table{
			Data: mvs[5],
			Lang: lang,
		},
		table{
			Data: mvs[0],
			Lang: lang,
		},
		table{
			Data: mvs[4],
			Lang: lang,
		},
		demo,
		conf.Currency,
	}

	var t *template.Template

	if isHX {
		w.Header().Set("HX-Trigger-After-Settle", "ecm-dashboard-reload")
		t = dashboardHxTpl
	} else {
		t = dashboardTpl
	}

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func DemoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(contexts.User).(users.User)
	if !ok {
		httperrors.Catch(w, ctx, "something_went_wrong", 400)
		return
	}

	_, err := user.ToggleDemo(ctx)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	w.Header().Set("HX-Redirect", "/admin/index")
	w.Write([]byte(""))
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var did string
		val, err := r.Cookie(cookies.Device)
		if err != nil || val.Value == "" {
			id := uuid.New()
			did = id.String()
		} else {
			did = val.Value
		}

		c := httphelpers.NewCookie(cookies.Device, did, int(conf.Cookie.MaxAge))
		http.SetCookie(w, &c)

		ctx := context.WithValue(r.Context(), contexts.Device, did)

		ua := useragent.Parse(r.Header.Get("User-Agent"))
		go Visit(ctx, ua, VisitData{
			URL:     r.URL.Path,
			Referer: r.Referer(),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
