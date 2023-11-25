package httpext

import "net/http"

func Redirect(w http.ResponseWriter, r *http.Request, url string, status int) {
	htmx := r.Header.Get("HX-Request") == "true"
	if htmx {
		w.Header().Set("HX-Redirect", url)
	} else {
		http.Redirect(w, r, url, http.StatusFound)
	}
}
