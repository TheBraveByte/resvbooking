package main

import (
	"encoding/gob"
	"fmt"
	"github.com/Akinleye007/resvbooking/pkg/models"
	"log"
	"net/http"
	"time"

	"github.com/Akinleye007/resvbooking/pkg/config"
	"github.com/Akinleye007/resvbooking/pkg/handlers"
	"github.com/Akinleye007/resvbooking/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

// the most likely place
// to use session is the handlers package
var app config.AppConfig
var session *scs.SessionManager

func main() {

	app.InProduction = false
	//Using session to keep track of data store from the form
	gob.Register(models.ReservationData{})

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	//Getting the templates cache
	tc, err := render.TemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	// storing the cache in the app config
	app.TempCache = tc
	app.Session = session

	//authorize using cache
	app.UseCache = false

	//Referencing the map store in the app AppConfig
	repo := handlers.NewRepository(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println("Starting the Server :8080")

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)

}
