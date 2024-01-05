package admin

import (
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/cookies"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"

	"golang.org/x/text/language"
)

var settingsTpl *template.Template

type Settings struct {
	Name    string
	Address string
	City    string
	Zipcode string
	Phone   string
	Email   string
	Logo    string

	// Active defines if the store is available or not
	Active bool

	// Guest allows to accept guest order
	Guest bool

	// Show quantity in product page
	Quantity bool

	// Enable the stock managment
	Stock bool

	// Number of days during which the product is considered 'new'
	New bool

	// Max items per page
	Items int

	// Mininimum order
	Min int

	// Redirect after  the product was added to the cart
	Redirect bool

	// Display last products when the quantity is under the amount.
	// Set to zero to disable this feature.
	LastProducts int

	// AdvancedSearch enables the advanced search
	AdvancedSearch bool

	// Cache enables the advanced search
	Cache bool

	// Google map key used for geolocation api
	GmapKey string
}

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
		Data  Settings
	}{
		lang,
		"settings",
		flash,
		Settings{},
	}

	if err = settingsTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
