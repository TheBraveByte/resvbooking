package main

import (
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

// var app *config.AppConfig

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
		fmt.Println("Hit the page")
		next.ServeHTTP(wr, rq)
	})
}

//NoSurf add csrf Protection to all POST request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//SessionLoad Loads and save the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
