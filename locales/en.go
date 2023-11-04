// Package locales provides locale resources for languages
package locales

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func LoadEn() {
	message.SetString(language.English, "home_title", "Homepage")
	message.SetString(language.English, "article_title_required", "The title is required")
	message.SetString(language.English, "article_description_required", "The title is required")
	message.SetString(language.English, "article_image_required", "The title is required")
	message.SetString(language.English, "csv_image_extension_missing", "The image extension is missing in %s.")
	message.SetString(language.English, "csv_image_extension_not_supported", "The image extension %s is not supported.")
	message.SetString(language.English, "csv_bad_file", "The file %s is not correct.")
	message.SetString(language.English, "csv_not_valid", "The csv is not valid.")
	message.SetString(language.English, "csv_line_error", "Found error at line %d. %s")
	message.SetString(language.English, "input_validation", "The field %s is not correct.")
	message.SetString(language.English, "input_required", "The field %s is required.")
	message.SetString(language.English, "user_logout_invalid", "The logout data are required.")
	message.SetString(language.English, "user_magic_code_required", "The magic code is required.")
	message.SetString(language.English, "user_email_invalid", "The email is invalid.")
	message.SetString(language.English, "user_magic_code_invalid", "The magic code is invalid.")
	message.SetString(language.English, "user_device_required", "The device is required.")
	message.SetString(language.English, "user_firstname_required", "The firstname is required.")
	message.SetString(language.English, "user_lastname_required", "The lastname is required.")
	message.SetString(language.English, "user_city_required", "The city is required.")
	message.SetString(language.English, "user_street_required", "The street is required.")
	message.SetString(language.English, "user_zipcode_required", "The zipcode is required.")
	message.SetString(language.English, "user_phone_required", "The phone is required.")
	message.SetString(language.English, "user_wptoken_required", "The token is required.")
	message.SetString(language.English, "user_sid_required", "The session id is required.")
	message.SetString(language.English, "user_not_found", "The user is not found.")
	message.SetString(language.English, "http_bad_status", "Received bad status %d from %s.")
	message.SetString(language.English, "something_went_wrong", "Something went wrong, please try again later.")
	message.SetString(language.English, "home_description", "Home description")
	message.SetString(language.English, "product_url", "/product")
	message.SetString(language.English, "product_pid_required", "The product pid is required.")
	message.SetString(language.English, "user_magik_link_url", "/magic-link.html")
	message.SetString(language.English, "cart_not_found", "Your session is expired, please refresh your page.")
	message.SetString(language.English, "cart_empty", "The cart is empty.")
	message.SetString(language.English, "order_not_found", "The order is not found.")
	message.SetString(language.English, "order_note_required", "The note is required.")
	message.SetString(language.English, "order_bad_status", "The status is not valid.")
	message.SetString(language.English, "order_created_email", "Hi Tralala,\nThank you for your order %s. Bla bla bla")
	message.SetString(language.English, "unauthorized", "Your are not authoried to process the request, this will be reported.")
	message.SetString(language.English, "mail_magic_link", "Hi\n, your order %s has been updated to %s.\nContact us if you need more information.\nThe team.")

}
