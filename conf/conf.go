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

// Languages contains the available languages in the application
var Languages = []string{"en"}

// ItemsPerPage is the number of items displayed per page or pagination
// Deprecated: Should be moved into the configuration
const ItemsPerPage = 12

// Database index for redis
const DatabaseIndex = 0

// Session duration in nanoseconds
const SessionDuration = time.Hour * 24 * 30

// Statistics duration in nanoseconds
// The statistics cannot be kept too long in order to avoid
// the database to be very big.
// A backup should be done and kept for the history.
// Also it could be a could idea to keep some tracking logs.
const StatisticsDuration = time.Hour * 24 * 30 * 3

// Cart duration in nanoseconds
const CartDuration = time.Hour * 24 * 7

// SessionIDCookie is the session id cookie name
const SessionIDCookie = "session_id"

// Magic link duration in nanoseconds
const MagicCodeDuration = time.Minute * 5

// AppURL is the application root URL
const AppURL = "http://localhost:8080"

// EmailUsername is the username for email sending
const EmailUsername = "a3a5f2d396a820"

// EmailPassword is the password for email sending
const EmailPassword = "12fcfd3c6edb95"

// EmailHost is the email host for email sending
const EmailHost = "sandbox.smtp.mailtrap.io"

// EmailPort is the email port for email sending
const EmailPort = "587"

// HasHomeDelivery enabled the "home" delivery if true
const HasHomeDelivery = true

// VapidPublicKey is the public key used for VAPID protocol
const VapidPublicKey = ""

// VapidPrivateKey is the private key used for VAPID protocol
const VapidPrivateKey = ""

// VapidEmail is the email used for VAPID protocol
const VapidEmail = ""

const WebsiteURL = "http://localhost"

// TagMaxDepth is the depth maximum used when looking for
// tags and links.
// Be careful, this setting is very dangerous and could impact badly
// the performance.
const TagMaxDepth = 3

const AdminPrefix = "/admin"

// Pagination returns the start items index and the
// end items index.
func Pagination(page int) (int, int) {
	if page == -1 {
		return 0, -1
	}
	return page * ItemsPerPage, page*ItemsPerPage + ItemsPerPage
}
