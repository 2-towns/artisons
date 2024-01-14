package main

import (
	"context"
	"gifthub/admin"
	"gifthub/admin/login"
	"gifthub/cache"
	"gifthub/conf"
	"gifthub/http/httperrors"
	"gifthub/http/router"
	"gifthub/http/security"
	"gifthub/http/seo"
	"gifthub/locales"
	"gifthub/logs"
	"gifthub/pages"
	"gifthub/stats"
	"gifthub/users"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

func adminRouter() http.Handler {
	r := chi.NewRouter()
	r.NotFound(httperrors.NotFound)
	r.Use(seo.BlockRobots)
	r.Use(users.AdminOnly)
	r.Use(security.Csrf)

	r.Get("/index.html", admin.Dashboard)
	r.Get("/products.html", admin.Products)
	r.Get("/products/add.html", admin.AddProductForm)
	r.Get("/products/{id}/edit.html", admin.EditProductForm)
	r.Get("/blog.html", admin.Blog)
	r.Get("/blog/add.html", admin.AddBlogForm)
	r.Get("/blog/{id}/edit.html", admin.EditBlogForm)
	r.Get("/orders.html", admin.Orders)
	r.Get("/orders/{id}/edit.html", admin.EditOrderForm)
	r.Get("/settings.html", admin.SettingsForm)
	r.Get("/products.html", admin.Products)
	r.Get("/seo.html", admin.Seo)
	r.Get("/seo/{id}/edit.html", admin.EditSeoForm)

	r.Post("/demo.html", stats.Demo)
	r.Post("/products/add.html", admin.AddProduct)
	r.Post("/products/{id}/edit.html", admin.EditProduct)
	r.Post("/products/{id}/delete.html", admin.DeleteProduct)
	r.Post("/blog/add.html", admin.AddBlog)
	r.Post("/blog/{id}/edit.html", admin.EditBlog)
	r.Post("/blog/{id}/delete.html", admin.DeleteBlog)
	r.Post("/orders/{id}/status.html", admin.UpdateOrderStatus)
	r.Post("/orders/{id}/note.html", admin.AddOrderNote)
	r.Post("/contact-settings.html", admin.EditContactSettings)
	r.Post("/shop-settings.html", admin.EditShopSettings)
	r.Post("/locale.html", admin.EditLocale)
	r.Post("/seo/{id}/edit.html", admin.EditSeo)

	return r
}

func main() {
	locales.LoadEn()
	logs.Init()
	security.LoadCsp()
	cache.Busting()
	log.Println(time.Now().Unix())
	l := httplog.Logger{
		Logger:  slog.Default(),
		Options: httplog.Options{},
	}

	// router.Use(middleware.RequestID)
	router.R.Use(middleware.RealIP)
	router.R.Use(httplog.RequestLogger(&l))
	router.R.Use(locales.Middleware)
	router.R.Use(security.Csrf)
	router.R.Use(security.Headers)
	router.R.Use(users.Middleware)
	router.R.Use(stats.Middleware)

	if conf.Debug {
		router.R.Use(seo.BlockRobots)
	}

	fs := http.FileServer(http.Dir("web/public"))
	router.R.Handle("/public/*", http.StripPrefix("/public/", fs))

	router.R.Get("/", pages.Home)
	router.R.Get("/auth/index.html", login.Form)
	router.R.Post("/auth/otp.html", login.Otp)
	router.R.Post("/auth/login.html", login.Login)

	router.R.Mount("/admin", adminRouter())

	slog.LogAttrs(context.Background(), slog.LevelInfo, "starting server on addr", slog.String("addr", conf.ServerAddr))

	http.ListenAndServe(conf.ServerAddr, router.R)
}

// fix url image
// Change the blog editor to display paragraph
