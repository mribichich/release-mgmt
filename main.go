package main

import (
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"github.com/mribichich/release-mgmt/controllers"
)

func main() {
	url := ":8080"

	log.Printf("listening at '%s' ...", url)

	// router := NewRouter()

	r := httprouter.New()

	session := getSession()
	ac := controllers.NewApplicationsController(session)
	rc := controllers.NewReleasesController(session)

	r.GET("/applications", ac.GetAll)
	r.GET("/applications/:name", ac.Get)
	r.POST("/applications", ac.Create)
	r.PUT("/applications/:name", ac.Update)
	r.DELETE("/applications/:name", ac.Delete)

	r.GET("/applications/:name/releases", rc.GetAll)
	r.GET("/applications/:name/releases/:version", rc.Get)
	r.POST("/applications/:name/releases", rc.Create)
	r.PUT("/applications/:name/releases/:version", rc.Update)
	r.DELETE("/applications/:name/releases/:version", rc.Delete)

	handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(url, handler))
}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}

	// Deliver session
	return s
}
