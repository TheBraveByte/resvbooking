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
	mux.Get("/contact", handlers.Repo.ContactPage)
	mux.Get("/junior-suite", handlers.Repo.JuniorSuitePage)
	mux.Get("/premium-suite", handlers.Repo.PremiumSuitePage)
	mux.Get("/deluxe-suite", handlers.Repo.DeluxeSuitePage)
	mux.Get("/penthouse-suite", handlers.Repo.PenthousePage)
	mux.Get("/executive-suite", handlers.Repo.ExecutivePage)
	mux.Get("/make-reservation", handlers.Repo.MakeReservationPage)

	mux.Get("/check-availability", handlers.Repo.CheckAvailabilityPage)

	mux.Post("/check-availability", handlers.Repo.PostCheckAvailabilityPage)

	//This allows files static files like images and icon to display in the html
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux

}
