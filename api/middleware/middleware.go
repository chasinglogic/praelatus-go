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
	"github.com/praelatus/praelatus/repo"
)

// Cache is the global SessionCache
var Cache repo.Cache

// ContentHeaders will set the content-type header for the API to application/json
func ContentHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path[:len("/api")] == "/api" {
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
				utils.APIErr(w, http.StatusBadRequest,
					"incorrect content-type")
				return
			}
			w.Header().Set("Content-Type", "application/json")
		}

		next.ServeHTTP(w, r)
	})
}

// LoadMw will wrap the given http.Handler in the DefaultMiddleware
func LoadMw(handler http.Handler) http.Handler {
	h := handler

	for _, m := range DefaultMiddleware {
		h = m(h)
	}

	return h
}

// DefaultMiddleware is the default middleware stack for Praelatus
var DefaultMiddleware = []func(http.Handler) http.Handler{
	ContentHeaders,
	Logger,
}
