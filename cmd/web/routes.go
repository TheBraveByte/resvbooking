package main

import (
	"github.com/dev-ayaa/resvbooking/pkg/config"
	"github.com/dev-ayaa/resvbooking/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
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

	mux.Get("/check-availability", handlers.Repo.CheckAvailabilityPage)
	mux.Post("/check-availability", handlers.Repo.PostCheckAvailabilityPage)

	mux.Get("/json-availability", handlers.Repo.JsonAvailabilityPage)
	mux.Post("/json-availability", handlers.Repo.JsonAvailabilityPage)

	mux.Get("/select-available-room/{id}", handlers.Repo.SelectAvailableRoom)

	mux.Get("/make-reservation", handlers.Repo.MakeReservationPage)
	mux.Post("/make-reservation", handlers.Repo.PostMakeReservationPage)
	mux.Get("/make-reservation-data", handlers.Repo.MakeReservationSummary)
	mux.Get("/book-now", handlers.Repo.BookRoomNow)

	//This allows files static files like images and icon to display in the html
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
