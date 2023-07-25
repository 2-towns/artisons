// Package conf manages the application configuration
package conf

import "time"

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

// Database index for redis
const DatabaseIndex = 0

// Session duration in nanoseconds
const SessionDuration = time.Hour * 24 * 30

// SessionIDCookie is the session id cookie name
const SessionIDCookie = "session_id"

// Magic link duration in nanoseconds
const MagicCodeDuration = time.Minute * 5

// AppURL is the application root URL
const AppURL = "http://localhost:8080"
