package httpext

import (
	"artisons/http/contexts"
	"net/http"
)

func Redirect(w http.ResponseWriter, r *http.Request, url string, status int) {
	isHX, _ := r.Context().Value(contexts.HX).(bool)

	if isHX {
		w.Header().Set("HX-Redirect", url)
	} else {
		http.Redirect(w, r, url, http.StatusFound)
	}
}
