package admin

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/http/forms"
	"artisons/http/httperrors"
	"artisons/http/pages"
	"artisons/products"
	"artisons/products/filters"
	"artisons/string/stringutil"
	"artisons/tags"
	"artisons/templates"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var productsTpl *template.Template
var productsHxTpl *template.Template
var productsFormTpl *template.Template

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

func ProductSave(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.FormValue("price") == "" {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the empty price")
		httperrors.HXCatch(w, ctx, "input:price")
		return
	}

	if r.FormValue("quantity") == "" {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot use the empty quantity")
		httperrors.HXCatch(w, ctx, "input:quantity")
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the price", slog.String("price", r.FormValue("price")), slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "input:price")
		return
	}

	quantity, err := strconv.ParseInt(r.FormValue("quantity"), 10, 64)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the quantity", slog.String("quantity", r.FormValue("quantity")), slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "input:quantity")
		return
	}

	var discount float64 = 0
	if r.FormValue("discount") != "" {
		val, err := strconv.ParseFloat(r.FormValue("discount"), 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the discount", slog.String("discount", r.FormValue("discount")), slog.String("error", err.Error()))
			httperrors.HXCatch(w, ctx, "input:discount")
			return
		}
		discount = val
	}

	var weight float64 = 0
	if r.FormValue("weight") != "" {
		val, err := strconv.ParseFloat(r.FormValue("weight"), 64)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse the weight", slog.String("weight", r.FormValue("discount")), slog.String("error", err.Error()))
			httperrors.HXCatch(w, ctx, "input:weight")
			return
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
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
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

	err = p.Validate(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	query := products.Query{Slug: p.Slug}
	res, err := products.Search(ctx, query, 0, 1)
	if err != nil || res.Total > 0 && (res.Products[0].ID != p.ID) {
		httperrors.HXCatch(w, ctx, "input:slug")
		return
	}

	images := []string{"image_1", "image_2", "image_3", "image_3"}
	files, err := forms.Upload(r, "products", images)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	p.Image1 = files[0]

	if files[1] != "" {
		p.Image2 = files[1]
	}

	if files[2] != "" {
		p.Image3 = files[2]
	}

	if files[3] != "" {
		p.Image4 = files[4]
	}

	if p.ID == "" && p.Image1 == "" {
		slog.LogAttrs(ctx, slog.LevelError, "cannot process the empty image one")
		httperrors.HXCatch(w, ctx, "input:image_1")
		return
	}

	_, err = p.Save(ctx)
	if err != nil {
		forms.RollbackUpload(ctx, files)
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	Success(w, "/admin/products.html")
}

func ProductList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := ctx.Value(contexts.Pagination).(pages.Paginator)

	qry := products.Query{}
	if p.Query != "" {
		qry.Keywords = db.Escape(p.Query)
	}

	res, err := products.Search(ctx, qry, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	t := productsTpl
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		t = productsHxTpl
	}

	data := pages.Datalist(ctx, res.Products)
	data.Pagination = p.Build(ctx, res.Total, len(res.Products))
	data.Page = "Products"

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func ProductForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var product products.Product

	if id != "" {
		var err error
		product, err = products.Find(ctx, id)

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse find the product", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}
	}

	data := pages.Dataform[products.Product](ctx, product)
	data.Page = "Products"

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
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	err := products.Delete(ctx, id)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	pages.UpdateQuery(r)

	ProductList(w, r)
}
