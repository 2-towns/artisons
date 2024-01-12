package admin

import (
	"gifthub/conf"
	"gifthub/http/contexts"
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
		append(templates.AdminUI,
			conf.WorkingSpace+"web/views/admin/icons/close.svg",
			conf.WorkingSpace+"web/views/admin/locales/locales.html",
			conf.WorkingSpace+"web/views/admin/settings/settings.html",
		)...)

	if err != nil {
		log.Panicln(err)
	}
}

func SettingsForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	data := struct {
		Lang    language.Tag
		Page    string
		Data    shops.Settings
		Locales []language.Tag
	}{
		lang,
		"settings",
		shops.Data,
		conf.LocalesSupported,
	}

	if err := settingsTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
