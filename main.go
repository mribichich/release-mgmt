package main

import (
	"log"
	"net/http"
)

func main() {
	url := ":8080"

	log.Printf("listening at '%s' ...", url)

	router := NewRouter()

	log.Fatal(http.ListenAndServe(url, router))
}
