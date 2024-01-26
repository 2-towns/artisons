package templates

import (
	"fmt"
	"gifthub/cache"
	"gifthub/conf"
	"gifthub/http/seo"
	"gifthub/images"
	"gifthub/locales"
	"html/template"
	"log"
	"slices"
	"strings"
	"time"

	"golang.org/x/text/language"
)

type Pagination struct {
	// True if the page is the first page in the pagination
	IsFirst bool

	// True if the page is the lastg page in the pagination
	IsLast bool

	// Pagination numbers availables
	Items []int

	// The max page number
	Max int

	// The current page
	Page int

	// The URL used to retrieve the previous / next page
	URL string

	// The corresponding start number of items displayed
	Start int

	// The corresponding end number of items displayed
	End int

	// The total items available across all the pages
	Total int

	Lang language.Tag
}

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

var AdminList = append(AdminUI, AdminSuccess...)

var Pages map[string]*template.Template = map[string]*template.Template{}
var pages = []string{"home", "wish", "wish-list"}

func init() {
	folder := fmt.Sprintf("%s/web/views/themes/%s", conf.WorkingSpace, conf.DefaultTheme)

	for _, value := range pages {
		tpl, err := Build("base.html").ParseFiles([]string{
			folder + "/base.html",
			folder + fmt.Sprintf("/%s.html", value),
		}...)

		if err != nil {
			log.Panicln(err)
		}

		Pages[value] = tpl
	}
}

func Build(name string) *template.Template {
	return template.New(name).Funcs(template.FuncMap{
		"translate":   locales.Translate,
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
		"metatitle": func(key string, id string) string {
			if id == "" {
				return seo.URLs[key].Title
			}

			return strings.Replace(seo.URLs[key].Title, "{{key}}", id, 1)
		},
		"metadescription": func(key string, id string) string {
			if id == "" {
				return seo.URLs[key].Description
			}

			return strings.Replace(seo.URLs[key].Description, "{{key}}", id, 1)
		},
		"url": func(key string, id string) string {
			if id == "" {
				return seo.URLs[key].URL
			}

			return strings.Replace(seo.URLs[key].URL, "{{id}}", id, 1)
		},
	})
}

// Paginate provides data for pagination template.
// The page parameter is the current page.
// The loaded parameter is the number of loaded items returned by Redis.
// The total is the total items available.
func Paginate(page int, loaded int, total int) Pagination {
	items := []int{}

	if page > 2 {
		items = append(items, page-2)
	}

	if page > 1 {
		items = append(items, page-1)
	}

	items = append(items, page)

	maxp := total / conf.ItemsPerPage

	if total%conf.ItemsPerPage > 0 {
		maxp++
	}

	if page+1 <= maxp {
		items = append(items, page+1)
	}

	if page+2 <= maxp {
		items = append(items, page+2)
	}

	start := (page - 1) * conf.ItemsPerPage
	end := start + loaded

	if loaded > 0 {
		start++
	}

	return Pagination{
		IsFirst: page == 1,
		IsLast:  page == maxp,
		Items:   items,
		Max:     maxp,
		Page:    page,
		Start:   start,
		End:     end,
		Total:   total,
	}
}
