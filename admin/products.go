package admin

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/products"
	"artisons/products/filters"
	"artisons/string/stringutil"
	"artisons/tags"
	"artisons/templates"
	"context"
	"errors"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const productsName = "Products"
const productsURL = "/admin/products.html"
const productsFolder = "products"

var productsTpl *template.Template
var productsHxTpl *template.Template
var productsFormTpl *template.Template

type productsFeature struct{}

func init() {
	var err error

	files := append(templates.AdminTable,
		conf.WorkingSpace+"web/views/admin/products/products-table.html",
	)

	productsTpl, err = templates.Build("base.html").ParseFiles(
		append(files, append(templates.AdminList,
			conf.WorkingSpace+"web/views/admin/products/products-actions.html",
			conf.WorkingSpace+"web/views/admin/products/products.html")...,
		)...)

	if err != nil {
		log.Panicln(err)
	}

	productsHxTpl, err = templates.Build("products-table.html").ParseFiles(files...)

	if err != nil {
		log.Panicln(err)
	}

	productsFormTpl, err = templates.Build("base.html").ParseFiles(
		append(templates.AdminUI,
			conf.WorkingSpace+"web/views/admin/icons/close.svg",
			conf.WorkingSpace+"web/views/admin/products/products-head.html",
			conf.WorkingSpace+"web/views/admin/products/products-scripts.html",
			conf.WorkingSpace+"web/views/admin/slug.html",
			conf.WorkingSpace+"web/views/admin/products/products-form.html",
		)...)

	if err != nil {
		log.Panicln(err)
	}
}

func (f productsFeature) ListTemplate(ctx context.Context) *template.Template {
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		return productsHxTpl
	}

	return productsTpl
}

func (f productsFeature) Search(ctx context.Context, q string, offset, num int) (searchResults[products.Product], error) {
	query := products.Query{}
	if q != "" {
		query.Keywords = db.Escape(q)
	}

	res, err := products.Search(ctx, query, offset, num)

	return searchResults[products.Product]{
		Total: res.Total,
		Items: res.Products,
	}, err
}

func (data productsFeature) Digest(ctx context.Context, r *http.Request) (products.Product, error) {
	if r.FormValue("price") == "" {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the empty price")
		return products.Product{}, errors.New("input:price")
	}

	if r.FormValue("quantity") == "" {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the empty quantity")
		return products.Product{}, errors.New("input:quantity")
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the price", slog.String("price", r.FormValue("price")), slog.String("error", err.Error()))
		return products.Product{}, errors.New("input:price")
	}

	quantity, err := strconv.ParseInt(r.FormValue("quantity"), 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the quantity", slog.String("quantity", r.FormValue("quantity")), slog.String("error", err.Error()))
		return products.Product{}, errors.New("input:quantity")
	}

	var discount float64 = 0
	if r.FormValue("discount") != "" {
		val, err := strconv.ParseFloat(r.FormValue("discount"), 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the discount", slog.String("discount", r.FormValue("discount")), slog.String("error", err.Error()))
			return products.Product{}, errors.New("input:discount")
		}
		discount = val
	}

	var weight float64 = 0
	if r.FormValue("weight") != "" {
		val, err := strconv.ParseFloat(r.FormValue("weight"), 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the weight", slog.String("weight", r.FormValue("discount")), slog.String("error", err.Error()))
			return products.Product{}, errors.New("input:weight")
		}
		weight = val
	}

	status := "online"

	if r.FormValue("status") != "on" {
		status = "offline"
	}

	filters, err := filters.Actives(ctx)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the filters", slog.String("error", err.Error()))
		return products.Product{}, errors.New("something went wrong")
	}

	meta := map[string][]string{}
	for _, val := range filters {
		meta[val.Key] = r.Form[val.Key]
	}

	p := products.Product{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Status:      status,
		Sku:         r.FormValue("sku"),
		Tags:        r.MultipartForm.Value["tags"],
		Price:       price,
		Discount:    discount,
		Weight:      weight,
		Quantity:    int(quantity),
		Meta:        meta,
	}

	if r.FormValue("slug") != "" {
		p.Slug = r.FormValue("slug")
	} else {
		p.Slug = stringutil.Slugify(p.Title)
	}

	if r.FormValue("image_2_delete") != "" {
		p.Image2 = "-"
	}

	if r.FormValue("image_3_delete") != "" {
		p.Image3 = "-"
	}

	if r.FormValue("image_4_delete") != "" {
		p.Image4 = "-"
	}

	p.ID = chi.URLParam(r, "id")

	return p, nil
}

func (f productsFeature) ID(ctx context.Context, id string) (interface{}, error) {
	return id, nil
}

func (f productsFeature) Find(ctx context.Context, id interface{}) (products.Product, error) {
	return products.Find(ctx, id.(string))
}

func (f productsFeature) Delete(ctx context.Context, id interface{}) error {
	return products.Delete(ctx, id.(string))
}

func (f productsFeature) IsImageRequired(p products.Product, key string) bool {
	return p.ID == "" && key == "image_1"
}

func (f productsFeature) UpdateImage(p *products.Product, key, image string) {
	switch key {
	case "image_1":
		p.Image1 = image
	case "image_2":
		p.Image2 = image
	case "image_3":
		p.Image3 = image
	case "image_4":
		p.Image4 = image
	}
}

func (f productsFeature) Validate(ctx context.Context, r *http.Request, data products.Product) error {
	query := products.Query{Slug: data.Slug}
	res, err := products.Search(ctx, query, 0, 1)
	if err != nil || res.Total > 0 && (res.Products[0].ID != data.ID) {
		return errors.New("input:slug")
	}

	return nil
}

func ProductSave(w http.ResponseWriter, r *http.Request) {
	digestSave[products.Product](w, r, save[products.Product]{
		Name:    productsName,
		URL:     productsURL,
		Feature: productsFeature{},
		Form:    multipartForm{},
		Images:  []string{"image_1", "image_2", "image_3", "image_4"},
		Folder:  productsFolder,
	})
}

func ProductList(w http.ResponseWriter, r *http.Request) {
	digestList[products.Product](w, r, list[products.Product]{
		Name:    productsName,
		URL:     productsURL,
		Feature: productsFeature{},
	})
}

func ProductForm(w http.ResponseWriter, r *http.Request) {
	data, err := digestForm[products.Product](w, r, Form[products.Product]{
		Name:    productsName,
		Feature: productsFeature{},
	})

	if err != nil {
		return
	}

	ctx := r.Context()
	t, err := tags.List(ctx, 0, 9999)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
	}

	f, err := filters.Actives(ctx)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
	}

	data.Extra = struct {
		Tags    []tags.Tag
		Filters []filters.Filter
	}{
		t.Tags,
		f,
	}

	if err := productsFormTpl.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func ProductDelete(w http.ResponseWriter, r *http.Request) {
	digestDelete[products.Product](w, r, delete[products.Product]{
		list: list[products.Product]{
			Name:    productsName,
			URL:     productsURL,
			Feature: productsFeature{},
		},
		Feature: productsFeature{},
	})
}
