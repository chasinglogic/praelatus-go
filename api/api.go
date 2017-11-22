// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package api has our router and HTTP handlers for all of the available api
// routes
package api

import (
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/praelatus/praelatus/repo"

	"github.com/praelatus/praelatus/api/middleware"
	"github.com/praelatus/praelatus/api/utils"
	"github.com/praelatus/praelatus/api/v1"
)

// Version of praelatus
var Version string

// Commit this version was built against
var Commit string

func routes(router *mux.Router) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}

// Routes will return the mux.Router which contains all of the api routes
func Routes() *mux.Router {
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
			utils.SendJSON(w, struct {
				Version string
				Commit  string
				GOOS    string
				GOARCH  string
			}{
				Version: Version,
				Commit:  Commit,
				GOOS:    runtime.GOOS,
				GOARCH:  runtime.GOARCH,
			})
		})

	v1r.HandleFunc("/routes", routes(v1r)).Methods("GET")
	api.HandleFunc("/routes", routes(api)).Methods("GET")
	return router
}

// New will start running the api on the given port
func New(r repo.Repo, mw middleware.Chain) http.Handler {
	v1.Repo = r
	return mw.Load(Routes())
}
