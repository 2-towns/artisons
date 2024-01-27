package admin

import (
	"context"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/http/seo"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const seoName = "SEO"
const seoURL = "/admin/seo.html"

var seoTpl *template.Template
var seoHxTpl *template.Template
var seoFormTpl *template.Template

type seoFeature struct{}

func init() {
	var err error

	files := append(templates.AdminTable,
		conf.WorkingSpace+"web/views/admin/seo/seo-table.html",
	)

	seoTpl, err = templates.Build("base.html").ParseFiles(
		append(files, append(templates.AdminList,
			conf.WorkingSpace+"web/views/admin/seo/seo.html",
		)...)...)

	if err != nil {
		log.Panicln(err)
	}

	seoHxTpl, err = templates.Build("seo-table.html").ParseFiles(files...)

	if err != nil {
		log.Panicln(err)
	}

	seoFormTpl, err = templates.Build("base.html").ParseFiles(
		append(templates.AdminUI,
			conf.WorkingSpace+"web/views/admin/seo/seo-form.html",
		)...)

	if err != nil {
		log.Panicln(err)
	}
}

func (f seoFeature) ListTemplate(ctx context.Context) *template.Template {
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		return seoHxTpl
	}

	return seoTpl
}

func (f seoFeature) Search(ctx context.Context, q string, offset, num int) (searchResults[seo.Content], error) {

	res := seo.List(ctx, offset, num)

	return searchResults[seo.Content]{
		Total: res.Total,
		Items: res.Content,
	}, nil
}

func (f seoFeature) Find(ctx context.Context, id interface{}) (seo.Content, error) {
	return seo.Find(ctx, id.(string))
}

func (f seoFeature) ID(ctx context.Context, id string) (interface{}, error) {
	return id, nil
}

func (f seoFeature) Validate(ctx context.Context, r *http.Request, data seo.Content) error {
	return nil
}

func (data seoFeature) Digest(ctx context.Context, r *http.Request) (seo.Content, error) {
	key := chi.URLParam(r, "id")

	c := seo.Content{
		Key:         key,
		URL:         r.FormValue("url"),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	return c, nil
}

func (f seoFeature) IsImageRequired(a seo.Content, key string) bool {
	return false
}

func (f seoFeature) UpdateImage(a *seo.Content, key, image string) {
}

func SeoList(w http.ResponseWriter, r *http.Request) {
	digestList[seo.Content](w, r, list[seo.Content]{
		Name:    seoName,
		URL:     seoURL,
		Feature: seoFeature{},
	})
}

func SeoForm(w http.ResponseWriter, r *http.Request) {
	data, err := digestForm[seo.Content](w, r, Form[seo.Content]{
		Name:    seoName,
		Feature: seoFeature{},
	})

	if err != nil {
		return
	}

	if err := seoFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func SeoSave(w http.ResponseWriter, r *http.Request) {
	digestSave[seo.Content](w, r, save[seo.Content]{
		Name:    seoName,
		URL:     seoURL,
		Feature: seoFeature{},
		Form:    urlEncodedForm{},
		Images:  []string{},
		Folder:  "",
	})
}
