package admin

import (
	"gifthub/admin/urls"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/httperrors"
	"gifthub/stats"
	"html/template"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type tableData struct {
	T    map[string]string
	Data []stats.MostValue
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
	}

	mvs, err := stats.MostValues(ctx, days)
	if err != nil {
		httperrors.Catch(w, ctx, "error_http_general")
		return
	}

	all, err := stats.GetAll(ctx, days)
	if err != nil {
		httperrors.Catch(w, ctx, "error_http_general")
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)

	files := []string{
		"web/views/admin/dashboard/dashboard.html",
		"web/views/admin/dashboard/dashboard-demo-button.html",
		"web/views/admin/dashboard/table-top-values.html",
		"web/views/admin/dashboard/table-most-values.html",
		"web/views/admin/icons/anchor.svg",
	}

	if r.Header.Get("HX-Request") != "true" {
		files = append([]string{
			"web/views/admin/base.html",
			"web/views/admin/ui.html",
			"web/views/admin/icons/home.svg",
			"web/views/admin/dashboard/dashboard-actions.html",
			"web/views/admin/dashboard/dashboard-head.html",
			"web/views/admin/dashboard/dashboard-scripts.html",
		}, files...)
	} else {
		files = append([]string{"web/views/admin/htmx.html"}, files...)
	}

	tpl, err := template.ParseFiles(files...)

	if err != nil {
		if r.Header.Get("HX-Request") == "true" {
			httperrors.Catch(w, ctx, err.Error())
		} else {
			httperrors.Page(w, ctx, err.Error(), 500)
		}
		return
	}

	demo := ctx.Value(contexts.Demo).(bool)
	p := message.NewPrinter(lang)

	t := map[string]string{
		"seo_title":                p.Sprintf("seo_dashboard_title"),
		"seo_description":          "",
		"title_page":               p.Sprintf("title_dashboard_page"),
		"label_logout":             p.Sprintf("label_dashboard_logout"),
		"label_home":               p.Sprintf("label_dashboard_home"),
		"label_newusers":           p.Sprintf("label_dashboard_newusers"),
		"label_revenue":            p.Sprintf("label_dashboard_revenue"),
		"label_visits":             p.Sprintf("label_dashboard_visits"),
		"label_mostvisited":        p.Sprintf("label_dashboard_mostvisited"),
		"label_pageviews":          p.Sprintf("label_dashboard_pageviews"),
		"label_topreferers":        p.Sprintf("label_dashboard_topreferers"),
		"label_topsystems":         p.Sprintf("label_dashboard_topsystems"),
		"label_topbrowsers":        p.Sprintf("label_dashboard_topbrowsers"),
		"label_bouncerate":         p.Sprintf("label_dashboard_bouncerate"),
		"label_pagename":           p.Sprintf("label_dashboard_pagename"),
		"label_visitors":           p.Sprintf("label_dashboard_visitors"),
		"label_unique":             p.Sprintf("label_dashboard_unique"),
		"label_link":               p.Sprintf("label_dashboard_link"),
		"label_name":               p.Sprintf("label_dashboard_name"),
		"label_count":              p.Sprintf("label_dashboard_count"),
		"label_totalearning":       p.Sprintf("label_dashboard_totalearning"),
		"label_totalorders":        p.Sprintf("label_dashboard_totalorders"),
		"label_mostproductssold":   p.Sprintf("label_dashboard_mostproductssold"),
		"label_mostproductsshared": p.Sprintf("label_dashboard_mostproductsshared"),
		"label_7days":              p.Sprintf("label_dashboard_7days"),
		"label_14days":             p.Sprintf("label_dashboard_14days"),
		"label_30days":             p.Sprintf("label_dashboard_30days"),
		"label_demoactivate":       p.Sprintf("label_dashboard_demoactivate"),
		"label_demodisable":        p.Sprintf("label_dashboard_demodisable"),
	}

	urls := map[string]string{
		"home":   urls.AdminPrefix,
		"logout": urls.Logout,
		"demo":   urls.Demo,
	}

	css := map[string]string{
		"home": "header-menu-active",
	}

	data := struct {
		T              map[string]string
		Urls           map[string]string
		CSS            map[string]string
		Data           stats.Data
		Days           int
		Referers       tableData
		Browsers       tableData
		Systems        tableData
		ProductsShared tableData
		Visits         tableData
		Products       tableData
		Demo           bool
		Currency       string
	}{
		t,
		urls,
		css,
		all,
		days,
		tableData{
			T:    t,
			Data: mvs[2],
		},
		tableData{
			T:    t,
			Data: mvs[1],
		},
		tableData{
			T:    t,
			Data: mvs[3],
		},
		tableData{
			T:    t,
			Data: mvs[5],
		},
		tableData{
			T:    t,
			Data: mvs[0],
		},
		tableData{
			T:    t,
			Data: mvs[4],
		},
		demo,
		conf.Currency,
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Trigger-After-Settle", "ecm-dashboard-reload")
	}

	tpl.Execute(w, &data)

}
