package main

import (
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/julienschmidt/httprouter"

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
	r.GET("/applications/:id", ac.Get)
	r.POST("/applications", ac.Create)
	r.DELETE("/applications/:id", ac.Delete)

	r.GET("/releases", rc.GetAll)
	r.GET("/releases/:id", rc.Get)
	r.POST("/releases", rc.Create)
	r.DELETE("/releases/:id", rc.Delete)

	log.Fatal(http.ListenAndServe(url, r))
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
