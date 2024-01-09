package admin

import (
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/cookies"
	"gifthub/shops"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"

	"golang.org/x/text/language"
)

var settingsTpl *template.Template

func init() {
	var err error

	settingsTpl, err = templates.Build("base.html").ParseFiles(
		conf.WorkingSpace+"web/views/admin/base.html",
		conf.WorkingSpace+"web/views/admin/ui.html",
		conf.WorkingSpace+"web/views/admin/icons/home.svg",
		conf.WorkingSpace+"web/views/admin/icons/building-store.svg",
		conf.WorkingSpace+"web/views/admin/icons/receipt.svg",
		conf.WorkingSpace+"web/views/admin/icons/settings.svg",
		conf.WorkingSpace+"web/views/admin/icons/article.svg",
		conf.WorkingSpace+"web/views/admin/settings/settings.html",
	)

	if err != nil {
		log.Panicln(err)
	}
}

func SettingsForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	flash := ""
	c, err := r.Cookie(cookies.FlashMessage)
	if err != nil && c != nil {
		flash = c.Value
	}

	data := struct {
		Lang  language.Tag
		Page  string
		Flash string
		Data  shops.Settings
	}{
		lang,
		"settings",
		flash,
		shops.Data,
	}

	if err = settingsTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
