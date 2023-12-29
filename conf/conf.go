// Package conf manages the application configuration
package conf

import (
	"os"
	"time"
)

// ImgProxyPath is the path to imgproxy folder
var ImgProxyPath = WorkingSpace + "web/images"

// ImgProxyURL is the url to imgproxy
var ImgProxyURL = "http://localhost:8080/"

// IsCurrencySupported returns true if the currency is supported
// in the application
func IsCurrencySupported(c string) bool {
	return c == "EUR"
}

const Currency = "EUR"

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

// OtpDuration in nanoseconds
const OtpDuration = time.Minute * 5

// OtpInterval is the time minimum between two otp attemps
const OtpInterval = time.Minute / 10

// OtpAttempts is the maximum attemps for an otp
const OtpAttempts = 3

// OtpDemo allows to use 111111 as otp.
// It should be used only for testing purpose.
const OtpDemo = true

// Cachebuster is used to bypass the assets cache
var Cachebuster = time.Now().Unix()

// AppURL is the application root URL
const AppURL = "http://localhost:8080"

var Email = struct {
	From     string
	Host     string
	Domain   string
	Username string
	Password string
	Port     string
	Dry      bool
}{
	From:     "hello@debugmail.io",
	Domain:   "debugmail.io",
	Host:     "sandbox.smtp.mailtrap.io",
	Username: "a3a5f2d396a820",
	Password: "12fcfd3c6edb95",
	Port:     "25",
	Dry:      os.Getenv("EMAIL_DRY") == "1",
}

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

// ServerAddr is the server start poi
const ServerAddr = ":8080"

// Disable robotx
const Debug = false

// DashboardItems give the numbers of items for most XX statistics
const DashboardMostItems = 5

// MaxUploadSize is the max size of the body for a
// multipart request. 10 Mb.
const MaxUploadSize = 1024 * 1024 * 10

// WorkingSpace is the project root folder. Mainly used for testing
var WorkingSpace = os.Getenv("WORKSPACE_DIR") + "/"

// ImagesAllowed defines the image extensions supported by file upload
var ImagesAllowed = []string{"image/jpg", "image/jpeg", "image/png"}

// UploadFolder is the upload folder destination
const UploadFolder = "web/upload"

var Cookie = struct {
	Domain string
	Secure bool
	MaxAge float64
}{
	Domain: "",
	Secure: os.Getenv("COOKIE_SECURE") == "1",
	// https://chromestatus.com/feature/4887741241229312
	MaxAge: time.Hour.Seconds() * 24 * 400,
}

// Pagination returns the start items index and the
// end items index.
func Pagination(page int) (int, int) {
	if page == -1 {
		return 0, -1
	}
	return page * ItemsPerPage, page*ItemsPerPage + ItemsPerPage
}
