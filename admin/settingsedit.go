package admin

import (
	"context"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/httperrors"
	"gifthub/http/httpext"
	"gifthub/shops"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/text/language"
)

var settingsAlertTpl *template.Template

func init() {
	var err error

	settingsAlertTpl, err = templates.Build("alert-success.html").ParseFiles(templates.AdminSuccess...)

	if err != nil {
		log.Panicln(err)
	}
}

func EditShopSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, contexts.HXTarget, "#alert-shop")

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	s := shops.ShopSettings{
		GmapKey:        r.FormValue("gmap_key"),
		Active:         r.FormValue("active") == "on",
		AdvancedSearch: r.FormValue("advanced_search") == "on",
		Cache:          r.FormValue("cache") == "on",
		Guest:          r.FormValue("guest") == "on",
		Quantity:       r.FormValue("quantity") == "on",
		New:            r.FormValue("new") == "on",
		Redirect:       r.FormValue("redirect") == "on",
	}

	items := r.FormValue("items")
	if r.FormValue("items") != "" {
		val, err := strconv.ParseInt(items, 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the items value", slog.String("items", items))
			httperrors.HXCatch(w, ctx, "input:items")
			return
		}

		s.Items = int(val)
	}

	min := r.FormValue("min")
	if r.FormValue("min") != "" {
		val, err := strconv.ParseInt(min, 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the min value", slog.String("min", min))
			httperrors.HXCatch(w, ctx, "input:min")
			return
		}

		s.Min = int(val)
	}

	last := r.FormValue("last_products")
	if r.FormValue("min") != "" {
		val, err := strconv.ParseInt(last, 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot use the last value", slog.String("last", last))
			httperrors.HXCatch(w, ctx, "input:last")
			return
		}

		s.LastProducts = int(val)
	}

	err := s.Validate(ctx)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the shop", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	err = s.Save(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	rid, _ := ctx.Value(middleware.RequestIDKey).(string)

	data := struct {
		Flash string
		Lang  language.Tag
		RID   string
	}{
		"The data has been saved successfully.",
		lang,
		rid,
	}

	w.Header().Set("HX-Reswap", "outerHTML show:#alert-shop:top")

	if err := settingsAlertTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func EditContactSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, contexts.HXTarget, "#alert-contact")

	if err := r.ParseMultipartForm(conf.MaxUploadSize); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	s := shops.Contact{
		Name:    r.FormValue("name"),
		Address: r.FormValue("address"),
		City:    r.FormValue("city"),
		Zipcode: r.FormValue("zipcode"),
		Phone:   r.FormValue("phone"),
		Email:   r.FormValue("email"),
	}

	err := s.Validate(ctx)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the shop", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	form := r.MultipartForm.File
	files, err := httpext.ProcessFiles(ctx, form, []string{"logo", "banner_1", "banner_2", "banner_3"})
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the shop", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if files["logo"] == nil && shops.Data.Logo == "" {
		slog.LogAttrs(ctx, slog.LevelInfo, "the logo is required")
		httperrors.HXCatch(w, ctx, "input:logo")
		return
	}

	images, err := httpext.Upload(ctx, files)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot update the files", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	if images["logo"] != "" {
		s.Logo = images["logo"]
	}

	del1 := r.FormValue("banner_1_delete")
	if images["banner_1"] != "" {
		s.Banner1 = images["banner_1"]
	} else if len(del1) > 0 && del1 != "" {
		s.Banner1 = ""
	}

	del2 := r.FormValue("banner_2_delete")
	if images["banner_1"] != "" {
		s.Banner2 = images["banner_2"]
	} else if len(del2) > 0 && del2 != "" {
		s.Banner2 = ""
	}

	del3 := r.FormValue("banner_3_delete")
	if images["banner_3"] != "" {
		s.Banner3 = images["banner_3"]
	} else if len(del3) > 0 && del3 != "" {
		s.Banner3 = ""
	}

	err = s.Save(ctx)
	if err != nil {
		httpext.RollbackUpload(ctx, []string{s.Logo, s.Banner1, s.Banner2, s.Banner3})
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	rid, _ := ctx.Value(middleware.RequestIDKey).(string)

	data := struct {
		Flash string
		Lang  language.Tag
		RID   string
	}{
		"The data has been saved successfully.",
		lang,
		rid,
	}

	w.Header().Set("HX-Reswap", "innerHTML show:#alert-contact:top")

	if err := settingsAlertTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
