package pages

import (
	"artisons/conf"
	"artisons/http/contexts"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
)

type listData[T any] struct {
	Lang       language.Tag
	Page       string
	Items      []T
	Empty      bool
	Currency   string
	Pagination Pagination
	Flash      string
}

type formData[T any] struct {
	Lang     language.Tag
	Page     string
	Data     T
	Currency string
	Extra    interface{}
}

type Paginator struct {
	Page   int
	Offset int
	Num    int
	Query  string
}

func Datalist[T any](ctx context.Context, data []T) listData[T] {
	lang := ctx.Value(contexts.Locale).(language.Tag)
	flash, _ := ctx.Value(contexts.Flash).(string)

	return listData[T]{
		Lang:     lang,
		Items:    data,
		Empty:    len(data) == 0,
		Currency: conf.Currency,
		Flash:    flash,
	}
}

func Dataform[T any](ctx context.Context, data T) formData[T] {
	lang := ctx.Value(contexts.Locale).(language.Tag)

	return formData[T]{
		Lang:     lang,
		Currency: conf.Currency,
		Data:     data,
	}
}

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

func (p Paginator) Build(ctx context.Context, total, loaded int) Pagination {
	items := []int{}

	if p.Page > 2 {
		items = append(items, p.Page-2)
	}

	if p.Page > 1 {
		items = append(items, p.Page-1)
	}

	items = append(items, p.Page)

	maxp := total / conf.ItemsPerPage

	if total%conf.ItemsPerPage > 0 {
		maxp++
	}

	if p.Page+1 <= maxp {
		items = append(items, p.Page+1)
	}

	if p.Page+2 <= maxp {
		items = append(items, p.Page+2)
	}

	start := (p.Page - 1) * conf.ItemsPerPage
	end := start + loaded

	if loaded > 0 {
		start++
	}

	return Pagination{
		IsFirst: p.Page == 1,
		IsLast:  p.Page == maxp,
		Items:   items,
		Max:     maxp,
		Page:    p.Page,
		Start:   start,
		End:     end,
		Total:   total,
	}
}

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var page = 1
		uquery := r.URL.Query()

		p := uquery.Get("page")
		if p != "" {
			if d, err := strconv.ParseInt(p, 10, 32); err == nil && d > 0 {
				page = int(d)
			}
		}

		q := uquery.Get("q")
		offset := 0
		if page > 0 {
			offset = (page - 1) * conf.ItemsPerPage
		}
		num := offset + conf.ItemsPerPage

		ctx := r.Context()
		ctx = context.WithValue(ctx, contexts.Pagination, Paginator{
			Page:   page,
			Offset: offset,
			Num:    num,
			Query:  q,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UpdateQuery(r *http.Request) {
	page := r.FormValue("page")
	r.URL.RawQuery = fmt.Sprintf("page=%s", page)

	query := r.FormValue("query")
	if query != "" {
		r.URL.RawQuery += fmt.Sprintf("&query=%s", query)
	}
}
