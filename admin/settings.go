package admin

import (
	"context"
	"errors"
	"gifthub/conf"
	"gifthub/http/contexts"
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

const settingsName = "Settings"
const settingsURL = "/admin/settings.html"
const settingsFolder = "settings"

var settingsTpl *template.Template
var settingsAlertTpl *template.Template

type settingsFeature struct{}
type settingsShopFeature struct{}
type settingsContactFeature struct{}

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

func (f settingsFeature) FormTemplate(ctx context.Context, w http.ResponseWriter) *template.Template {
	return settingsTpl
}

func (f settingsFeature) Find(ctx context.Context, id interface{}) (shops.Settings, error) {
	return shops.Data, nil
}

func (f settingsFeature) ID(ctx context.Context, id string) (interface{}, error) {
	return "settings", nil
}

func (data settingsShopFeature) Digest(ctx context.Context, r *http.Request) (shops.ShopSettings, error) {
	s := shops.ShopSettings{
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
			return shops.ShopSettings{}, errors.New("input:image_width")
		}

		s.ImageWidth = int(val)
	}

	height := r.FormValue("image_height")
	if height != "" {
		val, err := strconv.ParseInt(r.FormValue("image_height"), 10, 64)
		if err != nil {
			ctx = context.WithValue(ctx, contexts.HXTarget, "#alert-shop")
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the image height", slog.String("image_height", height), slog.String("error", err.Error()))
			return shops.ShopSettings{}, errors.New("input:image_height")
		}

		s.ImageHeight = int(val)
	}

	items := r.FormValue("items")
	if r.FormValue("items") != "" {
		val, err := strconv.ParseInt(items, 10, 64)
		if err != nil {
			ctx = context.WithValue(ctx, contexts.HXTarget, "#alert-shop")
			slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the items value", slog.String("items", items))
			return shops.ShopSettings{}, errors.New("input:items")

		}

		s.Items = int(val)
	}

	min := r.FormValue("min")
	if r.FormValue("min") != "" {
		val, err := strconv.ParseInt(min, 10, 64)
		if err != nil {
			ctx = context.WithValue(ctx, contexts.HXTarget, "#alert-shop")
			slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the min value", slog.String("min", min))
			return shops.ShopSettings{}, errors.New("input:min")
		}

		s.Min = int(val)
	}

	return s, nil
}

func (data settingsContactFeature) Digest(ctx context.Context, r *http.Request) (shops.Contact, error) {
	s := shops.Contact{
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

	return s, nil
}

func (f settingsContactFeature) IsImageRequired(s shops.Contact, key string) bool {
	return shops.Data.Logo == ""
}

func (f settingsContactFeature) UpdateImage(s *shops.Contact, key, image string) {
	switch key {
	case "logo":
		s.Logo = image
	case "banner_1":
		s.Banner1 = image
	case "banner_2":
		s.Banner2 = image
	case "banner_3":
		s.Banner3 = image
	}
}

func (f settingsShopFeature) IsImageRequired(a shops.ShopSettings, key string) bool {
	return false
}

func (f settingsShopFeature) UpdateImage(a *shops.ShopSettings, key, image string) {

}

func SettingsForm(w http.ResponseWriter, r *http.Request) {
	httpext.DigestForm[shops.Settings](w, r, httpext.Form[shops.Settings]{
		Feature: settingsFeature{},
	})
}

func SettingsShopSave(w http.ResponseWriter, r *http.Request) {
	httpext.DigestSave[shops.ShopSettings](w, r, httpext.Save[shops.ShopSettings]{
		Name:       settingsName,
		URL:        settingsURL,
		Feature:    settingsShopFeature{},
		Form:       httpext.UrlEncodedForm{},
		Images:     []string{},
		Folder:     "",
		NoRedirect: true,
	})

	w.Header().Set("HX-Reswap", "innerHTML show:#alert:top")

	ctx := r.Context()
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

	if err := settingsAlertTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func SettingsContactSave(w http.ResponseWriter, r *http.Request) {
	httpext.DigestSave[shops.Contact](w, r, httpext.Save[shops.Contact]{
		Name:       settingsName,
		URL:        settingsURL,
		Feature:    settingsContactFeature{},
		Form:       httpext.MultipartForm{},
		Images:     []string{"logo", "banner_1", "banner_2", "banner_3"},
		Folder:     settingsFolder,
		NoRedirect: true,
	})

	w.Header().Set("HX-Reswap", "innerHTML show:#alert:top")

	ctx := r.Context()
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

	if err := settingsAlertTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
