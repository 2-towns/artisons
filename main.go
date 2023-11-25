package main

import (
	"context"
	"fmt"
	"gifthub/admin"
	"gifthub/admin/login"
	"gifthub/admin/urls"
	"gifthub/conf"
	"gifthub/http/httperrors"
	"gifthub/http/security"
	"gifthub/http/seo"
	"gifthub/locales"
	"gifthub/logs"
	"gifthub/pages"
	"gifthub/users"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func get(router chi.Router, pathname string, handlerFn http.HandlerFunc) {
	for _, lang := range conf.Languages {
		l := language.Make(lang)
		p := message.NewPrinter(l)
		url := p.Sprint(pathname)
		router.Get(url, handlerFn)
	}
}

func post(router chi.Router, pathname string, handlerFn http.HandlerFunc) {
	for _, lang := range conf.Languages {
		l := language.Make(lang)
		p := message.NewPrinter(l)
		url := p.Sprint(pathname)
		router.Post(url, handlerFn)
	}
}

func adminRouter() http.Handler {
	r := chi.NewRouter()
	r.NotFound(httperrors.NotFound)
	r.Use(seo.BlockRobots)
	r.Use(users.AdminOnly)
	r.Use(security.Csrf)

	r.Get("/", admin.Dashboard)

	return r
}

func main() {
	locales.LoadEn()
	logs.Init()
	security.LoadCsp()

	logger := httplog.NewLogger("http", httplog.Options{
		LogLevel: slog.LevelDebug,
		// Concise:  true,
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

	if conf.Debug {
		router.Use(seo.BlockRobots)
	}

	fs := http.FileServer(http.Dir("web/public"))
	router.Handle("/public/*", http.StripPrefix("/public/", fs))

	router.Get("/", pages.Home)
	router.Get(urls.AuthPrefix, login.Form)
	router.Post(urls.Otp, login.Otp)
	router.Post(urls.Login, login.Login)

	router.Mount(urls.AdminPrefix, adminRouter())

	slog.LogAttrs(context.Background(), slog.LevelInfo, "starting server on addr", slog.String("addr", conf.ServerAddr))

	http.ListenAndServe(conf.ServerAddr, router)
}

// Add copy event for otp
// add back button
// remove otp label
// add cache
