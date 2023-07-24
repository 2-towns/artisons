// Package locales provides locale resources for languages
package locales

import (
	"context"
	"gifthub/http/httputil"
	"net/http"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// ContextKey is the context key used to store the lang
const ContextKey httputil.ContextKey = "lang"

// GetPage returns the common translations for a page
func GetPage(lang language.Tag, name string) map[string]string {
	p := message.NewPrinter(lang)

	return map[string]string{
		"title":                p.Sprintf(name + "_title"),
		"description":          p.Sprintf(name + "_description"),
		"about_us":             p.Sprintf("about_us"),
		"privacy_policy":       p.Sprintf("privacy_policy"),
		"terms_and_conditions": p.Sprintf("terms_and_conditions"),
		"contact_us":           p.Sprintf("contact_us"),
	}
}

// Languages contains the available languages in the application
// Deprecated: Should be moved in configuration
var Languages = "en"

// Default is the default language applied
// Deprecated: Should be moved in configuration
var Default = "en"

// Console is the default language for console
var Console language.Tag = language.English

// UntranslatedError contains the translation key
type UntranslatedError struct {
	Key string
}

// Error is here for error type compatibility
func (e UntranslatedError) Error() string {
	return e.Key
}

// TranslateError translates an error to a user friendly message
func TranslateError(e error, tag language.Tag) string {
	p := message.NewPrinter(tag)
	return p.Sprint(e.Error())
}

// Middleware load the detected language in the context.
// It looks into Accept-Language header and fallback
// to english language when the detected language is
// missing or not recognized.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		langs := strings.Split(r.Header.Get("Accept-Language"), "-")
		lang := langs[0]

		if strings.Contains(Languages, lang) {
			lang = Default
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