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

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

var tagsEditTpl *template.Template

func init() {
	var err error

	tagsEditTpl, err = templates.Build("base.html").ParseFiles(
		append(templates.AdminUI,
			conf.WorkingSpace+"web/views/admin/tags/tags-edit.html",
			conf.WorkingSpace+"web/views/admin/tags/tags-form.html",
		)...)

	if err != nil {
		log.Panicln(err)
	}

}

func EditTagForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	id := chi.URLParam(r, "id")

	t, err := tags.Find(ctx, id)
	if err != nil {
		httperrors.Page(w, ctx, err.Error(), 400)
		return
	}

	data := struct {
		Lang language.Tag
		Page string
		ID   string
		Data tags.Tag
	}{
		lang,
		"Tags",
		id,
		t,
	}

	if err := tagsEditTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func EditTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(conf.MaxUploadSize); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	id := chi.URLParam(r, "id")
	t, err := processTagFrom(ctx, *r.MultipartForm, id)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
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
