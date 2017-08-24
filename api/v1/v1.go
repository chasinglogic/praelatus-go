// Package v1 contains all of the routes and handlers for v1 of the
// API
package v1

import (
	"github.com/gorilla/mux"
	"github.com/praelatus/backend/config"
	mgo "gopkg.in/mgo.v2"
)

// Conn is the global database connection used in our HTTP handlers.
var Conn *mgo.Session

// convenience function for getting a collection by name
func getCollection(collName string) *mgo.Collection {
	return Conn.Copy().DB(config.DBName()).C(collName)
}

// Routes will set up the appropriate routes on the given mux.Router
func Routes(router *mux.Router) {
	fieldRouter(router)
	projectRouter(router)
	ticketRouter(router)
	userRouter(router)
	workflowRouter(router)
	miscRouter(router)
}
