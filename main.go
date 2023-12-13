package main

import (
	"context"
	"fmt"
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

	r.Post(urls.Map["demo"], stats.Demo)

	return r
}

func main() {
	locales.LoadEn()
	logs.Init()
	security.LoadCsp()
	cache.Busting()

	logger := httplog.NewLogger("http", httplog.Options{
		LogLevel: slog.LevelDebug,
		Concise:  true,
		// RequestHeaders:   true,
		// TimeFieldFormat:  time.RFC850,
		Tags: map[string]string{
			"debug": fmt.Sprintf("%t", conf.Debug),
		},
	})

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(httplog.RequestLogger(logger))
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

// migration => migrate
// calcule md5 js
// add cache
