// Package locales provides locale resources for languages
package locales

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func LoadEn() {
	message.SetString(language.English, "home_title", "Homepage")
	message.SetString(language.English, "csv_image_extension_missing", "The image extension is missing in %s.")
	message.SetString(language.English, "csv_image_extension_not_supported", "The image extension %s is not supported.")
	message.SetString(language.English, "csv_bad_file", "The file %s is not correct.")
	message.SetString(language.English, "csv_not_valid", "The csv is not valid.")
	message.SetString(language.English, "csv_line_error", "Found error at line %d. %s")
	message.SetString(language.English, "input_validation", "The field %s is not correct.")
	message.SetString(language.English, "input_required", "The field %s is required.")
	message.SetString(language.English, "user_password_required", "The password is required.")
	message.SetString(language.English, "user_username_exists", "The username already exists.")
	message.SetString(language.English, "http_bad_status", "Received bad status %d from %s.")
	message.SetString(language.English, "something_went_wrong", "Something went wrong, please try again later.")
	message.SetString(language.English, "home_description", "Home description")
	message.SetString(language.English, "product_url", "/product")

}
