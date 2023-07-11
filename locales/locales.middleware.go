package locales

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/text/language"
)

// Lang load the detected language in the context.
// It looks into Accept-Language header and fallback
// to english language when the detected language is
// missing or not recognized.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		langs := strings.Split(r.Header.Get("Accept-Language"), "-")
		lang := langs[0]

		if strings.Contains(Languages, lang) {
			lang = DefaultLanguage
		}

		var tag language.Tag

		switch lang {
		case "en":
			tag = language.English
		default:
			tag = language.English
		}

		// create new context from `r` request context, and assign key `"user"`
		// to value of `"123"`
		ctx := context.WithValue(r.Context(), ContextKey, tag)

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new `"user"` key
		// will be accessible from this point forward.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
