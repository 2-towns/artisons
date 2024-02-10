package templates

import (
	"artisons/cache"
	"artisons/conf"
	"artisons/images"
	"artisons/locales"
	"artisons/seo/urls"
	"fmt"
	"html/template"
	"log"
	"slices"
	"strings"
	"time"
)

type Image struct {
	Name  string
	Value string
}

var AdminTable = []string{
	conf.WorkingSpace + "web/views/admin/icons/arrow-right.svg",
	conf.WorkingSpace + "web/views/admin/icons/arrow-left.svg",
	conf.WorkingSpace + "web/views/admin/icons/trash.svg",
	conf.WorkingSpace + "web/views/admin/icons/edit.svg",
	conf.WorkingSpace + "web/views/admin/icons/question-mark.svg",
	conf.WorkingSpace + "web/views/admin/icons/success.svg",
	conf.WorkingSpace + "web/views/admin/pagination.html",
}

var AdminUI = []string{
	conf.WorkingSpace + "web/views/admin/base.html",
	conf.WorkingSpace + "web/views/admin/ui.html",
	conf.WorkingSpace + "web/views/admin/icons/home.svg",
	conf.WorkingSpace + "web/views/admin/icons/building-store.svg",
	conf.WorkingSpace + "web/views/admin/icons/receipt.svg",
	conf.WorkingSpace + "web/views/admin/icons/article.svg",
	conf.WorkingSpace + "web/views/admin/icons/settings.svg",
	conf.WorkingSpace + "web/views/admin/icons/seo.svg",
	conf.WorkingSpace + "web/views/admin/icons/tag.svg",
	conf.WorkingSpace + "web/views/admin/icons/filter.svg",
}

var AdminSuccess = []string{
	conf.WorkingSpace + "web/views/admin/icons/success.svg",
	conf.WorkingSpace + "web/views/admin/alert-success.html",
}

var AdminListHandler = append(AdminUI, AdminSuccess...)

var Pages map[string]*template.Template = map[string]*template.Template{}

func buildTemplate(key string, files []string) {
	folder := fmt.Sprintf("%sweb/views/themes/%s", conf.WorkingSpace, conf.DefaultTheme)

	f := []string{}

	if !strings.HasPrefix(key, "hx") {
		f = append(f, folder+"/base.html")
	}

	for _, file := range files {
		if strings.Contains(file, "/") {
			f = append(f, file)
		} else {
			f = append(f, folder+"/"+file)
		}
	}

	parts := strings.Split(f[0], "/")

	tpl, err := Build(parts[len(parts)-1]).ParseFiles(f...)

	if err != nil {
		log.Panicln(err)
	}

	Pages[key] = tpl
}

func init() {
	buildTemplate("home", []string{"home.html"})
	buildTemplate("login", []string{
		"login.html",
		fmt.Sprintf("%s/web/views/login.html", conf.WorkingSpace),
	})
	buildTemplate("wish", []string{"wish.html", "hx-wish.html"})
	buildTemplate("hx-wish", []string{"hx-wish.html"})
	buildTemplate("blog", []string{"blog.html", "hx-blog.html"})
	buildTemplate("hx-blog", []string{"hx-blog.html"})
	buildTemplate("static", []string{"static.html"})
	buildTemplate("orders", []string{"orders.html", "hx-orders.html"})
	buildTemplate("account", []string{"account.html"})
	buildTemplate("hx-orders", []string{"hx-orders.html"})
	buildTemplate("search", []string{"search.html", "hx-search.html"})
	buildTemplate("hx-search", []string{"hx-search.html"})
	buildTemplate("order", []string{"order.html"})
	buildTemplate("categories", []string{"categories.html"})
	buildTemplate("product", []string{"product.html"})
	buildTemplate("cart", []string{"cart.html", "hx-cart.html"})
	buildTemplate("hx-cart", []string{"hx-cart.html"})
	buildTemplate("address", []string{
		fmt.Sprintf("%s/web/views/address.html", conf.WorkingSpace),
	})
	buildTemplate("hx-success", []string{
		fmt.Sprintf("%s/web/views/success.html", conf.WorkingSpace),
	})
	buildTemplate("hx-input-error", []string{
		fmt.Sprintf("%s/web/views/input-error.html", conf.WorkingSpace),
	})
	buildTemplate("delivery", []string{
		fmt.Sprintf("%s/web/views/delivery.html", conf.WorkingSpace),
	})
}

func Build(name string) *template.Template {
	return template.New(name).Funcs(template.FuncMap{
		"translate":   locales.Translate,
		"uitranslate": locales.UITranslate,
		"cachebuster": cache.Buster,
		"date": func(t time.Time) string {
			return t.Format("02 Jan 2006")
		},
		"datetime": func(t time.Time) string {
			return t.Format("02 Jan 2006 15:04:05")
		},
		"twodigits": func(f float64) string {
			return fmt.Sprintf("%.2f", f)
		},
		"join": func(values []string, sep string) string {
			return strings.Join(values, sep)
		},
		"contains": func(values []string, value string) bool {
			return slices.Contains(values, value)
		},

		"image": func(id, width, height string, cachebuster time.Time) string {
			return images.URL(id, images.Options{
				Width:       width,
				Height:      height,
				Cachebuster: cachebuster.Unix(),
			})
		},
		"meta": func(key string, t string, id string) string {
			return strings.Replace(urls.Get(key, t), "{{key}}", id, 1)
		},
	})
}
