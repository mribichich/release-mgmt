package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/mribichich/release-mgmt/entities"
)

type (
	// ReleasesController represents the controller for operating on the Release resource
	ReleasesController struct {
		session *mgo.Session
	}
)

// NewReleasesController provides a reference to a ReleasesController with provided mongo session
func NewReleasesController(s *mgo.Session) *ReleasesController {
	return &ReleasesController{s}
}

// GetRelease retrieves an individual release resource
func (uc ReleasesController) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub release
	result := entities.Releases{}

	// Fetch release
	if err := uc.session.DB("test").C("releases").Find(bson.M{}).All(&result); err != nil {
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

// GetRelease retrieves an individual release resource
func (uc ReleasesController) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Stub release
	u := entities.Release{}

	// Fetch release
	if err := uc.session.DB("test").C("releases").FindId(oid).One(&u); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// CreateRelease creates a new release resource
func (uc ReleasesController) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an release to be populated from the body
	result := entities.Release{}

	// Populate the release data
	json.NewDecoder(r.Body).Decode(&result)

	// Add an Id
	result.Id = bson.NewObjectId()
	result.Timestamp = time.Now()

	// Write the release to mongo
	err := uc.session.DB("test").C("releases").Insert(result)

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

// RemoveRelease removes an existing release resource
func (uc ReleasesController) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove release
	if err := uc.session.DB("test").C("releases").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}
