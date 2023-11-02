package main

import (
	"gifthub/conf"
	"gifthub/locales"
	"gifthub/logs"
	"gifthub/pages"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	locales.LoadEn()
	logs.Init()
	// p := message.NewPrinter(language.English)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(locales.Middleware)

	fs := http.FileServer(http.Dir("web/public"))
	router.Handle("/public/*", http.StripPrefix("/public/", fs))

	slog.Info(conf.ImgProxyPath)

	router.Get("/", pages.Home)

	http.ListenAndServe(":8080", router)

}
