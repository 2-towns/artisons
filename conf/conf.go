// Package conf manages the application configuration
package conf

// ImgProxyPath is the path to imgproxy folder
var ImgProxyPath = "../web/images"

// IsCurrencySupported returns true if the currency is supported
// in the application
func IsCurrencySupported(c string) bool {
	return c == "EUR"
}

// DefaultMerchantId is the default merchant id
var DefaultMID = "1234"

// ItemsPerPage is the number of items displayed per page or pagination
// Deprecated: Should be moved into the configuration
const ItemsPerPage = 12
