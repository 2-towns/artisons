package urls

const (
	AdminPrefix = "/admin"
	AuthPrefix  = "/auth"
	A           = ""
)

var Map = map[string]string{
	"auth":                  AuthPrefix,
	"auth_logout":           AuthPrefix + "/logout.html",
	"auth_login":            AuthPrefix + "/login.html",
	"auth_otp":              AuthPrefix + "/otp.html",
	"admin":                 AdminPrefix,
	"admin_dashboard":       AdminPrefix,
	"dashboard":             "/",
	"admin_demo":            AdminPrefix + "/demo.html",
	"demo":                  "/demo.html",
	"admin_products":        AdminPrefix + "/products.html",
	"products":              "/products.html",
	"admin_edit_products":   AdminPrefix + "/products/:id/edit.html",
	"edit_products":         "/products/:id/edit.html",
	"admin_delete_products": AdminPrefix + "/products/:id/delete.html",
	"delete_products":       "/products/:id/delete.html",
}

func Get(name string) string {
	return Map[name]
}
