package admin

import (
	"artisons/conf"
	"artisons/http/cookies"
	"net/http"
	"time"
)

func Success(w http.ResponseWriter, url string) {
	cookie := &http.Cookie{
		Name:     cookies.FlashMessage,
		Value:    "The data has been saved successfully.",
		MaxAge:   int(time.Minute.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   conf.Cookie.Secure,
		Domain:   conf.Cookie.Domain,
	}

	http.SetCookie(w, cookie)
	w.Header().Set("HX-Redirect", url)
	w.Write([]byte(""))
}
