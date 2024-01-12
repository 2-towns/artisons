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
var blogCspPolicy = ""

func init() {
	var err error

	blogAddTpl, err = templates.Build("base.html").ParseFiles(
		append(templates.AdminUI,
			conf.WorkingSpace+"web/views/admin/blog/blog-head.html",
			conf.WorkingSpace+"web/views/admin/blog/blog-scripts.html",
			conf.WorkingSpace+"web/views/admin/blog/blog-add.html",
			conf.WorkingSpace+"web/views/admin/blog/blog-form.html",
		)...)

	if err != nil {
		log.Panicln(err)
	}

	blogCspPolicy = "default-src 'self'"
	blogCspPolicy += " https://maxcdn.bootstrapcdn.com/font-awesome/latest/css/font-awesome.min.css"
	blogCspPolicy += " https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.eot"
	blogCspPolicy += " https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.woff2"
	blogCspPolicy += " https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.ttf"
	blogCspPolicy += " https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.svg"
	blogCspPolicy += " https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.woff"
}

func AddBlogForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	data := struct {
		Lang    language.Tag
		Page    string
		ID      string
		Data    blogs.Article
		Locales []language.Tag
	}{
		lang,
		"blog",
		"",
		blogs.Article{},
		conf.LocalesSupported,
	}

	w.Header().Set("Content-Security-Policy", blogCspPolicy)

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
