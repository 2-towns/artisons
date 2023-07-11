package main

import (
	"gifthub/locales"
	"gifthub/pages"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}

	locales.LoadEn()

	log.Println("Will look for views in", pwd+"/web/views/*")

	// p := message.NewPrinter(language.English)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(locales.Middleware)

	fs := http.FileServer(http.Dir("web/public"))
	router.Handle("/public/*", http.StripPrefix("/public/", fs))

	router.Get("/", pages.Home)

	http.ListenAndServe(":8080", router)

}
