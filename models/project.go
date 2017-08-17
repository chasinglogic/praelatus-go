package models

import (
	"fmt"
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
	CreatedDate time.Time `json:"created-date"`
	Lead        string    `json:"lead"`
	Homepage    string    `json:"homepage,omitempty"`
	Repo        string    `json:"repo,omitempty"`
	TicketTypes []string  `json:"ticket-types"`
	// Map roles to permissions
	Permissions map[Role][]permission.Permission

	FieldScheme bson.ObjectId `json:"field-scheme"`

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

func (p *Project) HasTicketType(typeName string) bool {
	for _, t := range p.TicketTypes {
		if t == typeName {
			return true
		}
	}

	return false
}

// HasPermission will return a slice of projects for which the given user has
// the permission indicated out of the projects given.
func HasPermission(permName permission.Permission, user User, projects ...Project) []Project {

	// Skip perm checking for SysAdmins
	if user.IsAdmin {
		return projects
	}

	hasPermission := make([]Project, len(projects))
	i := 0

projects:
	for _, p := range projects {
		role := user.Permissions[p.Key]

		permissions := p.Permissions[role]
		anon := p.Permissions["Anonymous"]

		for _, perm := range permissions {
			if perm == permName {
				hasPermission[i] = p
				i++
				continue projects
			}
		}

		fmt.Println("Checking anon permissions")
		for _, perm := range anon {
			if perm == permName {
				hasPermission[i] = p
				i++
				continue projects
			}
		}

	}

	return hasPermission[:i]
}
