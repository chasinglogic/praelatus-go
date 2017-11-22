// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package mongo implements repo.Repo for a mongodb database.
package mongo

import (
	"log"

	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/models/permission"
	"github.com/praelatus/praelatus/repo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var dbLog = log.New(config.LogWriter(), "[MONGO] ", log.LstdFlags)
var cacheLog = log.New(config.LogWriter(), "[MONGO_CACHE] ", log.LstdFlags)

// MongoDB Collection and Database names
const (
	dbName = "praelatus"

	projects      = "projects"
	fieldSchemes  = "field_schemes"
	tickets       = "tickets"
	users         = "users"
	sessions      = "sessions"
	cache         = "cache"
	workflows     = "workflows"
	notifications = "notifications"
)

func mongoErr(e error) error {
	if e == nil {
		return e
	}

	dbLog.Println("[MONGO] ERROR:", e)

	// TODO: Catch other repo errors
	switch e.Error() {
	case "not found":
		return repo.ErrNotFound
	default:
		return e
	}
}

func permQuery(u *models.User) bson.M {
	if u == nil {
		u = &models.User{}
	}

	if u.IsAdmin {
		return bson.M{}
	}

	viewPerms := make([]bson.M, len(u.Roles)+1)
	for i, r := range u.Roles {
		viewPerms[i] = bson.M{
			"_id": r.Project,
			"permissions": bson.M{
				"role":       r.Role,
				"permission": permission.ViewProject,
			},
		}
	}

	viewPerms[len(u.Roles)] = bson.M{
		"public": true,
	}

	query := bson.M{
		"$or": viewPerms,
	}

	return query
}

// Repo contains all model specific repos.
type Repo struct {
	Conn          *mgo.Session
	tickets       ticketRepo
	users         userRepo
	projects      projectRepo
	fieldSchemes  fieldSchemeRepo
	workflows     workflowRepo
	notifications notificationRepo
}

// Fields returns the fieldSchemesRepo implementation for mongodb
func (r Repo) Fields() repo.FieldSchemeRepo {
	return r.fieldSchemes
}

// Tickets returns the ticketRepo implementation for mongodb
func (r Repo) Tickets() repo.TicketRepo {
	return r.tickets
}

// Projects returns the projectRepo implementation for mongodb
func (r Repo) Projects() repo.ProjectRepo {
	return r.projects
}

// Workflows returns the workflowRepo implementation for mongodb
func (r Repo) Workflows() repo.WorkflowRepo {
	return r.workflows
}

// Notifications returns the notificationRepo implementation for mongodb
func (r Repo) Notifications() repo.NotificationRepo {
	return r.notifications
}

// Users returns the userRepo implementation for mongodb
func (r Repo) Users() repo.UserRepo {
	return r.users
}

// Clean will remove all data from the database
func (r Repo) Clean() error {
	return r.Conn.DB(dbName).Run("dropDatabase", nil)
}

// Test will test connection to the database
func (r Repo) Test() error {
	return r.Conn.Ping()
}

// Init will setup the indexes on the database
func (r Repo) Init() error {
	return nil
}

// New will attempt to connect to the MongoDB instance at connURL and return
// the repo.Repo that is connected to it.
func New(connURL string) repo.Repo {
	conn, err := mgo.Dial(connURL)
	if err != nil {
		panic(err)
	}

	return Repo{
		Conn:          conn,
		tickets:       ticketRepo{conn},
		projects:      projectRepo{conn},
		workflows:     workflowRepo{conn},
		fieldSchemes:  fieldSchemeRepo{conn},
		users:         userRepo{conn},
		notifications: notificationRepo{conn},
	}
}
