package httphelpers

import (
	"artisons/conf"
	"artisons/http/cookies"
	"context"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/text/language"
)

type Paginator struct {
	Page   int
	Offset int
	Num    int
	Query  string
}

type List[T any] struct {
	Lang       language.Tag
	Page       string
	Items      []T
	Empty      bool
	Currency   string
	Pagination Pagination
	Flash      string
}

type Form[T any] struct {
	Lang     language.Tag
	Page     string
	Data     T
	Currency string
	Extra    interface{}
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

func Success(w http.ResponseWriter, url string) {
	c := NewCookie(cookies.FlashMessage, "The data has been saved successfully.", int(time.Minute.Seconds()))
	http.SetCookie(w, &c)

	w.Header().Set("HX-Redirect", url)
	w.Write([]byte(""))
}

func Flash(w http.ResponseWriter, r *http.Request) string {
	c, err := r.Cookie(cookies.FlashMessage)
	if err == nil && c != nil {
		flash := c.Value

		cookie := NewCookie(cookies.FlashMessage, flash, -1)
		http.SetCookie(w, &cookie)

		return flash
	}

	return ""
}

func BuildPaginator(r *http.Request) Paginator {
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

	return Paginator{
		Page:   page,
		Offset: offset,
		Num:    num,
		Query:  q,
	}
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

func NewCookie(name string, val string, max int) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    val,
		MaxAge:   max,
		Path:     "/",
		HttpOnly: true,
		Secure:   conf.Cookie.Secure,
		Domain:   conf.Cookie.Domain,
		SameSite: http.SameSiteStrictMode,
	}
}
