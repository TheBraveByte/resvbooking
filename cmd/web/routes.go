package main

import (
	"github.com/Akinleye007/resvbooking/pkg/config"
	"github.com/Akinleye007/resvbooking/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	//	mux := pat.New()
	//	mux.Get("/", http.HandlerFunc(handlers.Repo.HomePage))
	//	mux.Get("/about", http.HandlerFunc(handlers.Repo.AboutPage))
	//
	//	return mux
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.HomePage)
	mux.Get("/about", handlers.Repo.AboutPage)

	//This allows files static files like images and icon to display in the html
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*",http.StripPrefix("/static",fileServer))

	return mux

}
