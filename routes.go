package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"ReleaseIndex",
		"GET",
		"/releases",
		ReleaseIndex,
	},
	Route{
		"ReleaseShow",
		"GET",
		"/releases/{id}",
		ReleaseShow,
	},
	Route{
		"ReleaseCreate",
		"POST",
		"/releases",
		ReleaseCreate,
	},
}
