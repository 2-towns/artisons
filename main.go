package main

import (
	"gifthub/conf"
	"gifthub/locales"
	"gifthub/logs"
	"gifthub/pages"
	"gifthub/users"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

func adminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(users.AdminOnly)
	r.Get("/", pages.Home)

	return r
}

func main() {
	locales.LoadEn()
	logs.Init()

	logger := httplog.NewLogger("httplog-example", httplog.Options{
		// JSON:             true,
		LogLevel: slog.LevelDebug,
		Concise:  true,
		// RequestHeaders:   true,
		MessageFieldName: "message",
		TimeFieldFormat:  time.RFC850,
		Tags: map[string]string{
			"version": "v1.0-81aa4244d9fc8076a",
			"env":     "dev",
		},
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(httplog.RequestLogger(logger))
	router.Use(locales.Middleware)

	fs := http.FileServer(http.Dir("web/public"))
	router.Handle("/public/*", http.StripPrefix("/public/", fs))

	router.Get("/", pages.Home)
	router.Mount(conf.AdminPrefix, adminRouter())

	http.ListenAndServe(":8080", router)

}
