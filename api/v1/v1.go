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
