package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mribichich/release-mgmt/entities"
	"github.com/mribichich/release-mgmt/repositories"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func ReleaseIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	releases := repositories.FindAll()
	if err := json.NewEncoder(w).Encode(releases); err != nil {
		panic(err)
	}
}

func ReleaseShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	releaseId := vars["id"]
	fmt.Fprintln(w, "Release show: ", releaseId)
}

func ReleaseCreate(w http.ResponseWriter, r *http.Request) {
	var release entities.Release
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &release); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := repositories.RepoCreateRelease(release)
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}
