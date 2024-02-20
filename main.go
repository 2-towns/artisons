package main

import (
	"artisons/addresses"
	"artisons/auth"
	"artisons/blog"
	"artisons/cache"
	"artisons/carts"
	"artisons/conf"
	"artisons/http/contexts"
	"artisons/http/security"
	"artisons/locales"
	"artisons/logs"
	"artisons/orders"
	"artisons/products"
	"artisons/products/filters"
	"artisons/seo"
	"artisons/seo/urls"
	"artisons/shops"
	"artisons/shops/website"
	"artisons/stats"
	"artisons/string/slughttp"
	"artisons/tags"
	"artisons/users"
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		ctx := context.WithValue(r.Context(), contexts.RequestID, id.String())
		ctx = context.WithValue(ctx, contexts.Locale, conf.DefaultLocale)
		ctx = context.WithValue(ctx, contexts.HX, r.Header.Get("HX-Request") == "true")
		ctx = context.WithValue(ctx, contexts.Tracking, conf.EnableTrackingLog)
		ctx = context.WithValue(ctx, contexts.ThrowsWhenPaymentFailed, shops.Data.ThrowsWhenPaymentFailed)

		slog.LogAttrs(
			ctx,
			slog.LevelInfo,
			"new request",
			slog.String("path", r.URL.Path),
			slog.String("ua", r.Header.Get("User-Agent")),
			slog.String("referer", r.Header.Get("Referer")),
		)

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Powered-By", "WordPress")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set(
			"Accept-CH",
			"Sec-CH-Prefers-Color-Scheme, Device-Memory, Downlink, ECT",
		)
		w.Header().Set("Referrer-Policy", "strict-origin")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set(
			"Strict-Transport-Security",
			"max-age=63072000; includeSubDomains; preload",
		)
		w.Header().Set("Date", time.Now().Format(time.RFC1123))
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("X-XSS-Protection", "1")

		if conf.Debug {
			w.Header().Set("X-Robots-Tag", "noindex")
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func adminMux() *http.ServeMux {
	admin := http.NewServeMux()
	admin.HandleFunc("GET /admin/index", stats.Handler)
	admin.HandleFunc("GET /admin/products", products.AdminListHandlerHandler)
	admin.HandleFunc("GET /admin/products/add", products.AdminFormHandler)
	admin.HandleFunc("GET /admin/products/{id}/edit", products.AdminFormHandler)
	admin.HandleFunc("GET /admin/slug", slughttp.Handler)
	admin.HandleFunc("GET /admin/blog", blog.AdminListHandler)
	admin.HandleFunc("GET /admin/blog/add", blog.AdminFormHandler)
	admin.HandleFunc("GET /admin/blog/{id}/edit", blog.AdminFormHandler)
	admin.HandleFunc("GET /admin/tags", tags.AdminListHandler)
	admin.HandleFunc("GET /admin/tags/add", tags.AdminFormHandler)
	admin.HandleFunc("GET /admin/tags/{id}/edit", tags.AdminFormHandler)
	admin.HandleFunc("GET /admin/filters", filters.AdminListHandler)
	admin.HandleFunc("GET /admin/filters/add", filters.AdminFormHandler)
	admin.HandleFunc("GET /admin/filters/{id}/edit", filters.AdminFormHandler)
	admin.HandleFunc("GET /admin/orders", orders.OrderListHandler)
	admin.HandleFunc("GET /admin/orders/{id}/edit", orders.OrderFormHandler)
	admin.HandleFunc("GET /admin/settings", shops.SettingsFormHandler)
	admin.HandleFunc("GET /admin/seo", seo.AdminListHandler)
	admin.HandleFunc("GET /admin/seo/{id}/edit", seo.AdminFormHandler)
	admin.HandleFunc("POST /admin/demo", stats.DemoHandler)
	admin.HandleFunc("POST /admin/products/add", products.AdminSaveHandler)
	admin.HandleFunc("POST /admin/products/{id}/edit", products.AdminSaveHandler)
	admin.HandleFunc("POST /admin/products/{id}/delete", products.AdminDeleteHandler)
	admin.HandleFunc("POST /admin/tags/add", tags.AdminSaveHandler)
	admin.HandleFunc("POST /admin/tags/{id}/edit", tags.AdminSaveHandler)
	admin.HandleFunc("POST /admin/tags/{id}/delete", tags.AdminDeleteHandler)
	admin.HandleFunc("POST /admin/filters/add", filters.AdminSaveHandler)
	admin.HandleFunc("POST /admin/filters/{id}/edit", filters.AdminSaveHandler)
	admin.HandleFunc("POST /admin/filters/{id}/delete", filters.AdminDeleteHandler)
	admin.HandleFunc("POST /admin/blog/add", blog.AdminSaveHandler)
	admin.HandleFunc("POST /admin/blog/{id}/edit", blog.AdminSaveHandler)
	admin.HandleFunc("POST /admin/blog/{id}/delete", blog.AdminDeleteHandler)
	admin.HandleFunc("POST /admin/orders/{id}/status", orders.OrderUpdateStatus)
	admin.HandleFunc("POST /admin/orders/{id}/note", orders.OrderAddNoteHandler)
	admin.HandleFunc("POST /admin/contact-settings", shops.SettingsContactSave)
	admin.HandleFunc("POST /admin/shop-settings", shops.SettingsShopSave)
	admin.HandleFunc("POST /admin/seo/{id}/edit", seo.AdminSaveHandler)
	// csrf.With(forms.ParseForm).Post("/locale", admin.EditLocale)

	return admin
}

func websiteMux() *http.ServeMux {
	web := http.NewServeMux()
	web.HandleFunc("GET /", website.Home)
	web.HandleFunc("GET /blog", blog.ListHandler)
	web.HandleFunc("GET /blog/{slug}", blog.ArticleHandler)
	web.HandleFunc("GET /cart", carts.Handler)
	web.HandleFunc("GET /otp", auth.Formhandler)
	web.HandleFunc("GET /search", website.SearchHandler)
	web.HandleFunc("GET /"+urls.Get("product", "url")+"/{slug}", products.ProductHandler)
	web.HandleFunc("GET /"+urls.Get("terms", "url"), website.StaticHandler)
	web.HandleFunc("GET /"+urls.Get("about", "url"), website.StaticHandler)
	web.HandleFunc("GET /"+urls.Get("categories", "url"), website.CategoriesHandler)

	return web
}

func accountMux() *http.ServeMux {
	stat := http.NewServeMux()
	stat.HandleFunc("GET /account/index", users.AccountHandler)
	stat.HandleFunc("GET /account/address", users.AddressFormHandler)
	stat.HandleFunc("GET /account/wish", products.WishesHandler)

	account := http.NewServeMux()
	account.Handle("GET /", stats.Middleware(stat))
	account.HandleFunc("GET /account/orders", orders.OrdersHandler)
	account.HandleFunc("GET /account/orders/{id}/detail", orders.OrderHandler)
	account.HandleFunc("POST /account/address", users.AddressHandler)
	account.HandleFunc("POST /account/wish/{id}/add", products.WishHandler)
	account.HandleFunc("POST /account/wish/{id}/delete", products.UnWishHandler)

	return account
}

func main() {
	locales.LoadEn()
	logs.Init()
	// security.LoadCsp()
	cache.Busting()

	admin := adminMux()
	web := websiteMux()
	account := accountMux()

	app := http.NewServeMux()
	app.Handle("GET /admin/", users.AdminOnly(admin))
	app.Handle("GET /account/", users.AccountOnly(account))
	app.Handle("GET /", stats.Middleware(web))
	app.Handle("POST /admin/", users.AdminOnly(admin))
	app.Handle("POST /account/", users.AccountOnly(account))
	app.HandleFunc("GET /sso", auth.Formhandler)
	app.HandleFunc("GET /addresses", addresses.Handler)
	app.HandleFunc("GET /delivery", carts.DeliveryHandler)
	app.HandleFunc("GET /cart/address", carts.AddressFormHandler)
	app.HandleFunc("GET /payment", carts.PaymentHandler)
	app.HandleFunc("POST /payment", carts.PaymentProcessHandler)
	app.HandleFunc("POST /cart/address", carts.AddressHandler)
	app.HandleFunc("POST /otp", auth.OtpHandler)
	app.HandleFunc("POST /login", auth.LoginHandler)
	app.HandleFunc("POST /logout", auth.LogoutHandler)
	app.HandleFunc("POST /cart/{id}/add", carts.AddHandler)
	app.HandleFunc("POST /cart/{id}/delete", carts.DeleteHandler)
	app.HandleFunc("POST /delivery", carts.DeliverySetHandler)

	fs := http.FileServer(http.Dir("web/public"))
	mux := http.NewServeMux()
	mux.Handle("GET /public/", fs)
	mux.Handle("GET /css/", fs)
	mux.Handle("GET /js/", fs)
	mux.Handle("GET /icons/", fs)
	mux.Handle("GET /fonts/", fs)
	mux.Handle("GET /favicon.ico", fs)

	mux.Handle("GET /", handler(app))
	mux.Handle("POST /", handler(security.Csrf(app)))

	slog.Info("Starting server on", slog.String("address", conf.ServerAddr))
	err := http.ListenAndServe(conf.ServerAddr, mux)
	log.Fatal(err)
}
