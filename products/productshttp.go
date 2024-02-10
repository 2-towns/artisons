package products

import (
	"artisons/conf"
	"artisons/db"
	"artisons/http/contexts"
	"artisons/http/forms"
	"artisons/http/httperrors"
	"artisons/http/httphelpers"
	"artisons/products/filters"
	"artisons/shops"
	"artisons/string/stringutil"
	"artisons/tags"
	"artisons/tags/tree"
	"artisons/templates"
	"artisons/users"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/text/language"
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
		append(files, append(templates.AdminListHandler,
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

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	slug := r.PathValue("slug")
	query := Query{Slug: slug}

	res, err := Search(ctx, query, 0, 1)
	if err != nil {
		httperrors.Page(w, r.Context(), err.Error(), 400)
		return
	}

	if res.Total == 0 {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot find the product", slog.String("slug", slug))
		httperrors.Page(w, r.Context(), "oops the data is not found", 404)
		return
	}

	p := res.Products[0]

	wish := false
	user, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		wish = HasWish(ctx, user.ID, p.ID)
	}

	f, err := filters.Actives(ctx)
	if err != nil {
		httperrors.Page(w, r.Context(), err.Error(), 400)
		return
	}

	data := struct {
		Lang    language.Tag
		Shop    shops.Settings
		Product Product
		Tags    []tree.Leaf
		Wish    bool
		Filters []filters.Filter
	}{
		lang,
		shops.Data,
		p,
		tree.Tree,
		wish,
		f,
	}

	if err := templates.Pages["product"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AdminSaveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(conf.MaxUploadSize); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot parse the form", slog.String("error", err.Error()))
		httperrors.HXCatch(w, ctx, "something went wrong")
		return
	}

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

	p := Product{
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

	p.ID = r.PathValue("id")

	err = p.Validate(ctx)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	query := Query{Slug: p.Slug}
	res, err := Search(ctx, query, 0, 1)
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

	httphelpers.Success(w, "/admin/products")
}

func AdminListHandlerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p := httphelpers.BuildPaginator(r)

	qry := Query{}
	if p.Query != "" {
		qry.Keywords = db.Escape(p.Query)
	}

	res, err := Search(ctx, qry, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	t := productsTpl
	isHX, _ := ctx.Value(contexts.HX).(bool)
	if isHX {
		t = productsHxTpl
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := httphelpers.List[Product]{
		Lang:       lang,
		Items:      res.Products,
		Empty:      len(res.Products) == 0,
		Currency:   conf.Currency,
		Pagination: p.Build(ctx, res.Total, len(res.Products)),
		Page:       "Products",
		Flash:      httphelpers.Flash(w, r),
	}

	if err = t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func AdminFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	var product Product

	if id != "" {
		var err error
		product, err = Find(ctx, id)

		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot parse find the product", slog.Any("id", id), slog.String("error", err.Error()))
			httperrors.Page(w, ctx, "oops the data is not found", 404)
			return
		}
	}

	lang := ctx.Value(contexts.Locale).(language.Tag)
	data := httphelpers.Form[Product]{
		Data:     product,
		Lang:     lang,
		Currency: conf.Currency,
		Page:     "Products",
	}

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

func AdminDeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	err := Delete(ctx, id)
	if err != nil {
		httperrors.HXCatch(w, ctx, err.Error())
		return
	}

	p, _ := url.Parse(r.Header.Get("HX-Current-Url"))
	r.URL.Path = p.Path

	AdminListHandlerHandler(w, r)
}
