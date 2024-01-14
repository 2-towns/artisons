package admin

import (
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/cookies"
	"gifthub/http/httperrors"
	"gifthub/http/seo"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

var seoEditTpl *template.Template

func init() {
	var err error

	seoEditTpl, err = templates.Build("base.html").ParseFiles(
		append(templates.AdminUI,
			conf.WorkingSpace+"web/views/admin/seo/seo-edit.html",
		)...)

	if err != nil {
		log.Panicln(err)
	}
}

func EditSeoForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	id := chi.URLParam(r, "id")

	c, err := seo.Find(ctx, id)
	if err != nil {
		httperrors.Page(w, ctx, "oops the data is not found", 404)
		return
	}

	data := struct {
		Lang language.Tag
		Page string
		ID   string
		Data seo.Content
	}{
		lang,
		"SEO",
		id,
		c,
	}

	if err := seoEditTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func EditSeo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

	key := chi.URLParam(r, "id")
	c := seo.Content{
		Key:         key,
		URL:         r.FormValue("url"),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	err := c.Validate(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	err = c.Save(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     cookies.FlashMessage,
		Value:    "text_seo_editsuccess",
		MaxAge:   int(time.Minute.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   conf.Cookie.Secure,
		Domain:   conf.Cookie.Domain,
	}

	http.SetCookie(w, cookie)
	w.Header().Set("HX-Redirect", "/admin/seo.html")
	w.Write([]byte(""))
}
