package templates

import (
	"fmt"
	"gifthub/admin/urls"
	"gifthub/cache"
	"gifthub/conf"
	"gifthub/locales"
	"html/template"
	"strings"
	"time"

	"golang.org/x/text/language"
)

type Pagination struct {
	IsFirst bool
	IsLast  bool
	Items   []int
	Max     int
	Page    int
	URL     string
	Start   int
	End     int
	Total   int
	Lang    language.Tag
}

type Image struct {
	Name  string
	Value string
}

func Build(name string) *template.Template {
	return template.New(name).Funcs(template.FuncMap{
		"translate":   locales.Translate,
		"cachebuster": cache.Buster,
		"urls":        urls.Get,
		"date": func(t time.Time) string {
			return t.Format("02 Jan 2006")
		},
		"twodigits": func(f float64) string {
			return fmt.Sprintf("%.2f", f)
		},
		"join": func(values []string, sep string) string {
			return strings.Join(values, sep)
		},
	})
}

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
