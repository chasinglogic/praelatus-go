// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package api has our router and HTTP handlers for all of the available api
// routes
package api

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/praelatus/praelatus/repo"

	"github.com/praelatus/praelatus/api/middleware"
	"github.com/praelatus/praelatus/api/utils"
	"github.com/praelatus/praelatus/api/v1"
)

// API holds global API state
type API struct {
	router    *mux.Router
	Endpoints []Endpoint
	Version   string
	Commit    string
}

func (a API) routes(w http.ResponseWriter, r *http.Request) {
	rs := []string{}

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		rs = append(rs, t)
		return nil
	})

	utils.SendJSON(w, rs)
}

// Routes will return the mux.Router which contains all of the api routes
func (a API) Routes() http.Handler {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	v1r := api.PathPrefix("/v1").Subrouter()

	// setup v1 routes
	v1.Routes(v1r)

	// setup latest routes
	v1.Routes(api)

	// setup routes endpoints
	api.HandleFunc("/version",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(fmt.Sprintf("Praelatus %s#%s %s/%s",
				Version, Commit, runtime.GOOS, runtime.GOARCH)))
		})

	v1r.HandleFunc("/routes", routes(v1r)).Methods("GET")
	api.HandleFunc("/routes", routes(api)).Methods("GET")

	static := http.StripPrefix("/assets/",
		http.FileServer(http.Dir("client/assets/")))

	router.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path[len("api"):] == "api" {
				api.ServeHTTP(w, r)
			} else if r.URL.Path[len("static"):] == "static" {
				static.ServeHTTP(w, r)
			} else {
				w.Header().Set("Content-Type", "text/html")
				http.ServeFile(w, r, "client/index.html")
			}
		})

	return router
}

// New will start running the api on the given port
func New(r repo.Repo, cache repo.Cache) http.Handler {
	v1.Repo = r
	middleware.Cache = cache

	router := Routes()
	return middleware.LoadMw(router)
}
