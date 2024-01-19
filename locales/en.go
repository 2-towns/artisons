// Package locales provides locale resources for languages
package locales

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func LoadEn() {
	// Errors
	message.SetString(language.English, "the data is not found", "The data is not found.")
	message.SetString(language.English, "the data is invalid", "The data is invalid.")
	message.SetString(language.English, "the session is expired", "Your session is expired, please refresh your page.")
	message.SetString(language.English, "the file %s is not correct", "The file %s is not correct.")
	message.SetString(language.English, "the csv is invalid", "The csv is not valid.")
	message.SetString(language.English, "the image extension is missing in %s", "The image extension is missing in %s.")
	message.SetString(language.English, "the image extension %s is not supported", "The image extension %s is not supported.")
	message.SetString(language.English, "your are not authorized to process this request", "Your are not authorized to access to this page. This will be reported.")
	message.SetString(language.English, "the image %s cannot be downloaded got status %d", "The image %s cannot be downloaded, got status %d.")
	message.SetString(language.English, "oops the data is not found", "Oops..the data is not found.")
	message.SetString(language.English, "something went wrong", "Something went wrong, please try again later. Your request id is %s.")
	message.SetString(language.English, "the page is not found or not accessible anymore", "The page is not found or not accessible anymore.")
	message.SetString(language.English, "oops you just found an error page", "Oops..You just found an error page.")
	message.SetString(language.English, "the request id is %s", "The request id is %s.")
	message.SetString(language.English, "you need to wait before asking another otp", "You need to wait before asking another OTP.")
	message.SetString(language.English, "you reached the max tentatives", "You reached the max tentatives. The OTP is locked now.")
	message.SetString(language.English, "the OTP does not match", "The OTP does not match.")
	message.SetString(language.English, "the user is not found", "The user is not found.")

	// Data
	message.SetString(language.English, "created", "Created")
	message.SetString(language.English, "processing", "Processing")
	message.SetString(language.English, "delivering", "Delivering")
	message.SetString(language.English, "delivered", "Delivered")
	message.SetString(language.English, "canceled", "Canceled")
	message.SetString(language.English, "cash", "Cash")
	message.SetString(language.English, "card", "Card")
	message.SetString(language.English, "bitcoin", "Bitcoin")
	message.SetString(language.English, "wire", "Wire")
	message.SetString(language.English, "collect", "Collect")
	message.SetString(language.English, "home", "Home")

	// Emails
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

	// Texts
	message.SetString(language.English, "The data has been saved successfully.", "The data has been saved successfully.")
	message.SetString(language.English, "The cart is empty.", "The cart is empty.")
	message.SetString(language.English, "Activate demo", "Activate demo")
	message.SetString(language.English, "Disable demo", "Disable demo")
	message.SetString(language.English, "Take me home", "Take me home")
	message.SetString(language.English, "The extensions allowed are .jpg, .jpeg .png.", "The extensions allowed are .jpg, .jpeg .png.")
	message.SetString(language.English, "Add", "Add")
	message.SetString(language.English, "Add article", "Add article")
	message.SetString(language.English, "Address", "Address")
	message.SetString(language.English, "Add product", "Add product")
	message.SetString(language.English, "Articles", "Articles")
	message.SetString(language.English, "Back", "Back")
	message.SetString(language.English, "Banner", "Banner")
	message.SetString(language.English, "Blog", "Blog")
	message.SetString(language.English, "Cancel", "Cancel")
	message.SetString(language.English, "Create", "Create")
	message.SetString(language.English, "City", "City")
	message.SetString(language.English, "Count", "Count")
	message.SetString(language.English, "Delivery", "Delivery")
	message.SetString(language.English, "Delivering", "Delivering")
	message.SetString(language.English, "Description", "Description")
	message.SetString(language.English, "Edit", "Edit")
	message.SetString(language.English, "Email", "Email")
	message.SetString(language.English, "No results found.", "No results found.")
	message.SetString(language.English, "Geolocation", "Geolocation")
	message.SetString(language.English, "Key", "Key")
	message.SetString(language.English, "ID", "ID")
	message.SetString(language.English, "Image", "Image")
	message.SetString(language.English, "Images", "Images")
	message.SetString(language.English, "info", "Info")
	message.SetString(language.English, "Link", "Link")
	message.SetString(language.English, "List", "List")
	message.SetString(language.English, "Logo", "Logo")
	message.SetString(language.English, "Logout", "Logout")
	message.SetString(language.English, "Name", "Name")
	message.SetString(language.English, "Note", "Note")
	message.SetString(language.English, "Notes", "Notes")
	message.SetString(language.English, "Page views", "Page views")
	message.SetString(language.English, "Phone", "Phone")
	message.SetString(language.English, "Offline", "Offline")
	message.SetString(language.English, "Online", "Online")
	message.SetString(language.English, "Optional", "optional")
	message.SetString(language.English, "Orders", "Orders")
	message.SetString(language.English, "Showing %d to %d of %d entries", "Showing %d to %d of %d entries")
	message.SetString(language.English, "Payment", "Payment")
	message.SetString(language.English, "Id, status, delivery or payment", "Id, status, delivery or payment")
	message.SetString(language.English, "Price", "Price")
	message.SetString(language.English, "Products", "Products")
	message.SetString(language.English, "search by title id or sku", "Search by title, id or sku")
	message.SetString(language.English, "Quantity", "Quantity")
	message.SetString(language.English, "Save", "Save")
	message.SetString(language.English, "See", "See")
	message.SetString(language.English, "Seo", "SEO")
	message.SetString(language.English, "Settings", "Settings")
	message.SetString(language.English, "status", "Status")
	message.SetString(language.English, "Success", "Success")
	message.SetString(language.English, "Title", "Title")
	message.SetString(language.English, "The product URL will be generated from this title.", "The product URL will be generated from this title.")
	message.SetString(language.English, "Total", "Total")
	message.SetString(language.English, "Translations", "Translations")
	message.SetString(language.English, "Updated at", "Updated at")
	message.SetString(language.English, "Value", "Value")
	message.SetString(language.English, "Verify", "Verify")
	message.SetString(language.English, "Quantity", "Quantity")
	message.SetString(language.English, "zipcode", "Zipcode")
	message.SetString(language.English, "Add your Google Map API key to enable to geolocation for the customers instead of typing the address.", "Add your Google Map API key to enable to geolocation for the customers instead of typing the address.")
	message.SetString(language.English, "Use the same text as in the template.", "Use the same text as in the template..")
	message.SetString(language.English, "Sign in", "Sign in")
	message.SetString(language.English, "Please enter your email to receive a code and connect to the application.", "Please enter your email to receive a code and connect to the application..")
	message.SetString(language.English, "Login to your account", "Login to your account")
	message.SetString(language.English, "Add a note to the order. The customer will be able to see it.", "Add a note to the order. The customer will be able to see it.. The customer will be able to see it.")
	message.SetString(language.English, "We've sent you an otp code to %s.", "We've sent you an otp code to %s.")
	message.SetString(language.English, "Please click the link to confirm your address or enter the otp code below.", "Please click the link to confirm your address or enter the otp code below.")
	message.SetString(language.English, "Check your inbox", "Check your inbox")
	message.SetString(language.English, "Discount", "Discount")
	message.SetString(language.English, "Most products shared", "Most products shared")
	message.SetString(language.English, "Most products sold", "Most products sold")
	message.SetString(language.English, "Most visited pages", "Most visited pages")
	message.SetString(language.English, "SKU", "SKU")
	message.SetString(language.English, "An internal reference to manage your own data. Only alphanumerics characters are allowed.", "An internal reference to manage your own data. Only alphanumerics characters are allowed.")
	message.SetString(language.English, "Changing the order status will trigger automatic action like notifying the customer.", "Changing the order status will trigger automatic action like notifying the customer.")
	message.SetString(language.English, "summer topQuality", "summer topQuality")
	message.SetString(language.English, "Tags", "Tags")
	message.SetString(language.English, "Separate your tags by a space. Your can separate the words inside a tag by using a uppercase letter.", "Separate your tags by a space. Your can separate the words inside a tag by using a uppercase letter.")
	message.SetString(language.English, "Weight", "Weight")
	message.SetString(language.English, "The product weight in grams.", "The product weight in grams.")
	message.SetString(language.English, "SEO means Search Engine Optimization.", "SEO means Search Engine Optimization.")
	message.SetString(language.English, "The SEO aimed at improving the visibility of a website on search engines.", "The SEO aimed at improving the visibility of a website on search engines.")
	message.SetString(language.English, "Enable shop", "Enable shop")
	message.SetString(language.English, "Choose to make your shop active after you have made changes.", "Choose to make your shop active after you have made changes.")
	message.SetString(language.English, "Use it when you need to perform a maintenance on your shop.", "Use it when you need to perform a maintenance on your shop.")
	message.SetString(language.English, "Enable cache", "Enable cache")
	message.SetString(language.English, "Cache the response into Redis to accelerate the page loading.", "Cache the response into Redis to accelerate the page loading.")
	message.SetString(language.English, "Enable guest checkout", "Enable guest checkout")
	message.SetString(language.English, "Allow orders to be made by non-registered users.", "Allow orders to be made by non-registered users.")
	message.SetString(language.English, "Products per page", "Products per page")
	message.SetString(language.English, "Indicate how many products are displayed on the pages.", "Indicate how many products are displayed on the pages.")
	message.SetString(language.English, "This option is particularly useful for promoting purchases.", "This option is particularly useful for promoting purchases.")
	message.SetString(language.English, "Minimum purchase amount", "Minimum purchase amount")
	message.SetString(language.English, "Indicates the minimum amount that must be in the shopping cart to submit an order.", "Indicates the minimum amount that must be in the shopping cart to submit an order.")
	message.SetString(language.English, "If the amount is not reached, your customer can not complete their purchase.", "If the amount is not reached, your customer can not complete their purchase.")
	message.SetString(language.English, "If you do not want to activate this feature, enter '0' in the field.", "If you do not want to activate this feature, enter '0' in the field.")
	message.SetString(language.English, "Display available quantities on product page", "Display available quantities on product page")
	message.SetString(language.English, "If enabled, your visitors can see the quantities of each object available in stock.", "By enabling this feature, your visitors can see the quantities of each object available in stock.")
	message.SetString(language.English, "Number of days for new products", "Number of days for new products")
	message.SetString(language.English, "Define how many days the product will be considered as 'new'", "Define how many days the product will be considered as 'new'")
	message.SetString(language.English, "Redirect after adding product", "Redirect after adding product")
	message.SetString(language.English, "Redirect the client to the cart page after adding a product to the cart.", "Redirect the client to the cart page after adding a product to the cart.")
	message.SetString(language.English, "Total earning", "Total earning")
	message.SetString(language.English, "Total orders", "Total orders")
	message.SetString(language.English, "Last 7 days", "Last 7 days")
	message.SetString(language.English, "Last 14 days", "Last 14 days")
	message.SetString(language.English, "Last 30 days", "Last 30 days")
	message.SetString(language.English, "Active users", "Active users")
	message.SetString(language.English, "Bounce rate", "Bounce rate")
	message.SetString(language.English, "Top referers", "Top referers")
	message.SetString(language.English, "Top browsers", "Top browsers")
	message.SetString(language.English, "Top systems", "Top systems")
	message.SetString(language.English, "Unique visitors", "Unique visitors")
	message.SetString(language.English, "Visitors", "Visitors")
	message.SetString(language.English, "Total visits", "Total visits")
	message.SetString(language.English, "You can override a translation used in the template by providing you own translation", "You can override a translation used in the template by providing you own translation")
	message.SetString(language.English, "New clients", "New clients")
	message.SetString(language.English, "Access to your products", "Access to your dashboard")
	message.SetString(language.English, "Login to access to  your dashboard", "")
	message.SetString(language.English, "Login to your account", "Login to your account")
	message.SetString(language.English, "Brand color", "Brand color")
	message.SetString(language.English, "Define your brand color used in the theme.", "Define your brand color used in the theme.")
	message.SetString(language.English, "Default image width", "Default image width")
	message.SetString(language.English, "Default image height", "Default image height")
	message.SetString(language.English, "If empty, the theme will try to optimize the image display.", "If empty, the theme will try to optimize the image display.")
	message.SetString(language.English, "Enable by default the fuzzy search", "Enable by default the fuzzy search")
	message.SetString(language.English, "The fuzzy search will return the approximated sequences as product titles...etc", "The fuzzy search will return the approximated sequences as product titles...etc")
	message.SetString(language.English, "Enable by default the exact match search", "Enable by default the exact match search")
	message.SetString(language.English, "The exact match value will look the specific keywords entered by the user.", "The exact match value will look the specific keywords entered by the user.")
	message.SetString(language.English, "The default search behaviour seaches for one the the keywords entered by the user.", "The default search behaviour seaches for one the the keywords entered by the user.")
	message.SetString(language.English, "For the product page, use {{key}} to insert the product title.", "For the product page, use {{key}} to insert the product title.")
	message.SetString(language.English, "For the product page, use {{key}} to insert the beginning of the product description.", "For the product page, use {{key}} to insert the beginning of the product description.")
	message.SetString(language.English, "Url", "Url")
	message.SetString(language.English, "Be careful ! By changing this, the url will need to be reindexed by the search engines.", "Be careful ! By changing this, the url will need to be reindexed by the search engines.")
	message.SetString(language.English, "Add tag", "Add tag")
	message.SetString(language.English, "Label", "Label")
	message.SetString(language.English, "Score", "Score")
	message.SetString(language.English, "Tags", "Tags")
	message.SetString(language.English, "Children", "Children")
	message.SetString(language.English, "Root", "Root")
	message.SetString(language.English, "The tag identifier can contains only alphanumeric characters, no space.", "The tag identifier can contains only alphanumeric characters, no space.")
	message.SetString(language.English, "You cannot change it after the creation.", "You cannot change it after the creation.")
	message.SetString(language.English, "To link a product with this tag, you need to add this key in the product tags.", "To link a product with this tag, you need to add this key in the product tags.")
	message.SetString(language.English, "Define the text displayed on the website.", "Define the text displayed on the website.")
	message.SetString(language.English, "Enter the keys of the children tags separated by ';'.", "Enter the keys of the children tags separated by ';'.")
	message.SetString(language.English, "If the tag does not exist, it will be ignored.", "If the tag does not exist, it will be ignored.")
	message.SetString(language.English, "Define the score inside the root tag list.", "Define the score inside the root tag list.")
	message.SetString(language.English, "The higher the score, the more it will appear first.", "The higher the score, the more it will appear first.")
	message.SetString(language.English, "the tag exists already", "The tag exists already.")

}
