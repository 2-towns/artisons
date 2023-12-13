package admin

import (
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/httperrors"
	"gifthub/stats"
	"gifthub/templates"
	"log/slog"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
)

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

	files := []string{
		"web/views/admin/dashboard/dashboard.html",
		"web/views/admin/dashboard/table-top-values.html",
		"web/views/admin/dashboard/table-most-values.html",
		"web/views/admin/icons/anchor.svg",
		"web/views/admin/icons/building-store.svg",
	}

	if !isHX {
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

	tpl, err := templates.Build(lang, isHX).ParseFiles(files...)

	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	demo := ctx.Value(contexts.Demo).(bool)

	data := struct {
		Page           string
		Data           stats.Data
		Days           int
		Referers       []stats.MostValue
		Browsers       []stats.MostValue
		Systems        []stats.MostValue
		ProductsShared []stats.MostValue
		Visits         []stats.MostValue
		Products       []stats.MostValue
		Demo           bool
		Currency       string
	}{
		"dashboard",
		all,
		days,
		mvs[2],
		mvs[1],
		mvs[3],
		mvs[5],
		mvs[0],
		mvs[4],
		demo,
		conf.Currency,
	}

	if isHX {
		w.Header().Set("HX-Trigger-After-Settle", "ecm-dashboard-reload")
	}

	if err = tpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}

}
