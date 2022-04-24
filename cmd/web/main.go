package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/dev-ayaa/resvbooking/pkg/config"
	"github.com/dev-ayaa/resvbooking/pkg/driver"
	"github.com/dev-ayaa/resvbooking/pkg/handlers"
	"github.com/dev-ayaa/resvbooking/pkg/helpers"
	"github.com/dev-ayaa/resvbooking/pkg/models"
	"github.com/dev-ayaa/resvbooking/pkg/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

// the most likely place
// to use session is the handlers package
var app config.AppConfig
var session *scs.SessionManager
var infoLogger *log.Logger
var errorLogger *log.Logger
var dataSourceName string

func main() {

	db, err := run()
	if err != nil {
		log.Fatal("Failed to run the Application........")
	}

	//Close the Database
	defer func(db *driver.DB) {
		_ = db.PSQL.Close()
	}(db)

	fmt.Println("..............connecting the mail server in channels..............")

	defer close(app.MailChannel)

	fmt.Println("Starting the Server :8080")

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)

}

func run() (*driver.DB, error) {
	//Using session to keep track of data store in Models

	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.User{})
	gob.Register(models.RoomRestriction{})
	//gob.Register([]models.Room{})

	mailChannel := make(chan models.MailData)
	app.MailChannel = mailChannel
	mailRoutes()

	app.InProduction = false

	infoLogger = log.New(os.Stdout, "INFO ::\t", log.LstdFlags)
	app.InfoLog = infoLogger

	errorLogger = log.New(os.Stdout, "ERROR ::\t", log.LstdFlags|log.Lshortfile)
	app.ErrorLog = errorLogger

	session = scs.New()
	session.Lifetime = 24 * time.Hour              // how to keep the session of users
	session.Cookie.Persist = true                  //To keep cookies
	session.Cookie.SameSite = http.SameSiteLaxMode //if the user visit the same sites again
	session.Cookie.Secure = app.InProduction       // is the application in production or development

	//Getting the templates cache
	tc, err := render.TemplateCache()
	fmt.Println(tc, err)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Cannot create template cache")
	}

	// storing the cache in the app config
	app.TempCache = tc
	app.Session = session

	//authorize using cache
	app.UseCache = false

	//connection the database to the Application
	log.Println(".........Connecting to the database.........")
	dataSourceName = "host=localhost port=5432  dbname=Resvbooking user=postgres password=dev-ayaa"
	db, err := driver.ConnectSqlDb(dataSourceName)

	if err != nil {
		log.Fatal("Error Connecting to the database.....")
		return nil, err
	}

	//Referencing the map store in the app AppConfig
	repo := handlers.NewRepository(&app, db)
	handlers.NewHandlers(repo)
	helpers.NewHelper(&app)

	render.NewTemplates(&app)

	return db, nil

}
