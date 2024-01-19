package admin

import (
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/cookies"
	"gifthub/http/httperrors"
	"gifthub/http/httpext"
	"gifthub/tags"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/text/language"
)

var tagsAddTpl *template.Template

func init() {
	var err error

	tagsAddTpl, err = templates.Build("base.html").ParseFiles(
		append(templates.AdminUI,
			conf.WorkingSpace+"web/views/admin/tags/tags-add.html",
			conf.WorkingSpace+"web/views/admin/tags/tags-form.html",
		)...)

	if err != nil {
		log.Panicln(err)
	}
}

func AddTagForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	data := struct {
		Lang language.Tag
		Page string
		ID   string
		Data tags.Tag
	}{
		lang,
		"Tags",
		"",
		tags.Tag{},
	}

	if err := tagsAddTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AddTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(conf.MaxUploadSize); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	id := ""
	t, err := processTagFrom(ctx, *r.MultipartForm, id)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	exists, err := tags.Exists(ctx, t.Key)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	if exists {
		httperrors.Alert(w, ctx, "the tag exists already")
		return
	}

	err = t.Save(ctx)
	if err != nil {
		httpext.RollbackUpload(ctx, []string{t.Image})
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     cookies.FlashMessage,
		Value:    "The data has been saved successfully.",
		MaxAge:   int(time.Minute.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   conf.Cookie.Secure,
		Domain:   conf.Cookie.Domain,
	}

	http.SetCookie(w, cookie)
	w.Header().Set("HX-Redirect", "/admin/tags.html")
	w.Write([]byte(""))
}
