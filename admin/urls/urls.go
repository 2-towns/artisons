package urls

const (
	AdminPrefix = "/admin"
	AuthPrefix  = "/auth"
)

var Map = map[string]string{
	"auth":            AuthPrefix,
	"auth_logout":     AuthPrefix + "/logout.html",
	"auth_login":      AuthPrefix + "/login.html",
	"auth_otp":        AuthPrefix + "/otp.html",
	"admin_dashboard": AdminPrefix,
	"dashboard":       "/",
	"admin_demo":      AdminPrefix + "/demo.html",
	"demo":            "/demo.html",
	"admin_products":  AdminPrefix + "/products.html",
	"products":        "/products.html",
}

func Get(name string) string {
	return Map[name]
}
