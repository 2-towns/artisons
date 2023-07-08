// Package util provides a set of utility function used by other packages
package util

import (
	"fmt"
	"gifthub/conf"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// ItemsPerPage is the number of items displayed per page or pagination
// Deprecated: Should be moved into the configuration
const ItemsPerPage = 12

// Product is the product representation in the applicatio
type Product struct {
	ID          string            `redis:"id"` // ID is an unique identifier
	Title       string            `redis:"title"`
	Image       string            `redis:"image"`
	Description string            `redis:"description"`
	Price       float64           `redis:"price"`
	Slug        string            `redis:"slug"`
	Links       []string          // Links contains the linked product IDs
	Meta        map[string]string // Meta contains the product options.
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

// Languages contains the available languages in the application
// Deprecated: Should be moved in configuration
var Languages = []string{"en"}

// DefaultLanguage is the default language applied
// Deprecated: Should be moved in configuration
var DefaultLanguage = "en"

// Contains return true if a slice contains the string passed
// in parameter.
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// ContextKey is the type of key used for the context.
// It is necessary to create a specific type for the context, but
// it does not bring added value.
type ContextKey string

// ContextLangKey is the context key used to store the lang
const ContextLangKey ContextKey = "lang"

// Slugify returns the slug representation of a title
func Slugify(title string) string {
	return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
}

// CsvLines representes the lines of a csv
type CsvLines [][]string

var r *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString provides a random unique string
func RandomString() string {
	length := 16
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

// ImageExtensions is the allowed extensions in the application
const ImageExtensions = "jpg jpeg png"

const (
	Online  = "online"  // Make th product available in the application
	Offline = "offline" // Hide th product  in the application
)

// GetImagePath returns the imgproxy path for a file
// Later on, the method should be improve to generate subfolders path,
// if the products are more than the unix file limit
func GetImagePath(pid string, index int) (string, string) {
	folder := fmt.Sprintf("%s/%s", conf.ImgProxyPath, pid)
	return folder, fmt.Sprintf("%s/%d", folder, index)
}
