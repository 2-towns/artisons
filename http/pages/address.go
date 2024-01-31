// Package pages provides the application pages
package pages

import (
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/shops"
	"artisons/tags"
	"artisons/templates"
	"artisons/users"
	"log/slog"
	"net/http"

	"golang.org/x/text/language"
)

func AddressForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	user := ctx.Value(contexts.User).(users.User)

	data := struct {
		Lang language.Tag
		Shop shops.Settings
		Tags []tags.Leaf
		User users.User
	}{
		lang,
		shops.Data,
		tags.Tree,
		user,
	}

	if err := templates.Pages["address"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func Address(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	user := ctx.Value(contexts.User).(users.User)

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	a := users.Address{
		Firstname:     r.FormValue("firstname"),
		Lastname:      r.FormValue("lastname"),
		Street:        r.FormValue("street"),
		Complementary: r.FormValue("complementary"),
		City:          r.FormValue("city"),
		Zipcode:       r.FormValue("zipcode"),
		Phone:         r.FormValue("phone"),
	}

	if err := a.Validate(ctx); err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if err := a.Save(ctx, user.ID); err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	data := struct {
		Lang           language.Tag
		SuccessMessage string
	}{
		lang,
		"The address has been saved successfully.",
	}

	if err := templates.Pages["hx-success"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
