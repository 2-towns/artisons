package admin

import (
	"gifthub/conf"
	"gifthub/http/cookies"
	"gifthub/http/httperrors"
	"gifthub/http/httpext"
	"gifthub/shops"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func EditSettingsShop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "error_http_general")
		return
	}

	settings := shops.Settings{}

	items := r.FormValue("items")
	if r.FormValue("items") != "" {
		val, err := strconv.ParseInt(items, 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the items value", slog.String("items", items))
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}

		settings.Items = int(val)
	}

	min := r.FormValue("min")
	if r.FormValue("min") != "" {
		val, err := strconv.ParseInt(min, 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the items value", slog.String("min", min))
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}

		settings.Min = int(val)
	}

	last := r.FormValue("last_products")
	if r.FormValue("min") != "" {
		val, err := strconv.ParseInt(last, 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot use the items value", slog.String("last", last))
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}

		settings.LastProducts = int(val)
	}

	settings.Active = r.FormValue("active") == "on"
	settings.AdvancedSearch = r.FormValue("advanced_search") == "on"
	settings.Cache = r.FormValue("cache") == "on"
	settings.Guest = r.FormValue("guest") == "on"
	settings.Quantity = r.FormValue("quantity") == "on"
	settings.Stock = r.FormValue("stock") == "on"
	settings.New = r.FormValue("new") == "on"
	settings.Redirect = r.FormValue("redirect") == "on"

	err := settings.Validate(ctx)
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
		httperrors.HXCatch(w, ctx, "input_logo_required")
		return
	}

	images, err := httpext.Upload(ctx, files)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot update the files", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "error_http_general")
		return
	}

	if images["logo"] != "" {
		settings.Logo = images["logo"]
	}

	del1 := r.FormValue("banner_1_delete")
	if images["banner_1"] != "" {
		settings.Banner1 = images["banner_1"]
	} else if len(del1) > 0 && del1 != "" {
		settings.Banner1 = ""
	}

	del2 := r.FormValue("banner_2_delete")
	if images["banner_1"] != "" {
		settings.Banner2 = images["banner_2"]
	} else if len(del2) > 0 && del2 != "" {
		settings.Banner2 = ""
	}

	del3 := r.FormValue("banner_3_delete")
	if images["banner_3"] != "" {
		settings.Banner3 = images["banner_3"]
	} else if len(del3) > 0 && del3 != "" {
		settings.Banner3 = ""
	}

	err = settings.Save(ctx)
	if err != nil {
		httpext.RollbackUpload(ctx, []string{settings.Logo, settings.Banner1, settings.Banner2, settings.Banner3})
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     cookies.FlashMessage,
		Value:    "text_settings_editsuccess",
		MaxAge:   int(time.Minute.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   conf.Cookie.Secure,
		Domain:   conf.Cookie.Domain,
	}

	http.SetCookie(w, cookie)
	w.Header().Set("HX-Redirect", "/admin/settings.html")
	w.Write([]byte(""))
}

func EditSettingsContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(conf.MaxUploadSize); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "error_http_general")
		return
	}

	settings := shops.Settings{
		Name:    r.FormValue("name"),
		Address: r.FormValue("address"),
		City:    r.FormValue("city"),
		Zipcode: r.FormValue("zipcode"),
		Phone:   r.FormValue("phone"),
		Email:   r.FormValue("email"),
		GmapKey: r.FormValue("gmap_key"),
	}

	log.Println("name", settings)

	items := r.FormValue("items")
	if r.FormValue("items") != "" {
		val, err := strconv.ParseInt(items, 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the items value", slog.String("items", items))
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}

		settings.Items = int(val)
	}

	min := r.FormValue("min")
	if r.FormValue("min") != "" {
		val, err := strconv.ParseInt(min, 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the items value", slog.String("min", min))
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}

		settings.Min = int(val)
	}

	last := r.FormValue("last_products")
	if r.FormValue("min") != "" {
		val, err := strconv.ParseInt(last, 10, 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot use the items value", slog.String("last", last))
			httperrors.HXCatch(w, ctx, err.Error())
			return
		}

		settings.LastProducts = int(val)
	}

	settings.Active = r.FormValue("active") == "on"
	settings.AdvancedSearch = r.FormValue("advanced_search") == "on"
	settings.Cache = r.FormValue("cache") == "on"
	settings.Guest = r.FormValue("guest") == "on"
	settings.Quantity = r.FormValue("quantity") == "on"
	settings.Stock = r.FormValue("stock") == "on"
	settings.New = r.FormValue("new") == "on"
	settings.Redirect = r.FormValue("redirect") == "on"

	err := settings.Validate(ctx)
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
		httperrors.HXCatch(w, ctx, "input_logo_required")
		return
	}

	images, err := httpext.Upload(ctx, files)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot update the files", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "error_http_general")
		return
	}

	if images["logo"] != "" {
		settings.Logo = images["logo"]
	}

	del1 := r.FormValue("banner_1_delete")
	if images["banner_1"] != "" {
		settings.Banner1 = images["banner_1"]
	} else if len(del1) > 0 && del1 != "" {
		settings.Banner1 = ""
	}

	del2 := r.FormValue("banner_2_delete")
	if images["banner_1"] != "" {
		settings.Banner2 = images["banner_2"]
	} else if len(del2) > 0 && del2 != "" {
		settings.Banner2 = ""
	}

	del3 := r.FormValue("banner_3_delete")
	if images["banner_3"] != "" {
		settings.Banner3 = images["banner_3"]
	} else if len(del3) > 0 && del3 != "" {
		settings.Banner3 = ""
	}

	err = settings.Save(ctx)
	if err != nil {
		httpext.RollbackUpload(ctx, []string{settings.Logo, settings.Banner1, settings.Banner2, settings.Banner3})
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     cookies.FlashMessage,
		Value:    "text_settings_editsuccess",
		MaxAge:   int(time.Minute.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   conf.Cookie.Secure,
		Domain:   conf.Cookie.Domain,
	}

	http.SetCookie(w, cookie)
	w.Header().Set("HX-Redirect", "/admin/settings.html")
	w.Write([]byte(""))
}
