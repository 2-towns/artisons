// Package locales provides locale resources for languages
package locales

import (
	"gifthub/http/httputil"

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

// DefaultLanguage is the default language applied
// Deprecated: Should be moved in configuration
var DefaultLanguage = "en"
