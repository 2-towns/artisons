package main

import (
	"gifthub/conf"
	"gifthub/locales"
	"gifthub/pages"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	locales.LoadEn()

	// p := message.NewPrinter(language.English)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(locales.Middleware)

	fs := http.FileServer(http.Dir("web/public"))
	router.Handle("/public/*", http.StripPrefix("/public/", fs))

	log.Println(conf.ImgProxyPath)

	router.Get("/", pages.Home)

	http.ListenAndServe(":8080", router)

}
