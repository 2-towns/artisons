// Package locales provides locale resources for languages
package locales

import (
	"context"
	"gifthub/conf"
	"gifthub/http/contexts"
	"net/http"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var trans = map[language.Tag]*message.Printer{
	language.English: message.NewPrinter(language.English),
}

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

// Default is the default language applied
// Deprecated: Should be moved in configuration
var Default = language.English

// Console is the default language for console
var Console language.Tag = language.English

// Middleware load the detected language in the context.
// It looks into Accept-Language header and fallback
// to english language when the detected language is
// missing or not recognized.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		langs := strings.Split(r.Header.Get("Accept-Language"), "-")
		lang := langs[0]
		var tag language.Tag

		if !slices.Contains(conf.Languages, lang) {
			tag = Default
		} else {
			tag = language.Make(lang)
		}

		// create new context from `r` request context, and assign key `"user"`
		// to value of `"123"`
		ctx := context.WithValue(r.Context(), contexts.Locale, tag)

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new `"user"` key
		// will be accessible from this point forward.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type TranslatorFunc func(message.Reference, ...interface{}) string

func Translator(lang language.Tag) TranslatorFunc {
	p := message.NewPrinter(lang)

	return p.Sprintf
}

func Translate(l language.Tag, msg string, attr ...interface{}) string {
	t := trans[l]

	if t == nil {
		t = trans[Default]
	}

	return t.Sprintf(msg, attr...)
}
