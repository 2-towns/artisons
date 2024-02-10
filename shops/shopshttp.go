package shops

import (
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/forms"
	"artisons/http/httperrors"
	"artisons/http/httphelpers"
	"artisons/templates"
	"context"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
)

var settingsTpl *template.Template
var settingsAlertTpl *template.Template

func init() {
	var err error

	settingsTpl, err = templates.Build("base.html").ParseFiles(
		append(templates.AdminUI,
			conf.WorkingSpace+"web/views/admin/icons/close.svg",
			conf.WorkingSpace+"web/views/admin/settings/settings.html",
		)...)

	if err != nil {
		log.Panicln(err)
	}

	settingsAlertTpl, err = templates.Build("alert-success.html").ParseFiles(templates.AdminSuccess...)

	if err != nil {
		log.Panicln(err)
	}
}

func SettingsFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := httphelpers.Form[Settings]{
		Data:     Data,
		Lang:     lang,
		Currency: conf.Currency,
		Page:     "Settings",
	}

	if err := settingsTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func SettingsShopSave(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	s := ShopSettings{
		GmapKey:          r.FormValue("gmap_key"),
		Color:            r.FormValue("color"),
		Active:           r.FormValue("active") == "on",
		Cache:            r.FormValue("cache") == "on",
		Guest:            r.FormValue("guest") == "on",
		Quantity:         r.FormValue("quantity") == "on",
		New:              r.FormValue("new") == "on",
		Redirect:         r.FormValue("redirect") == "on",
		FuzzySearch:      r.FormValue("fuzzy_search") == "on",
		ExactMatchSearch: r.FormValue("exact_match_search") == "on",
	}

	width := r.FormValue("image_width")
	if width != "" {
		val, err := strconv.ParseInt(r.FormValue("image_width"), 10, 64)
		if err != nil {
			ctx = context.WithValue(ctx, contexts.HXTarget, "#alert-shop")
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the image width", slog.String("image_width", width), slog.String("error", err.Error()))
			httperrors.HXCatch(w, ctx, "input:image_width")
			return
		}

		s.ImageWidth = int(val)
	}

	height := r.FormValue("image_height")
	if height != "" {
		val, err := strconv.ParseInt(r.FormValue("image_height"), 10, 64)
		if err != nil {
			ctx = context.WithValue(ctx, contexts.HXTarget, "#alert-shop")
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the image height", slog.String("image_height", height), slog.String("error", err.Error()))
			httperrors.HXCatch(w, ctx, "input:image_height")
			return
		}

		s.ImageHeight = int(val)
	}

	items := r.FormValue("items")
	if r.FormValue("items") != "" {
		val, err := strconv.ParseInt(items, 10, 64)
		if err != nil {
			ctx = context.WithValue(ctx, contexts.HXTarget, "#alert-shop")
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
			ctx = context.WithValue(ctx, contexts.HXTarget, "#alert-shop")
			slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the min value", slog.String("min", min))
			httperrors.HXCatch(w, ctx, "input:min")
			return
		}

		s.Min = int(val)
	}

	err := s.Validate(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	_, err = s.Save(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	w.Header().Set("HX-Reswap", "innerHTML show:#alert:top")

	lang := ctx.Value(contexts.Locale).(language.Tag)
	rid, _ := ctx.Value(contexts.RequestID).(string)

	data := struct {
		Flash string
		Lang  language.Tag
		RID   string
	}{
		"The data has been saved successfully.",
		lang,
		rid,
	}

	if err := settingsAlertTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func SettingsContactSave(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(conf.MaxUploadSize); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	s := Contact{
		Name:    r.FormValue("name"),
		Address: r.FormValue("address"),
		City:    r.FormValue("city"),
		Zipcode: r.FormValue("zipcode"),
		Phone:   r.FormValue("phone"),
		Email:   r.FormValue("email"),
	}

	if r.FormValue("banner_1_delete") != "" {
		s.Banner1 = "-"
	}

	if r.FormValue("banner_2_delete") != "" {
		s.Banner2 = "-"
	}

	if r.FormValue("banner_3_delete") != "" {
		s.Banner3 = "-"
	}

	err := s.Validate(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	images := []string{"logo", "banner_1", "banner_2", "banner_3"}
	files, err := forms.Upload(r, "blog", images)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if Data.Logo == "" && files[0] == "" {
		slog.LogAttrs(ctx, slog.LevelError, "cannot process the empty logo")
		httperrors.HXCatch(w, ctx, "input:logo")
		return
	}

	if files[1] != "" {
		s.Banner1 = files[1]
	}

	if files[2] != "" {
		s.Banner2 = files[2]
	}

	if files[3] != "" {
		s.Banner3 = files[3]
	}

	_, err = s.Save(ctx)
	if err != nil {
		forms.RollbackUpload(ctx, files)
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	w.Header().Set("HX-Reswap", "innerHTML show:#alert:top")

	lang := ctx.Value(contexts.Locale).(language.Tag)
	rid, _ := ctx.Value(contexts.RequestID).(string)

	data := struct {
		Flash string
		Lang  language.Tag
		RID   string
	}{
		"The data has been saved successfully.",
		lang,
		rid,
	}

	if err := settingsAlertTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
