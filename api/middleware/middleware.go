// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package middleware contains the HTTP middleware used in the api as
// well as utility functions for interacting with them
package middleware

import (
	"net/http"
	"strings"

	"github.com/praelatus/praelatus/api/utils"
)

// Middleware is a function which modifies the behavior of a http.Handler
type Middleware func(http.Handler) http.Handler

// Chain is a slice of Middlewares that can be applied with Load
type Chain []Middleware

// Load will load all middleware in the chain for next returning the wrapped
// handler
func (c Chain) Load(next http.Handler) http.Handler {
	h := next

	for _, m := range c {
		h = m(h)
	}

	return h
}

// ContentHeaders will set the content-type header for the API to application/json
func ContentHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path[:len("/api")] == "/api" {
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") && r.Method != "GET" {
				utils.APIErr(w, http.StatusBadRequest,
					"incorrect content-type")
				return
			}

			w.Header().Set("Content-Type", "application/json")
		}

		next.ServeHTTP(w, r)
	})
}

// CORS allows cross origin requests to the server. Note: By default it allows
// all origins so can be insecure.
// TODO: Make the origins configurable
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Add("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
			w.Header().Add("Access-Control-Expose-Headers", "X-Praelatus-Token, Content-Type")
			w.Header().Add("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")

			if r.Method == "OPTIONS" {
				w.Write([]byte{})
				return
			}

			next.ServeHTTP(w, r)
		})
}

// Default is the default middleware stack for Praelatus
var Default = Chain{
	ContentHeaders,
	Logger,
	CORS,
}
