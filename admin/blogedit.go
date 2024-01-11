package admin

import (
	"fmt"
	"gifthub/blogs"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/cookies"
	"gifthub/http/httperrors"
	"gifthub/http/httpext"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

var blogEditTpl *template.Template

func init() {
	var err error

	blogEditTpl, err = templates.Build("base.html").ParseFiles([]string{
		conf.WorkingSpace + "web/views/admin/base.html",
		conf.WorkingSpace + "web/views/admin/ui.html",
		conf.WorkingSpace + "web/views/admin/icons/home.svg",
		conf.WorkingSpace + "web/views/admin/icons/close.svg",
		conf.WorkingSpace + "web/views/admin/icons/building-store.svg",
		conf.WorkingSpace + "web/views/admin/icons/receipt.svg",
		conf.WorkingSpace + "web/views/admin/icons/settings.svg",
		conf.WorkingSpace + "web/views/admin/icons/article.svg",
		conf.WorkingSpace + "web/views/admin/icons/close.svg",
		conf.WorkingSpace + "web/views/admin/blog/blog-head.html",
		conf.WorkingSpace + "web/views/admin/blog/blog-scripts.html",
		conf.WorkingSpace + "web/views/admin/blog/blog-add.html",
		conf.WorkingSpace + "web/views/admin/blog/blog-form.html",
	}...)

	if err != nil {
		log.Panicln(err)
	}
}

func EditBlogForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	id := chi.URLParam(r, "id")

	iid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.String("id", id), slog.String("error", err.Error()))
		httperrors.Page(w, ctx, "error_http_blognotfound", 404)
	}

	p, err := blogs.Find(ctx, iid)
	if err != nil {
		httperrors.Page(w, ctx, "error_http_blognotfound", 404)
		return
	}

	data := struct {
		Lang language.Tag
		Page string
		ID   string
		Data blogs.Article
	}{
		lang,
		"blog",
		id,
		p,
	}

	policy := fmt.Sprintf("default-src 'self' https://unpkg.com/easymde/dist/easymde.min.js https://unpkg.com/easymde/dist/easymde.min.css https://maxcdn.bootstrapcdn.com/font-awesome/latest/css/font-awesome.min.css https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.eot  https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.woff2?v=4.7.0  https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.woff2?v=4.7.0 https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.ttf?v=4.7.0  https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.svg?v=4.7.0#fontawesomeregular https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.woff?v=4.7.0;")
	w.Header().Set("Content-Security-Policy", policy)

	if err := blogEditTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func EditBlog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(conf.MaxUploadSize); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "error_http_general")
		return
	}

	id := chi.URLParam(r, "id")
	p, err := processBlogFrom(ctx, *r.MultipartForm, id)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	err = p.Save(ctx)
	if err != nil {
		httpext.RollbackUpload(ctx, []string{p.Image})
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     cookies.FlashMessage,
		Value:    "text_blog_editsuccess",
		MaxAge:   int(time.Minute.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   conf.Cookie.Secure,
		Domain:   conf.Cookie.Domain,
	}

	http.SetCookie(w, cookie)
	w.Header().Set("HX-Redirect", "/admin/blog.html")
	w.Write([]byte(""))
}
