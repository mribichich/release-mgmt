package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/mribichich/release-mgmt/entities"
	"github.com/mribichich/release-mgmt/models"
)

type (
	// ReleasesController represents the controller for operating on the Release resource
	ReleasesController struct {
		session *mgo.Session
		// releasesCollection     *mgo.Collection
		applicationsCollection *mgo.Collection
	}
)

// NewReleasesController provides a reference to a ReleasesController with provided mongo session
func NewReleasesController(s *mgo.Session) *ReleasesController {
	// releasesCollection := s.DB("test").C("releases")
	applicationsCollection := s.DB("test").C("applications")
	return &ReleasesController{s, applicationsCollection}
}

// GetAll retrieves an individual release resource
func (ctrl ReleasesController) GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	applicationName := p.ByName("name")
	applicationName = strings.ToLower(applicationName)

	var application entities.Application

	query := bson.M{"name": bson.RegEx{fmt.Sprintf("^%s$", applicationName), "i"}}
	err := ctrl.applicationsCollection.Find(query).One(&application)

	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(404)
			fmt.Fprintf(w, "%s", "error: application not found")
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "%s", err.Error())
		}
		return
	}

	// uj, _ := json.Marshal(application.Releases)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	// fmt.Fprintf(w, "%s", uj)
	json.NewEncoder(w).Encode(application.Releases)
}

// Get retrieves an individual release resource
func (ctrl ReleasesController) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	applicationName := p.ByName("name")
	releaseVersion := p.ByName("version")

	applicationName = strings.ToLower(applicationName)

	var application entities.Application

	query := bson.M{"name": bson.RegEx{fmt.Sprintf("^%s$", applicationName), "i"}}
	err := ctrl.applicationsCollection.Find(query).One(&application)

	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(404)
			fmt.Fprintf(w, "%s", "error: application not found")
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "%s", err.Error())
		}
		return
	}

	var release entities.Release

	for _, rel := range application.Releases {
		if rel.Version == releaseVersion {
			release = rel
			break
		}
	}

	if release == (entities.Release{}) {
		w.WriteHeader(404)
		fmt.Fprintf(w, "%s", "error: release not found")
		return
	}

	uj, _ := json.Marshal(release)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// Create creates a new release resource
func (ctrl ReleasesController) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	applicationName := p.ByName("name")
	applicationName = strings.ToLower(applicationName)

	var application entities.Application

	query := bson.M{"name": bson.RegEx{fmt.Sprintf("^%s$", applicationName), "i"}}
	err := ctrl.applicationsCollection.Find(query).One(&application)

	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(404)
			fmt.Fprintf(w, "%s", "error: application not found")
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "%s", err.Error())
		}
		return
	}

	releasePost := models.ReleasePost{}

	json.NewDecoder(r.Body).Decode(&releasePost)

	newRelease := entities.Release{}
	newRelease.Id = bson.NewObjectId()
	newRelease.Timestamp = time.Now()
	newRelease.Version = releasePost.Version

	application.Releases = append(application.Releases, newRelease)

	err = ctrl.applicationsCollection.UpdateId(application.Id, application)

	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(201)
}

// Update a release resource
func (ctrl ReleasesController) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	applicationName := p.ByName("name")
	releaseVersion := p.ByName("version")

	applicationName = strings.ToLower(applicationName)

	var application entities.Application

	query := bson.M{"name": bson.RegEx{fmt.Sprintf("^%s$", applicationName), "i"}}
	err := ctrl.applicationsCollection.Find(query).One(&application)

	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(404)
			fmt.Fprintf(w, "%s", "error: application not found")
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "%s", err.Error())
		}
		return
	}

	var release entities.Release
	var releaseIndex int

	for i, rel := range application.Releases {
		if rel.Version == releaseVersion {
			release = rel
			releaseIndex = i
			break
		}
	}

	if release == (entities.Release{}) {
		w.WriteHeader(404)
		fmt.Fprintf(w, "%s", "error: release not found")
		return
	}

	releasePost := models.ReleasePost{}

	json.NewDecoder(r.Body).Decode(&releasePost)

	release.Version = releasePost.Version

	application.Releases[releaseIndex] = release

	err = ctrl.applicationsCollection.UpdateId(application.Id, application)

	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(404)
			fmt.Fprintf(w, "%s", "error: application not found")
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "%s", err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(release)
}

// Delete removes an existing release resource
func (ctrl ReleasesController) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	applicationName := p.ByName("name")
	releaseVersion := p.ByName("version")

	applicationName = strings.ToLower(applicationName)

	var application entities.Application

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

	var found bool

	for _, rel := range application.Releases {
		if rel.Version == releaseVersion {
			found = true
			break
		}
	}

	if !found {
		w.WriteHeader(404)
		fmt.Fprintf(w, "%s", "error: release not found")
		return
	}

	query = bson.M{"$pull": bson.M{"releases": bson.M{"version": releaseVersion}}}

	if err := ctrl.applicationsCollection.UpdateId(application.Id, query); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	w.WriteHeader(200)
}
