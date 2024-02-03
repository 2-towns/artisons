package httpext

// TODO remove
import (
	"artisons/conf"
	"artisons/http/contexts"
	"net/http"
	"strconv"
)

type pagination struct {
	Page   int
	Offset int
	Num    int
	Query  string
}

func Redirect(w http.ResponseWriter, r *http.Request, url string, status int) {
	isHX, _ := r.Context().Value(contexts.HX).(bool)

	if isHX {
		w.Header().Set("HX-Redirect", url)
	} else {
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func Pagination(r *http.Request) pagination {
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

	return pagination{
		Page:   page,
		Offset: offset,
		Num:    num,
		Query:  q,
	}
}
