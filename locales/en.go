package locales

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func LoadEn() {
	message.SetString(language.English, "home_title", "Homepage")
	message.SetString(language.English, "home_description", "Home description")
	message.SetString(language.English, "product_url", "/product")
}
