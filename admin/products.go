package admin

import (
	"context"
	"errors"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/http/contexts"
	"gifthub/http/httpext"
	"gifthub/products"
	"gifthub/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

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

func (f productsFeature) Search(ctx context.Context, q string, offset, num int) (httpext.SearchResults[products.Product], error) {
	query := products.Query{}
	if q != "" {
		query.Keywords = db.Escape(q)
	}

	res, err := products.Search(ctx, query, offset, num)

	return httpext.SearchResults[products.Product]{
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

	p := products.Product{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Status:      r.FormValue("status"),
		Sku:         r.FormValue("sku"),
		Tags:        strings.Split(r.FormValue("tags"), ";"),
		Price:       price,
		Discount:    discount,
		Weight:      weight,
		Quantity:    int(quantity),
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

func (f productsFeature) FormTemplate(ctx context.Context, w http.ResponseWriter) *template.Template {
	return productsFormTpl
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

func ProductSave(w http.ResponseWriter, r *http.Request) {
	httpext.DigestSave[products.Product](w, r, httpext.Save[products.Product]{
		Name:    productsName,
		URL:     productsURL,
		Feature: productsFeature{},
		Form:    httpext.MultipartForm{},
		Images:  []string{"image_1", "image_2", "image_3", "image_4"},
		Folder:  productsFolder,
	})
}

func ProductList(w http.ResponseWriter, r *http.Request) {
	httpext.DigestList[products.Product](w, r, httpext.List[products.Product]{
		Name:    productsName,
		URL:     productsURL,
		Feature: productsFeature{},
	})
}

func ProductForm(w http.ResponseWriter, r *http.Request) {
	httpext.DigestForm[products.Product](w, r, httpext.Form[products.Product]{
		Name:    productsName,
		Feature: productsFeature{},
	})
}

func ProductDelete(w http.ResponseWriter, r *http.Request) {
	httpext.DigestDelete[products.Product](w, r, httpext.Delete[products.Product]{
		List: httpext.List[products.Product]{
			Name:    productsName,
			URL:     productsURL,
			Feature: productsFeature{},
		},
		Feature: productsFeature{},
	})
}
