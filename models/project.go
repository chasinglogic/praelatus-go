package models

import (
	"time"

	"github.com/praelatus/backend/models/permission"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Role is an alias type to make it's use more clear inside of other models.
type Role string

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
	Permissions map[Role][]permission.Permission

	FieldScheme bson.ObjectId `json:"fieldScheme"`

	// Map ticket types to workflow ID's
	WorkflowScheme map[string]bson.ObjectId

	Icon *mgo.GridFile `json:"-"`
}

func (p *Project) String() string {
	return jsonString(p)
}

// GetWorkflow will return the ID of the workflow to use for tickets of the
// given type for this project.
func (p *Project) GetWorkflow(ticketType string) bson.ObjectId {
	for t, worfklowID := range p.WorkflowScheme {
		if t == ticketType {
			return worfklowID
		}
	}

	if defaultWorkflow, ok := p.WorkflowScheme[""]; ok {
		return defaultWorkflow
	}

	return ""
}
