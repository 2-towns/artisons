package admin

import (
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/httperrors"
	"gifthub/stats"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

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
		conf.WorkingSpace + "web/views/admin/icons/building-store.svg",
	}

	dashboardTpl, err = templates.Build("base.html").ParseFiles(append([]string{
		conf.WorkingSpace + "web/views/admin/base.html",
		conf.WorkingSpace + "web/views/admin/ui.html",
		conf.WorkingSpace + "web/views/admin/icons/home.svg",
		conf.WorkingSpace + "web/views/admin/dashboard/dashboard-actions.html",
		conf.WorkingSpace + "web/views/admin/dashboard/dashboard-head.html",
		conf.WorkingSpace + "web/views/admin/dashboard/dashboard-scripts.html",
	}, files...)...)

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
	Data []stats.MostValue
	Lang language.Tag
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
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

	mvs, err := stats.MostValues(ctx, days)
	if err != nil {
		httperrors.Catch(w, ctx, "error_http_general", 500)
		return
	}

	all, err := stats.GetAll(ctx, days)
	if err != nil {
		httperrors.Catch(w, ctx, "error_http_general", 500)
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	isHX, _ := ctx.Value(contexts.HX).(bool)
	demo := ctx.Value(contexts.Demo).(bool)
	data := struct {
		Lang           language.Tag
		Page           string
		Data           stats.Data
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
		"dashboard",
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
