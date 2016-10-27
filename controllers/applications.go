package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/mribichich/release-mgmt/entities"
)

type (
	// ApplicationsController represents the controller for operating on the Application resource
	ApplicationsController struct {
		session *mgo.Session
	}
)

// NewApplicationsController provides a reference to a ApplicationsController with provided mongo session
func NewApplicationsController(s *mgo.Session) *ApplicationsController {
	return &ApplicationsController{s}
}

// GetApplication retrieves an individual application resource
func (uc ApplicationsController) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub application
	result := entities.Applications{}

	// Fetch application
	if err := uc.session.DB("test").C("applications").Find(bson.M{}).All(&result); err != nil {
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

// GetApplication retrieves an individual application resource
func (uc ApplicationsController) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Stub application
	result := entities.Application{}

	// Fetch application
	if err := uc.session.DB("test").C("applications").FindId(oid).One(&result); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(result)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// CreateApplication creates a new application resource
func (uc ApplicationsController) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an application to be populated from the body
	result := entities.Application{}

	// Populate the application data
	json.NewDecoder(r.Body).Decode(&result)

	// Add an Id
	result.Id = bson.NewObjectId()

	// Write the application to mongo
	err := uc.session.DB("test").C("applications").Insert(result)

	if err != nil {
		log.Fatal(err)
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(result)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

// RemoveApplication removes an existing application resource
func (uc ApplicationsController) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove application
	if err := uc.session.DB("test").C("applications").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}
