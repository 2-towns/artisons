package main

import (
	"artisons/admin"
	"artisons/admin/auth"
	"artisons/cache"
	"artisons/conf"
	"artisons/http/httperrors"
	"artisons/http/pages"
	"artisons/http/router"
	"artisons/http/security"
	"artisons/http/seo"
	"artisons/locales"
	"artisons/logs"
	"artisons/stats"
	"artisons/users"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

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
	r.Get("/products.html", admin.ProductList)
	r.Get("/products/add.html", admin.ProductForm)
	r.Get("/products/{id}/edit.html", admin.ProductForm)
	r.Get("/slug.html", admin.Slug)
	r.Get("/blog.html", admin.BlogList)
	r.Get("/blog/add.html", admin.BlogForm)
	r.Get("/blog/{id}/edit.html", admin.BlogForm)
	r.Get("/tags.html", admin.TagsList)
	r.Get("/tags/add.html", admin.TagsForm)
	r.Get("/tags/{id}/edit.html", admin.TagsForm)
	r.Get("/filters.html", admin.FiltersList)
	r.Get("/filters/add.html", admin.FiltersForm)
	r.Get("/filters/{id}/edit.html", admin.FiltersForm)
	r.Get("/orders.html", admin.OrdersList)
	r.Get("/orders/{id}/edit.html", admin.OrdersForm)
	r.Get("/settings.html", admin.SettingsForm)
	r.Get("/seo.html", admin.SeoList)
	r.Get("/seo/{id}/edit.html", admin.SeoForm)

	r.Post("/demo.html", stats.Demo)
	r.Post("/products/add.html", admin.ProductSave)
	r.Post("/products/{id}/edit.html", admin.ProductSave)
	r.Post("/products/{id}/delete.html", admin.ProductDelete)
	r.Post("/tags/add.html", admin.TagsSave)
	r.Post("/tags/{id}/edit.html", admin.TagsSave)
	r.Post("/tags/{id}/delete.html", admin.TagsDelete)
	r.Post("/filters/add.html", admin.FiltersSave)
	r.Post("/filters/{id}/edit.html", admin.FiltersSave)
	r.Post("/filters/{id}/delete.html", admin.FiltersDelete)
	r.Post("/blog/add.html", admin.BlogSave)
	r.Post("/blog/{id}/edit.html", admin.BlogSave)
	r.Post("/blog/{id}/delete.html", admin.BlogDelete)
	r.Post("/orders/{id}/status.html", admin.OrdersUpdateStatus)
	r.Post("/orders/{id}/note.html", admin.OrdersAddNote)
	r.Post("/contact-settings.html", admin.SettingsContactSave)
	r.Post("/shop-settings.html", admin.SettingsShopSave)
	r.Post("/locale.html", admin.EditLocale)
	r.Post("/seo/{id}/edit.html", admin.SeoSave)

	return r
}

func main() {
	locales.LoadEn()
	logs.Init()
	// security.LoadCsp()
	cache.Busting()

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

	log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!", seo.URLs["product"].URL)

	router.R.Handle("/public/*", http.StripPrefix("/public/", fs))
	router.R.Get("/", pages.Home)
	router.R.Get("/wish.html", pages.Wishes)
	router.R.Get("/blog.html", pages.Blog)
	router.R.Get("/blog/{slug}.html", pages.Article)
	router.R.Get(fmt.Sprintf("/%s/{slug}.html", seo.URLs["product"].URL), pages.Product)
	router.R.Get(fmt.Sprintf("/%s.html", seo.URLs["terms"].URL), pages.Static)
	router.R.Get(fmt.Sprintf("/%s.html", seo.URLs["about"].URL), pages.Static)
	router.R.Get("/sso.html", auth.Form)
	router.R.Get("/otp.html", auth.Form)
	router.R.Get("/addresses.html", pages.Addresses)

	router.R.Route("/account", func(r chi.Router) {
		r.Use(users.AccountOnly)
		r.Get("/index.html", pages.Account)
		r.Get("/orders.html", pages.Orders)
	})

	router.R.Post("/wish/{id}/add.html", pages.Wish)
	router.R.Post("/wish/{id}/delete.html", pages.UnWish)
	router.R.Post("/otp.html", auth.Otp)
	router.R.Post("/login.html", auth.Login)
	router.R.Post("/logout.html", auth.Logout)

	router.R.Mount("/admin", adminRouter())

	slog.LogAttrs(context.Background(), slog.LevelInfo, "starting server on addr", slog.String("addr", conf.ServerAddr))

	http.ListenAndServe(conf.ServerAddr, router.R)
}
