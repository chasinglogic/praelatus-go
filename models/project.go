package models

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Role is an alias type to make it's use more clear inside of other models.
type Role string

// Permission is an alias type to make it's use more clear inside of other
// models. for available permissions view
// github.com/praelatus/backend/permissions
type Permission string

// Project is the model used to represent a project in the database.
type Project struct {
	Key         string    `json:"key" bson:"_id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"createdDate"`
	Lead        string    `json:"lead"`
	Homepage    string    `json:"homepage,omitempty"`
	Repo        string    `json:"repo,omitempty"`
	TicketTypes []string  `json:"ticketTypes"`
	// Map roles to permissions
	Permissions map[Role][]Permission

	FieldScheme bson.ObjectId `json:"fieldScheme"`

	Icon *mgo.GridFile `json:"-"`
}

func (p *Project) String() string {
	return jsonString(p)
}
