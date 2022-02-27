package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/dev-ayaa/resvbooking/pkg/config"
	"github.com/dev-ayaa/resvbooking/pkg/handlers"
	"github.com/dev-ayaa/resvbooking/pkg/models"
	"github.com/dev-ayaa/resvbooking/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

// the most likely place
// to use session is the handlers package
var app config.AppConfig
var session *scs.SessionManager

func main() {

	err := run()
	if err != nil {
		log.Fatal("Failed to run the Application........")
	}
	//Using session to keep track of data store from the form

	fmt.Println("Starting the Server :8080")

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)

}

func run() error {

	gob.Register(models.ReservationData{})

	app.InProduction = false

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

	//Referencing the map store in the app AppConfig
	repo := handlers.NewRepository(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	return nil

}
