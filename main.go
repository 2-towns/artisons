package main

import (
	"artisons/admin"
	"artisons/admin/auth"
	"artisons/cache"
	"artisons/carts"
	"artisons/conf"
	"artisons/http/flash"
	"artisons/http/forms"
	"artisons/http/htmx"
	"artisons/http/httperrors"
	"artisons/http/pages"
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
	"github.com/go-chi/httplog/v2"
)

var l httplog.Logger = httplog.Logger{
	Logger:  slog.Default(),
	Options: httplog.Options{},
}

func adminRouter() http.Handler {
	r := chi.NewRouter()
	r.NotFound(httperrors.NotFound)

	r.Use(users.Middleware)
	r.Use(users.AdminOnly)
	r.Use(seo.BlockRobots)

	r.Get("/index.html", admin.Dashboard)
	r.With(pages.Paginate).With(flash.Middleware).Get("/products.html", admin.ProductList)
	r.Get("/products/add.html", admin.ProductForm)
	r.Get("/products/{id}/edit.html", admin.ProductForm)
	r.Get("/slug.html", admin.Slug)
	r.With(pages.Paginate).With(flash.Middleware).Get("/blog.html", admin.BlogList)
	r.With(forms.ParseOptionalID).Get("/blog/add.html", admin.BlogForm)
	r.With(forms.ParseID).Get("/blog/{id}/edit.html", admin.BlogForm)
	r.With(pages.Paginate).With(flash.Middleware).Get("/tags.html", admin.TagList)
	r.Get("/tags/add.html", admin.TagForm)
	r.Get("/tags/{id}/edit.html", admin.TagForm)
	r.With(pages.Paginate).With(flash.Middleware).Get("/filters.html", admin.FilterList)
	r.Get("/filters/add.html", admin.FilterForm)
	r.Get("/filters/{id}/edit.html", admin.FilterForm)
	r.With(pages.Paginate).With(flash.Middleware).Get("/orders.html", admin.OrderList)
	r.Get("/orders/{id}/edit.html", admin.OrderForm)
	r.Get("/settings.html", admin.SettingsForm)
	r.With(pages.Paginate).With(flash.Middleware).Get("/seo.html", admin.SeoList)
	r.Get("/seo/{id}/edit.html", admin.SeoForm)

	r.With(security.Csrf).Post("/demo.html", stats.Demo)
	r.With(security.Csrf).With(forms.ParseMultipartForm).Post("/products/add.html", admin.ProductSave)
	r.With(security.Csrf).With(forms.ParseMultipartForm).Post("/products/{id}/edit.html", admin.ProductSave)
	r.With(security.Csrf).With(forms.ParseForm).With(pages.Paginate).Post("/products/{id}/delete.html", admin.ProductDelete)
	r.With(security.Csrf).With(forms.ParseMultipartForm).Post("/tags/add.html", admin.TagSave)
	r.With(security.Csrf).With(forms.ParseMultipartForm).Post("/tags/{id}/edit.html", admin.TagSave)
	r.With(security.Csrf).With(forms.ParseForm).With(pages.Paginate).Post("/tags/{id}/delete.html", admin.TagDelete)
	r.With(security.Csrf).With(forms.ParseForm).Post("/filters/add.html", admin.FilterSave)
	r.With(security.Csrf).With(forms.ParseForm).Post("/filters/{id}/edit.html", admin.FilterSave)
	r.With(security.Csrf).With(forms.ParseForm).With(pages.Paginate).Post("/filters/{id}/delete.html", admin.FilterDelete)
	r.With(security.Csrf).With(forms.ParseOptionalID).With(forms.ParseMultipartForm).Post("/blog/add.html", admin.BlogSave)
	r.With(security.Csrf).With(forms.ParseID).With(forms.ParseMultipartForm).Post("/blog/{id}/edit.html", admin.BlogSave)
	r.With(security.Csrf).With(forms.ParseID).With(forms.ParseForm, pages.Paginate).Post("/blog/{id}/delete.html", admin.BlogDelete)
	r.With(security.Csrf).With(forms.ParseForm).Post("/orders/{id}/status.html", admin.OrderUpdateStatus)
	r.With(security.Csrf).With(forms.ParseForm).Post("/orders/{id}/note.html", admin.OrderAddNote)
	r.With(security.Csrf).With(forms.ParseMultipartForm).Post("/contact-settings.html", admin.SettingsContactSave)
	r.With(security.Csrf).With(forms.ParseForm).Post("/shop-settings.html", admin.SettingsShopSave)
	r.With(security.Csrf).Post("/locale.html", admin.EditLocale)
	r.With(security.Csrf).With(forms.ParseForm).Post("/seo/{id}/edit.html", admin.SeoSave)

	return r
}

func main() {
	locales.LoadEn()
	logs.Init()
	// security.LoadCsp()
	cache.Busting()

	var router = chi.NewRouter()

	fs := http.FileServer(http.Dir("web/public"))
	router.Handle("/public/*", http.StripPrefix("/public/", fs))

	router.Group(func(r chi.Router) {
		// router.Use(middleware.RealIP)
		r.Use(httplog.RequestLogger(&l))
		r.Use(security.Headers)
		r.Use(users.Domain)
		r.Use(htmx.Middleware)
		r.Use(locales.Middleware)
		r.Use(carts.Middleware)

		if conf.Debug {
			r.Use(seo.BlockRobots)
		}

		r.Mount("/admin", adminRouter())

		log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!", fmt.Sprintf("/%s.html", seo.URLs["terms"].URL))
		r.With(stats.Middleware).Get("/", pages.Home)
		r.With(stats.Middleware).Get("/blog.html", pages.Blog)
		r.With(stats.Middleware).Get("/blog/{slug}.html", pages.Article)
		r.With(stats.Middleware).Get(fmt.Sprintf("/%s/{slug}.html", seo.URLs["product"].URL), pages.Product)
		r.With(stats.Middleware).Get(fmt.Sprintf("/%s.html", seo.URLs["terms"].URL), pages.Static)
		r.With(stats.Middleware).Get(fmt.Sprintf("/%s.html", seo.URLs["about"].URL), pages.Static)
		r.With(stats.Middleware).Get(fmt.Sprintf("/%s.html", seo.URLs["categories"].URL), pages.Categories)
		r.With(stats.Middleware).Get("/addresses.html", pages.Addresses)
		r.With(stats.Middleware).Get("/search.html", pages.Search)
		r.With(stats.Middleware).Get("/cart.html", pages.Cart)
		r.With(users.Middleware).Get("/sso.html", auth.Form)
		r.With(users.Middleware).With(stats.Middleware).Get("/otp.html", auth.Form)

		r.With(security.Csrf).Post("/cart/{id}/add.html", pages.CartAdd)
		r.With(security.Csrf).Post("/otp.html", auth.Otp)
		r.With(security.Csrf).Post("/login.html", auth.Login)
		r.With(security.Csrf).Post("/logout.html", auth.Logout)

		r.Route("/account", func(r chi.Router) {
			r.Use(users.Middleware)
			r.Use(users.AccountOnly)

			r.Get("/index.html", pages.Account)
			r.Get("/address.html", pages.AddressForm)
			r.Get("/orders.html", pages.Orders)
			r.Get("/orders/{id}/detail.html", pages.Orders)
			r.With(stats.Middleware).Get("/wish.html", pages.Wishes)

			r.With(security.Csrf).Post("/address.html", pages.Address)
			r.With(security.Csrf).Post("/wish/{id}/add.html", pages.Wish)
			r.With(security.Csrf).Post("/wish/{id}/delete.html", pages.UnWish)
		})
	})

	slog.LogAttrs(context.Background(), slog.LevelInfo, "starting server on addr", slog.String("addr", conf.ServerAddr))

	http.ListenAndServe(conf.ServerAddr, router)
}
