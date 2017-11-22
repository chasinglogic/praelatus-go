// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package v1 contains all of the routes and handlers for v1 of the
// API
package v1

import (
	"github.com/gorilla/mux"
	"github.com/praelatus/praelatus/repo"
)

// Repo is the global database connection
var Repo repo.Repo

// Routes will set up the appropriate routes on the given mux.Router
func Routes(router *mux.Router) {
	fieldRouter(router)
	projectRouter(router)
	ticketRouter(router)
	userRouter(router)
	workflowRouter(router)
	miscRouter(router)
}
