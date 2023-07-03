package util

import (
	"html/template"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var Templates = template.Must(template.ParseFiles("views/base.html", "views/home.html"))

const ItemsPerPage = 12

type Product struct {
	ID          string  `redis:"id"`
	Title       string  `redis:"title"`
	Image       string  `redis:"image"`
	Description string  `redis:"description"`
	Price       float64 `redis:"price"`
	Slug        string  `redis:"slug"`
	Links       []string
	Meta        map[string]string
}

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

var Languages = []string{"en"}

var DefaultLanguage = "en"

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

type ContextKey string

const ContextLangKey ContextKey = "lang"

func Slugify(title string) string {
	return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
}
