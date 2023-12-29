package main

import (
	"context"
	"gifthub/admin"
	"gifthub/admin/login"
	"gifthub/admin/urls"
	"gifthub/cache"
	"gifthub/conf"
	"gifthub/http/httperrors"
	"gifthub/http/security"
	"gifthub/http/seo"
	"gifthub/locales"
	"gifthub/logs"
	"gifthub/pages"
	"gifthub/stats"
	"gifthub/users"
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

	r.Get(urls.Map["dashboard"], admin.Dashboard)
	r.Get(urls.Map["products"], admin.Products)
	r.Get("/products/add.html", admin.AddProductForm)
	r.Get("/products/{id}/edit.html", admin.EditProductForm)

	r.Post(urls.Map["demo"], stats.Demo)
	r.Post("/products/add.html", admin.AddProduct)
	r.Post("/products/{id}/edit.html", admin.EditProduct)
	r.Post("/products/{id}/delete.html", admin.DeleteProduct)

	return r
}

func main() {
	locales.LoadEn()
	logs.Init()
	security.LoadCsp()
	cache.Busting()

	l := httplog.Logger{
		Logger:  slog.Default(),
		Options: httplog.Options{},
	}

	router := chi.NewRouter()

	// router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(httplog.RequestLogger(&l))
	router.Use(locales.Middleware)
	router.Use(security.Csrf)
	router.Use(security.Headers)
	router.Use(users.Middleware)
	router.Use(stats.Middleware)

	if conf.Debug {
		router.Use(seo.BlockRobots)
	}

	fs := http.FileServer(http.Dir("web/public"))
	router.Handle("/public/*", http.StripPrefix("/public/", fs))

	router.Get("/", pages.Home)
	router.Get(urls.Map["auth"], login.Form)
	router.Post(urls.Map["auth_otp"], login.Otp)
	router.Post(urls.Map["auth_login"], login.Login)

	router.Mount(urls.AdminPrefix, adminRouter())

	slog.LogAttrs(context.Background(), slog.LevelInfo, "starting server on addr", slog.String("addr", conf.ServerAddr))

	http.ListenAndServe(conf.ServerAddr, router)
}

// searchin input error bien positioné + disparaitree
