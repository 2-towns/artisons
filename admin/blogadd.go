package admin

import (
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
	"time"

	"golang.org/x/text/language"
)

var blogAddTpl *template.Template

func init() {
	var err error

	blogAddTpl, err = templates.Build("base.html").ParseFiles([]string{
		conf.WorkingSpace + "web/views/admin/base.html",
		conf.WorkingSpace + "web/views/admin/ui.html",
		conf.WorkingSpace + "web/views/admin/icons/home.svg",
		conf.WorkingSpace + "web/views/admin/icons/building-store.svg",
		conf.WorkingSpace + "web/views/admin/icons/receipt.svg",
		conf.WorkingSpace + "web/views/admin/icons/settings.svg",
		conf.WorkingSpace + "web/views/admin/icons/article.svg",
		conf.WorkingSpace + "web/views/admin/icons/close.svg",
		conf.WorkingSpace + "web/views/admin/blog/blog-add.html",
		conf.WorkingSpace + "web/views/admin/blog/blog-form.html",
	}...)

	if err != nil {
		log.Panicln(err)
	}
}

func AddBlogForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	data := struct {
		Lang language.Tag
		Page string
		ID   string
		Data blogs.Article
	}{
		lang,
		"blog",
		"",
		blogs.Article{},
	}

	if err := blogAddTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AddBlog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(conf.MaxUploadSize); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "error_http_general")
		return
	}

	id := ""
	a, err := processBlogFrom(ctx, *r.MultipartForm, id)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	err = a.Save(ctx)
	if err != nil {
		httpext.RollbackUpload(ctx, []string{a.Image})
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     cookies.FlashMessage,
		Value:    "text_blog_addsuccess",
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
