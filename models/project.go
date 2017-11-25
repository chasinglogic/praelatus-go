// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package models

import (
	"time"

	"github.com/praelatus/praelatus/models/permission"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Role is an alias type to make it's use more clear inside of other models.
type Role string

// WorkflowMapping maps a ticket type to a workflow
type WorkflowMapping struct {
	Workflow   bson.ObjectId `json:"workflow"`
	TicketType string        `json:"ticket_type"`
}

// RolePermission maps a role to a permission on a project
type RolePermission struct {
	Role       Role                  `json:"role"`
	Permission permission.Permission `json:"permission"`
}

// Project is the model used to represent a project in the database.
type Project struct {
	Key         string           `json:"key" bson:"_id" required:"true"`
	Name        string           `json:"name" required:"true"`
	CreatedDate time.Time        `json:"createdDate"`
	Lead        string           `json:"lead"`
	Homepage    string           `json:"homepage,omitempty"`
	Repo        string           `json:"repo,omitempty"`
	TicketTypes []string         `json:"ticketTypes"`
	Public      bool             `json:"public"`
	Permissions []RolePermission `json:"permissions"`

	FieldScheme bson.ObjectId `json:"fieldScheme"`

	// Map ticket types to workflow ID's
	WorkflowScheme []WorkflowMapping `json:"workflowScheme"`

	Icon *mgo.GridFile `json:"-"`
}

func (p Project) String() string {
	return jsonString(p)
}

// GetWorkflow will return the ID of the workflow to use for tickets of the
// given type for this project.
func (p Project) GetWorkflow(ticketType string) bson.ObjectId {
	var defaultWorkflow bson.ObjectId

	for _, mapping := range p.WorkflowScheme {
		if mapping.TicketType == ticketType {
			return mapping.Workflow
		} else if mapping.TicketType == "" {
			defaultWorkflow = mapping.Workflow
		}
	}

	return defaultWorkflow
}

// HasTicketType is used to validate whether the given ticket type exists for this project
func (p Project) HasTicketType(typeName string) bool {
	for _, t := range p.TicketTypes {
		if t == typeName {
			return true
		}
	}

	return false
}

// GetPermsForRoles will take the given roles and return a slice of Permissions that those roles have
// NOTE: It does not remove duplicates so some permissions may exist more than once
func (p Project) GetPermsForRoles(roles ...Role) permission.Permissions {
	perms := make(permission.Permissions, 0)

	for _, role := range roles {
		for _, mappings := range p.Permissions {
			if role == mappings.Role {
				perms = perms.Add(mappings.Permission)
			}
		}
	}

	if p.Public {
		perms = perms.Add(permission.ViewProject)
	}

	return perms
}

// HasPermission will return a slice of projects for which the given user has
// the permission indicated out of the projects given.
func HasPermission(permName permission.Permission, user User,
	projects ...Project) []Project {

	if user.IsAdmin {
		return projects
	}

	hasPerm := make([]Project, 0)

	for _, p := range projects {
		roles := user.RolesForProject(p)
		perms := p.GetPermsForRoles(roles...)
		if perms.Contains(permName) {
			hasPerm = append(hasPerm, p)
		}
	}

	return hasPerm
}
