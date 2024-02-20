package website

import (
	"artisons/blog"
	"artisons/http/contexts"
	"artisons/http/httperrors"
	"artisons/http/httphelpers"
	"artisons/products"
	"artisons/shops"
	"artisons/tags/tree"
	"artisons/templates"
	"artisons/users"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/language"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	slug := strings.Replace(r.URL.Path, ".html", "", 1)
	slug = strings.Replace(slug, "/", "", 1)

	query := blog.Query{Slug: slug, Type: "cms"}
	res, err := blog.Search(ctx, query, 0, 1)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot get the article", slog.String("slug", slug), slog.String("error", err.Error()))
		httperrors.Page(w, r.Context(), err.Error(), 400)
		return
	}

	if res.Total == 0 {
		slog.LogAttrs(ctx, slog.LevelInfo, "cannot find the article", slog.String("slug", slug))
		httperrors.Page(w, r.Context(), "oops the data is not found", 404)
		return
	}

	s := res.Articles[0]

	data := struct {
		Lang    language.Tag
		Shop    shops.Settings
		Article blog.Article
		Tags    []tree.Leaf
	}{
		lang,
		shops.Data,
		s,
		tree.Tree,
	}

	if err := templates.Pages["static"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	data := struct {
		Lang language.Tag
		Shop shops.Settings
		Tags []tree.Leaf
	}{
		lang,
		shops.Data,
		tree.Tree,
	}

	w.Header().Set("Content-Type", "text/html")

	if err := templates.Pages["categories"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

// Home loads the most recent products in order to
// display them on the home page.
func Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)

	p := []products.Product{}
	// if err != nil {
	// 	slog.LogAttrs(ctx, slog.LevelError, "cannot get the products", slog.String("error", err.Error()))
	// 	httperrors.Page(w, r.Context(), "something went wrong", 400)
	// 	return
	// }

	wishes := []string{}
	user, ok := ctx.Value(contexts.User).(users.User)
	if ok {
		var err error
		wishes, err = products.Wishes(ctx, user.ID)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "cannot get the wishes", slog.String("error", err.Error()))
		}
	}

	data := struct {
		Lang     language.Tag
		Shop     shops.Settings
		Products []products.Product
		Tags     []tree.Leaf
		Wishes   []string
	}{
		lang,
		shops.Data,
		p,
		tree.Tree,
		wishes,
	}

	if err := templates.Pages["home"].Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lang := ctx.Value(contexts.Locale).(language.Tag)
	p := httphelpers.BuildPaginator(r)

	q := r.URL.Query()
	var min float32 = 0
	if q.Has("min") {
		if val, err := strconv.ParseFloat(q.Get("min"), 32); err == nil {
			min = float32(val)
		}
	}

	var max float32 = 0
	if q.Has("max") {
		if val, err := strconv.ParseFloat(q.Get("max"), 32); err == nil {
			max = float32(val)
		}
	}

	meta := map[string][]string{}
	for key, val := range q {
		if key == "min" || key == "max" || key == "q" || key == "tags" {
			continue
		}

		meta[key] = val
	}

	query := products.Query{
		PriceMin: min,
		PriceMax: max,
		Keywords: q.Get("q"),
		Tags:     q["tags"],
		Meta:     meta,
	}

	res, err := products.Search(ctx, query, p.Offset, p.Num)
	if err != nil {
		httperrors.Catch(w, ctx, err.Error(), 500)
		return
	}

	pag := p.Build(ctx, res.Total, len(res.Products))

	data := struct {
		Lang       language.Tag
		Shop       shops.Settings
		Tags       []tree.Leaf
		Products   []products.Product
		Empty      bool
		Pagination httphelpers.Pagination
	}{
		lang,
		shops.Data,
		tree.Tree,
		res.Products,
		len(res.Products) == 0,
		pag,
	}

	var t *template.Template
	isHX, _ := ctx.Value(contexts.HX).(bool)

	if isHX {
		t = templates.Pages["hx-search"]
	} else {
		t = templates.Pages["search"]
	}

	w.Header().Set("Content-Type", "text/html")

	if err := t.Execute(w, &data); err != nil {
		slog.Error("cannot render the template", slog.String("error", err.Error()))
	}
}
