package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/mribichich/release-mgmt/entities"
	"github.com/mribichich/release-mgmt/models"
)

type (
	// ApplicationsController represents the controller for operating on the Application resource
	ApplicationsController struct {
		session                *mgo.Session
		applicationsCollection *mgo.Collection
	}
)

// NewApplicationsController provides a reference to a ApplicationsController with provided mongo session
func NewApplicationsController(s *mgo.Session) *ApplicationsController {
	applicationsCollection := s.DB("test").C("applications")
	return &ApplicationsController{s, applicationsCollection}
}

// GetAll retrieves an individual application resource
func (ctrl ApplicationsController) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	result := entities.Applications{}

	// Fetch application
	if err := ctrl.applicationsCollection.Find(bson.M{}).All(&result); err != nil {
		w.WriteHeader(500)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(result)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// Get retrieves an individual application resource
func (ctrl ApplicationsController) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	applicationName := p.ByName("name")
	applicationName = strings.ToLower(applicationName)

	var result *entities.Application

	query := bson.M{"name": bson.RegEx{fmt.Sprintf("^%s$", applicationName), "i"}}

	if err := ctrl.applicationsCollection.Find(query).One(&result); err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(404)
			fmt.Fprintf(w, "%s", "error: application not found")
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "%s", err.Error())
		}
		return
	}

	uj, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// Create creates a new application resource
func (ctrl ApplicationsController) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	applicationPost := models.ApplicationPost{}

	json.NewDecoder(r.Body).Decode(&applicationPost)

	applicationPost.Name = strings.ToLower(applicationPost.Name)

	query := bson.M{"name": bson.RegEx{fmt.Sprintf("^%s$", applicationPost.Name), "i"}}
	n, err := ctrl.applicationsCollection.Find(query).Count()

	if err != nil {
		w.WriteHeader(500)
		return
	}

	if n != 0 {
		w.WriteHeader(400)
		errorModel := "error: application already exists"
		// errorModelJson, _ := json.Marshal(errorModel)
		fmt.Fprintf(w, "%s", errorModel)
		return
	}

	newApplication := entities.Application{
		Id:   bson.NewObjectId(),
		Name: applicationPost.Name}

	err = ctrl.applicationsCollection.Insert(newApplication)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(newApplication)
}

// Update an individual application resource
func (ctrl ApplicationsController) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	applicationName := p.ByName("name")
	applicationName = strings.ToLower(applicationName)

	var application *entities.Application

	query := bson.M{"name": bson.RegEx{fmt.Sprintf("^%s$", applicationName), "i"}}

	if err := ctrl.applicationsCollection.Find(query).One(&application); err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(404)
			fmt.Fprintf(w, "%s", "error: application not found")
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "%s", err.Error())
		}
		return
	}

	applicationPost := models.ApplicationPost{}
	json.NewDecoder(r.Body).Decode(&applicationPost)
	applicationPost.Name = strings.ToLower(applicationPost.Name)

	application.Name = applicationPost.Name

	if err := ctrl.applicationsCollection.UpdateId(application.Id, application); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(application)
}

// Delete removes an existing application resource
func (ctrl ApplicationsController) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	applicationName := p.ByName("name")
	applicationName = strings.ToLower(applicationName)

	query := bson.M{"name": bson.RegEx{fmt.Sprintf("^%s$", applicationName), "i"}}

	if err := ctrl.applicationsCollection.Remove(query); err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(404)
			fmt.Fprintf(w, "%s", "error: application not found")
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "%s", err.Error())
		}
		return
	}

	w.WriteHeader(200)
}
