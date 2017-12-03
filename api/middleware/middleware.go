// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
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

// Middleware is any function which modifies the behavior of a http.Handler
type Middleware func(http.Handler) http.Handler

// Load the middleware for the given handler
func (m Middleware) Load(next http.Handler) http.Handler {
	return m(next)
}

// Chain is a middleware chain
type Chain []Middleware

// Load the middleware in this chain for the given handler
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

// Default is the default middleware stack for Praelatus
var Default = Chain{
	ContentHeaders,
	Logger,
	CORS,
}
