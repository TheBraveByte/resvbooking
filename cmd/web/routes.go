package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/dev-ayaa/resvbooking/pkg/config"
	"github.com/dev-ayaa/resvbooking/pkg/handlers"
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
	mux.Get("/book-room-now", handlers.Repo.BookRoomNow)
	mux.Get("/login", handlers.Repo.LoginPage)
	mux.Post("/login", handlers.Repo.PostLoginPage)
	mux.Get("/logout", handlers.Repo.LogOutPage)

	//setting up the admin page

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Authenticate)
		mux.Get("/dashboard", handlers.Repo.AdminPage)
		mux.Get("/admin-new-reservation", handlers.Repo.AdminNewReservation)
		mux.Get("/admin-all-reservation", handlers.Repo.AdminAllReservation)
		mux.Get("/admin-reservation-calendar", handlers.Repo.AdminReservationCalendar)
		mux.Post("/admin-reservation-calendar", handlers.Repo.PostAdminReservationCalendar)

		mux.Get("/admin-show-reservation/{src}/{id}/show", handlers.Repo.AdminShowReservation)
		mux.Post("/admin-show-reservation/{src}/{id}", handlers.Repo.PostAdminShowReservation)

		mux.Get("/admin/admin-delete-reservation/{src}/{id}/done", handlers.Repo.AdminDeleteReservation)
		mux.Get("/admin/admin-process-reservation/{src}/{id}/done", handlers.Repo.AdminProcessReservation)

	})
	//This allows files static files like images and icon to display in the html/tmpl
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
