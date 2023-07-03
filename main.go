package main

import (
	"log"
	"mustafir/locales"
	"mustafir/middlewares"
	"mustafir/routes"
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

	log.Println("Will look for views in", pwd+"/views/*")

	// p := message.NewPrinter(language.English)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middlewares.Lang)

	fs := http.FileServer(http.Dir("public"))
	router.Handle("/public/*", http.StripPrefix("/public/", fs))

	router.Get("/", routes.HomeRoute)

	http.ListenAndServe(":8080", router)

}
