package locales

// var localesTpl *template.Template

// func init() {
// 	var err error

// 	localesTpl, err = templates.Build("locales-success.html").ParseFiles(
// 		append(templates.AdminSuccess,
// 			conf.WorkingSpace+"web/views/admin/locales/locales-success.html",
// 			conf.WorkingSpace+"web/views/admin/locales/locales.html",
// 		)...,
// 	)

// 	if err != nil {
// 		log.Panicln(err)
// 	}
// }

// func EditLocale(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	v := Value{
// 		Key:   r.FormValue("key"),
// 		Value: r.FormValue("value"),
// 	}

// 	err := v.Validate(ctx)
// 	if err != nil {
// 		httperrors.HXCatch(w, ctx, err.Error())
// 		return
// 	}

// 	err = v.Save(ctx)
// 	if err != nil {
// 		httperrors.HXCatch(w, ctx, err.Error())
// 		return
// 	}

// 	lang := ctx.Value(contexts.Locale).(language.Tag)

// 	data := struct {
// 		Flash string
// 		Lang  language.Tag
// 	}{
// 		"The data has been saved successfully.",
// 		lang,
// 	}

// 	w.Header().Set("HX-Reswap", "outerHTML show:#alert:top")

// 	if err := localesTpl.Execute(w, &data); err != nil {
// 		slog.Error("cannot render the template", slog.String("error", err.Error()))
// 	}
// }
