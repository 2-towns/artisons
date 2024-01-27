package admin

import (
	"context"
	"errors"
	"gifthub/blog"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/contexts"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const blogName = "CMS"
const blogURL = "/admin/blog.html"
const blogFolder = "blog"

var blogTpl *template.Template
var blogHxTpl *template.Template
var blogFormTpl *template.Template
var blogCspPolicy = ""

type blogFeature struct{}

func init() {
	var err error

	files := append(templates.AdminTable,
		conf.WorkingSpace+"web/views/admin/blog/blog-table.html",
	)

	blogTpl, err = templates.Build("base.html").ParseFiles(
		append(files, append(templates.AdminList,
			conf.WorkingSpace+"web/views/admin/blog/blog-actions.html",
			conf.WorkingSpace+"web/views/admin/blog/blog.html",
		)...)...)

	if err != nil {
		log.Panicln(err)
	}

	blogHxTpl, err = templates.Build("blog-table.html").ParseFiles(files...)

	if err != nil {
		log.Panicln(err)
	}

	blogFormTpl, err = templates.Build("base.html").ParseFiles(
		append(templates.AdminUI,
			conf.WorkingSpace+"web/views/admin/blog/blog-head.html",
			conf.WorkingSpace+"web/views/admin/blog/blog-scripts.html",
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

func (f blogFeature) ListTemplate(ctx context.Context) *template.Template {
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		return blogHxTpl
	}

	return blogTpl
}

func (f blogFeature) Search(ctx context.Context, q string, offset, num int) (searchResults[blog.Article], error) {
	query := blog.Query{}
	if q != "" {
		query.Keywords = db.Escape(q)
	}

	res, err := blog.Search(ctx, query, offset, num)

	return searchResults[blog.Article]{
		Total: res.Total,
		Items: res.Articles,
	}, err
}

func (data blogFeature) Digest(ctx context.Context, r *http.Request) (blog.Article, error) {
	status := "online"

	if r.FormValue("status") != "on" {
		status = "offline"
	}

	a := blog.Article{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Status:      status,
		Type:        "blog",
	}

	id := chi.URLParam(r, "id")
	if id != "" {
		i, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.String("id", id), slog.String("error", err.Error()))
			return blog.Article{}, errors.New("input:id")
		}

		a.ID = int(i)
	}

	return a, nil
}

func (f blogFeature) ID(ctx context.Context, id string) (interface{}, error) {
	if id == "" {
		return 0, nil
	}

	val, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the id", slog.String("id", id), slog.String("error", err.Error()))
		return 0, errors.New("oops the data is not found")
	}

	return int(val), nil
}

func (f blogFeature) Find(ctx context.Context, id interface{}) (blog.Article, error) {
	return blog.Find(ctx, id.(int))
}

func (f blogFeature) Delete(ctx context.Context, id interface{}) error {
	deletable, err := blog.Deletable(ctx, id.(int))
	if err != nil {
		return err
	}

	if !deletable {
		return errors.New("the item cannot be deleted")
	}

	return blog.Delete(ctx, id.(int))
}

func (f blogFeature) IsImageRequired(a blog.Article, key string) bool {
	return a.ID == 0
}

func (f blogFeature) UpdateImage(a *blog.Article, key, image string) {
	a.Image = image
}

func (f blogFeature) Validate(ctx context.Context, r *http.Request, data blog.Article) error {
	return nil
}

func BlogSave(w http.ResponseWriter, r *http.Request) {
	digestSave[blog.Article](w, r, save[blog.Article]{
		Name:    blogName,
		URL:     blogURL,
		Feature: blogFeature{},
		Form:    multipartForm{},
		Images:  []string{"image"},
		Folder:  blogFolder,
	})
}

func BlogList(w http.ResponseWriter, r *http.Request) {
	digestList[blog.Article](w, r, list[blog.Article]{
		Name:    blogName,
		URL:     blogURL,
		Feature: blogFeature{},
	})
}

func BlogForm(w http.ResponseWriter, r *http.Request) {
	data, err := digestForm[blog.Article](w, r, Form[blog.Article]{
		Name:    blogName,
		Feature: blogFeature{},
	})

	if err != nil {
		return
	}

	w.Header().Set("Content-Security-Policy", blogCspPolicy)

	if err := blogFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func BlogDelete(w http.ResponseWriter, r *http.Request) {
	digestDelete[blog.Article](w, r, delete[blog.Article]{
		list: list[blog.Article]{
			Name:    blogName,
			URL:     blogURL,
			Feature: blogFeature{},
		},
		Feature: blogFeature{},
	})
}
