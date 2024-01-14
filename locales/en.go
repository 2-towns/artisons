// Package locales provides locale resources for languages
package locales

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func LoadEn() {
	message.SetString(language.English, "", "")
	message.SetString(language.English, "the data is not found", "The data is not found.")

	message.SetString(language.English, "dynamic_status_created", "Created")
	message.SetString(language.English, "dynamic_status_processing", "Processing")
	message.SetString(language.English, "dynamic_status_delivering", "Delivering")
	message.SetString(language.English, "dynamic_status_delivered", "Delivered")
	message.SetString(language.English, "dynamic_status_canceled", "Canceled")
	message.SetString(language.English, "dynamic_payment_cash", "Cash")
	message.SetString(language.English, "dynamic_payment_card", "Card")
	message.SetString(language.English, "dynamic_payment_bitcoin", "Bitcoin")
	message.SetString(language.English, "dynamic_payment_wire", "Wire")
	message.SetString(language.English, "dynamic_delivery_collect", "Collect")
	message.SetString(language.English, "dynamic_delivery_home", "Home")
	message.SetString(language.English, "email_order_confirmation", "Hi %s,\nWoo hoo! Your order is on its way. Your order details can be found below.\n")
	message.SetString(language.English, "email_order_confirmationdate", "Order date: %s\n")
	message.SetString(language.English, "email_order_confirmationfooter", "\nSee you around,\nThe Customer Experience Team at gifthub shop")
	message.SetString(language.English, "email_order_confirmationid", "Order ID: %s\n")
	message.SetString(language.English, "email_order_confirmationsummary", "Here is your order summary:\n\n")
	message.SetString(language.English, "email_order_confirmationtotal", "Order total: %.2f\n\n")
	message.SetString(language.English, "email_otp_login", "Hi,\r\nYou have requested us to send an otp to sign into our application.\r\nPlease use the verification code below to sign in.\r\n\r\n%s\r\n\r\nThe OTP can only be used on the device you initiated the request.\r\nIf you didn't request this, you can ignore this email.\r\n\r\nThanks,\r\nThe support team")
	message.SetString(language.English, "email_otp_subject", "🔒 Your OTP code")
	message.SetString(language.English, "email_order_subject", "Order confirmation %s")
	message.SetString(language.English, "email_order_update", "Order update %s")
	message.SetString(language.English, "email_order_track", "Track your order by clicking on the following link: %s.\n")
	message.SetString(language.English, "error_cart_notfound", "Your session is expired, please refresh your page.")
	message.SetString(language.English, "error_csv_badfile", "The file %s is not correct.")
	message.SetString(language.English, "error_csv_fileinvalid", "The csv is not valid.")
	message.SetString(language.English, "error_csv_imageextensionmissing", "The image extension is missing in %s.")
	message.SetString(language.English, "error_csv_imageextensionnotsupported", "The image extension %s is not supported.")
	message.SetString(language.English, "error_csv_line", "Found error at line %d. %s")
	message.SetString(language.English, "error_login_device", "The device is required.")
	message.SetString(language.English, "error_http_badstatus", "Received bad status %d from %s.")
	message.SetString(language.English, "error_http_blognotfound", "Oops..the article is not found.")
	message.SetString(language.English, "error_http_csrf", "The request is invalid. This will be reported.")
	message.SetString(language.English, "error_http_general", "Something went wrong, please try again later. Your request id is %s.")
	message.SetString(language.English, "error_http_keynotfound", "The keys is not found.")
	message.SetString(language.English, "error_http_notfound", "The page is not found or not accessible anymore.")
	message.SetString(language.English, "error_http_page", "Oops..You just found an error page.")
	message.SetString(language.English, "error_http_productnotfound", "Oops..the product is not found.")
	message.SetString(language.English, "error_http_requestid", "The request id is %s.")
	message.SetString(language.English, "error_http_unauthorized", "Your are not authorized to access to this page. This will be reported.")
	message.SetString(language.English, "error_order_notfound", "The order is not found.")
	message.SetString(language.English, "error_otp_interval", "You need to wait before asking another OTP.")
	message.SetString(language.English, "error_otp_invalid", "The OTP is invalid.")
	message.SetString(language.English, "error_otp_locked", "You reached the max tentatives. The OTP is locked now.")
	message.SetString(language.English, "error_otp_mismatch", "The OTP does not match.")
	message.SetString(language.English, "error_session_idrequired", "The session id is required.")
	message.SetString(language.English, "error_session_notfound", "The user is not found.")
	message.SetString(language.English, "error_user_logout", "The logout data are required.")
	message.SetString(language.English, "input_address_invalid", "The address is invalid.")
	message.SetString(language.English, "input_city_required", "The city is required.")
	message.SetString(language.English, "input_city_invalid", "The city is invalid.")
	message.SetString(language.English, "input_description_invalid", "The description is invalid.")
	message.SetString(language.English, "input_discount_invalid", "The discount is invalid.")
	message.SetString(language.English, "input_email_invalid", "The email is invalid.")
	message.SetString(language.English, "input_email_notadmin", "The email is invalid.")
	message.SetString(language.English, "input_firstname_required", "The firstname is required.")
	message.SetString(language.English, "input_id_required", "The id is required.")
	message.SetString(language.English, "input_image_invalid", "The image is invalid.")
	message.SetString(language.English, "input_image_required", "The image is required.")
	message.SetString(language.English, "input_image_1_required", "The image is required.")
	message.SetString(language.English, "input_image_1_invalid", "The image is invalid.")
	message.SetString(language.English, "input_image_2_invalid", "The image is invalid.")
	message.SetString(language.English, "input_image_3_invalid", "The image is invalid.")
	message.SetString(language.English, "input_image_1_invalid", "The image is invalid.")
	message.SetString(language.English, "input_images_required", "The images are required.")
	message.SetString(language.English, "input_images_invalid", "The images are invalid.")
	message.SetString(language.English, "input_items_invalid", "The items is invalid.")
	message.SetString(language.English, "input_key_invalid", "The key is invalid.")
	message.SetString(language.English, "input_lang_invalid", "The lang is invalid.")
	message.SetString(language.English, "input_locale_invalid", "The locale is invalid.")
	message.SetString(language.English, "input_last_invalid", "The last quantity is invalid.")
	message.SetString(language.English, "input_logo_invalid", "The logo is invalid.")
	message.SetString(language.English, "input_logo_required", "The logo is required.")
	message.SetString(language.English, "input_min_invalid", "The min is invalid.")
	message.SetString(language.English, "input_name_invalid", "The name is invalid.")
	message.SetString(language.English, "input_quantity_invalid", "The quantity is invalid.")
	message.SetString(language.English, "input_lastname_required", "The lastname is required.")
	message.SetString(language.English, "input_otp_required", "The otp is required.")
	message.SetString(language.English, "input_name_required", "The tag is required.")
	message.SetString(language.English, "input_name_reserved", "This keyword is reserved.")
	message.SetString(language.English, "input_note_required", "The note is required.")
	message.SetString(language.English, "input_phone_required", "The phone is required.")
	message.SetString(language.English, "input_pid_required", "The pid is required.")
	message.SetString(language.English, "input_price_invalid", "The price is invalid.")
	message.SetString(language.English, "input_score_required", "The score is required.")
	message.SetString(language.English, "input_search_invalid", "The search input is invalid.")
	message.SetString(language.English, "input_sku_invalid", "The sku is invalid.")
	message.SetString(language.English, "input_status_invalid", "The status is invalid.")
	message.SetString(language.English, "input_street_required", "The street is required.")
	message.SetString(language.English, "input_tag_notfound", "The tag is not found.")
	message.SetString(language.English, "input_taglabel_invalid", "The tag is invalid, must be lowercase alpha characters only.")
	message.SetString(language.English, "input_tagname_invalid", "The tag is invalid, must be lowercase alpha characters only.")
	message.SetString(language.English, "input_tagparent_invalid", "The parent is invalid, must be lowercase alpha characters only.")
	message.SetString(language.English, "input_title_invalid", "The title is invalid.")
	message.SetString(language.English, "input_value_invalid", "The value is invalid.")
	message.SetString(language.English, "input_weight_invalid", "The weight is invalid.")
	message.SetString(language.English, "input_wptoken_required", "The token is required.")
	message.SetString(language.English, "input_zipcode_invalid", "The zipcode is invalid.")
	message.SetString(language.English, "input_zipcode_required", "The zipcode is required.")
	message.SetString(language.English, "text_blog_addsuccess", "The article has been created successfully.")
	message.SetString(language.English, "text_blog_editsuccess", "The article has been updated successfully.")
	message.SetString(language.English, "text_cart_empty", "The cart is empty.")
	message.SetString(language.English, "text_demo_activate", "Activate demo")
	message.SetString(language.English, "text_demo_disable", "Disable demo")
	message.SetString(language.English, "text_error_button", "Take me home")
	message.SetString(language.English, "text_images_sub", "The extensions allowed are .jpg, .jpeg .png.")
	message.SetString(language.English, "text_general_add", "Add")
	message.SetString(language.English, "text_general_addarticle", "Add article")
	message.SetString(language.English, "text_general_address", "Address")
	message.SetString(language.English, "text_general_addproduct", "Add product")
	message.SetString(language.English, "text_general_articles", "Articles")
	message.SetString(language.English, "text_general_back", "Back")
	message.SetString(language.English, "text_general_banner", "Banner")
	message.SetString(language.English, "text_general_blog", "Blog")
	message.SetString(language.English, "text_general_blogsearchplaceholder", "Id, title or description")
	message.SetString(language.English, "text_general_cancel", "Cancel")
	message.SetString(language.English, "text_general_canceled", "Canceled")
	message.SetString(language.English, "text_general_create", "Create")
	message.SetString(language.English, "text_general_created", "Created")
	message.SetString(language.English, "text_general_city", "City")
	message.SetString(language.English, "text_general_count", "Count")
	message.SetString(language.English, "text_general_delete", "Delete")
	message.SetString(language.English, "text_general_delivered", "Delivered")
	message.SetString(language.English, "text_general_delivery", "Delivery")
	message.SetString(language.English, "text_general_delivering", "Delivering")
	message.SetString(language.English, "text_general_description", "Description")
	message.SetString(language.English, "text_general_edit", "Edit")
	message.SetString(language.English, "text_general_email", "Email")
	message.SetString(language.English, "text_general_empty", "No results found.")
	message.SetString(language.English, "text_general_error", "Error")
	message.SetString(language.English, "text_general_geolocation", "Geolocation")
	message.SetString(language.English, "text_general_home", "Home")
	message.SetString(language.English, "text_general_key", "Key")
	message.SetString(language.English, "text_general_id", "ID")
	message.SetString(language.English, "text_general_image", "Image")
	message.SetString(language.English, "text_general_images", "Images")
	message.SetString(language.English, "text_general_info", "Info")
	message.SetString(language.English, "text_general_lang", "Lang")
	message.SetString(language.English, "text_general_link", "Link")
	message.SetString(language.English, "text_general_list", "List")
	message.SetString(language.English, "text_general_locale", "Locale")
	message.SetString(language.English, "text_general_localesupdated", "The translation has been updated.")
	message.SetString(language.English, "text_general_logo", "Logo")
	message.SetString(language.English, "text_general_logout", "Logout")
	message.SetString(language.English, "text_general_name", "Name")
	message.SetString(language.English, "text_general_note", "Note")
	message.SetString(language.English, "text_general_notes", "Notes")
	message.SetString(language.English, "text_general_pageviews", "Page views")
	message.SetString(language.English, "text_general_phone", "Phone")
	message.SetString(language.English, "text_general_offline", "Offline")
	message.SetString(language.English, "text_general_online", "Online")
	message.SetString(language.English, "text_general_optional", "optional")
	message.SetString(language.English, "text_general_orders", "Orders")
	message.SetString(language.English, "text_general_pagination", "Showing %d to %d of %d entries")
	message.SetString(language.English, "text_general_payment", "Payment")
	message.SetString(language.English, "text_general_ordersnoteadded", "The note has been added successfully.")
	message.SetString(language.English, "text_general_orderssearchplaceholder", "Id, status, delivery or payment")
	message.SetString(language.English, "text_general_ordersstatusupdated", "The status has been updated successfully.")
	message.SetString(language.English, "text_general_price", "Price")
	message.SetString(language.English, "text_general_products", "Products")
	message.SetString(language.English, "text_general_productssearchplaceholder", "Search by title, id or sku")
	message.SetString(language.English, "text_general_quantity", "Quantity")
	message.SetString(language.English, "text_general_save", "Save")
	message.SetString(language.English, "text_general_see", "See")
	message.SetString(language.English, "text_general_send", "Send")
	message.SetString(language.English, "text_general_seo", "SEO")
	message.SetString(language.English, "text_general_settings", "Settings")
	message.SetString(language.English, "text_general_status", "Status")
	message.SetString(language.English, "text_general_success", "Success")
	message.SetString(language.English, "text_general_title", "Title")
	message.SetString(language.English, "text_general_titlesub", "The product URL will be generated from this title.")
	message.SetString(language.English, "text_general_total", "Total")
	message.SetString(language.English, "text_general_translations", "Translations")
	message.SetString(language.English, "text_general_updatedat", "Updated at")
	message.SetString(language.English, "text_general_value", "Value")
	message.SetString(language.English, "text_general_verify", "Verify")
	message.SetString(language.English, "text_general_quality", "Quantity")
	message.SetString(language.English, "text_general_zipcode", "Zipcode")
	message.SetString(language.English, "text_geolocation_sub", "Add your Google Map API key to enable to geolocation for the customers instead of typing the address.")
	message.SetString(language.English, "text_locales_keysub", "Use  the same text as in the template.")
	message.SetString(language.English, "text_login_otp", "Otp")
	message.SetString(language.English, "text_login_signin", "Sign in")
	message.SetString(language.English, "text_login_signinsub", "Please enter your email to receive a code and connect to the application.")
	message.SetString(language.English, "text_login_sub", "Login to your account")
	message.SetString(language.English, "text_note_sub", "Add a note to the order. The customer will be able to see it.")
	message.SetString(language.English, "text_orders_processing", "Processing")
	message.SetString(language.English, "text_otp_message", "We've sent you an otp code to %s. Please click the link to confirm your address or enter the otp code below.")
	message.SetString(language.English, "text_otp_sub", "Check your inbox")
	message.SetString(language.English, "text_products_addsuccess", "The product has been created successfully.")
	message.SetString(language.English, "text_products_discount", "Discount")
	message.SetString(language.English, "text_products_editsuccess", "The product has been updated successfully.")
	message.SetString(language.English, "text_products_mostshared", "Most products shared")
	message.SetString(language.English, "text_products_mostsold", "Most products sold")
	message.SetString(language.English, "text_products_mostvisited", "Most visited pages")
	message.SetString(language.English, "text_products_sku", "SKU")
	message.SetString(language.English, "text_products_skusub", "An internal reference to manage your own data. Only alphanumerics characters are allowed.")
	message.SetString(language.English, "text_products_statussub", "Changing the order status will trigger automatic action like notifying the customer. So be careful before updating it !")
	message.SetString(language.English, "text_products_tagsplaceholder", "summer topQuality")
	message.SetString(language.English, "text_products_tags", "Tags")
	message.SetString(language.English, "text_products_tagssub", "Separate your tags by a space. Your can separate the words inside a tag by using a uppercase letter.")
	message.SetString(language.English, "text_products_weight", "Weight")
	message.SetString(language.English, "text_products_weightsub", "The product weight in grams.")
	message.SetString(language.English, "text_seo_message", "SEO means Search Engine Optimization. It represents a set of techniques aimed at improving the visibility of a website on search engines. This section helps you improve the presence of your store on web searches, and therefore reach more potential customers.")
	message.SetString(language.English, "text_settings_active", "Enable shop")
	message.SetString(language.English, "text_settings_activesub", "Choose to make your shop active after you have made changes. Use it when you need to performe a maintenance on your shop.")
	message.SetString(language.English, "text_settings_advanced", "Advanced search")
	message.SetString(language.English, "text_settings_advancedsub", "Enable the advanced search allowing the user to select multiple criteria.")
	message.SetString(language.English, "text_settings_cache", "Enable cache")
	message.SetString(language.English, "text_settings_cachesub", "Cache the response into Redis to accelerate the page loading.")
	message.SetString(language.English, "text_settings_guest", "Enable guest checkout")
	message.SetString(language.English, "text_settings_guestsub", "Allow orders to be made by non-registered users.")
	message.SetString(language.English, "text_settings_items", "Products per page")
	message.SetString(language.English, "text_settings_itemssub", "Indicate how many products are displayed on the pages.")
	message.SetString(language.English, "text_settings_last", "Last quantity")
	message.SetString(language.English, "text_settings_lastsub", "You can display an alert when a stock of your products gets low. This option is particularly useful for promoting purchases. To configure this feature, enter the field value at which an alert message should appear on your store.")
	message.SetString(language.English, "text_settings_min", "Minimum purchase amount")
	message.SetString(language.English, "text_settings_minsub", "Indicates the minimum amount that must be in the shopping cart to submit an order. If the amount in this field is not reached, your customer can not complete their purchase. If you do not want to activate this feature, enter '0' in the field.")
	message.SetString(language.English, "text_settings_quantity", "Display available quantities on product page")
	message.SetString(language.English, "text_settings_quantitysub", "By enabling this feature, your visitors can see the quantities of each object available in stock. Displaying this information can be used to stimulate sales in the case where the quantity in stock is low.")
	message.SetString(language.English, "text_settings_new", "Number of days for new products")
	message.SetString(language.English, "text_settings_newsub", "Define how many days the product will be considered as 'new'")
	message.SetString(language.English, "text_settings_redirect", "Redirect after adding product")
	message.SetString(language.English, "text_settings_redirectsub", "When a product is added to the shopping cart and the AJAX version of the cart mode is disabled, the client can be directed to the shopping cart summary or stay in the current page.")
	message.SetString(language.English, "text_settings_stock", "Enable stock management")
	message.SetString(language.English, "text_settings_stocksub", "By default you should leave this feature enabled. This affects the entire inventory management of your store.")
	message.SetString(language.English, "text_settings_editsuccess", "The success has been updated successfully.")
	message.SetString(language.English, "text_statistics_totalearning", "Total earning")
	message.SetString(language.English, "text_statistics_totalorders", "Total orders")
	message.SetString(language.English, "text_statistics_7days", "Last 7 days")
	message.SetString(language.English, "text_statistics_14days", "Last 14 days")
	message.SetString(language.English, "text_statistics_30days", "Last 30 days")
	message.SetString(language.English, "text_statistics_activeusers", "Active users")
	message.SetString(language.English, "text_statistics_bouncerate", "Bounce rate")
	message.SetString(language.English, "text_statistics_topreferers", "Top referers")
	message.SetString(language.English, "text_statistics_topbrowsers", "Top browsers")
	message.SetString(language.English, "text_statistics_topsystems", "Top systems")
	message.SetString(language.English, "text_statistics_unique", "Unique visitors")
	message.SetString(language.English, "text_statistics_visitors", "Visitors")
	message.SetString(language.English, "text_statistics_visits", "Total visits")
	message.SetString(language.English, "text_translations_message", "You can manage the text translation in this section. The translation key is the one used into the template, so you can override it or translate it.")
	message.SetString(language.English, "text_users_new", "New clients")
	message.SetString(language.English, "seo_dashboard_title", "Access to your products")
	message.SetString(language.English, "seo_dashboard_title", "Access to your dashboard")
	message.SetString(language.English, "seo_login_description", "")
	message.SetString(language.English, "seo_login_title", "Login to your account")

}
