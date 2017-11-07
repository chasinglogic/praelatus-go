// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package repo provides definitions for abstracting away database interaction
// in Praelatus
package repo

import (
	"errors"

	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/ql/ast"
)

// Errors
var (
	ErrLoginRequired          error = errors.New("you must be logged in")
	ErrAdminRequired                = errors.New("you must be an administrator")
	ErrUnauthorized                 = errors.New("permission denied")
	ErrNotFound                     = errors.New("not found")
	ErrInvalidTicketType            = errors.New("invalid ticket type for project")
	ErrInvalidFieldsForTicket       = errors.New("invalid fields for ticket of that type for project")
)

// TicketRepo handles storing, retrieving, updating, and creating tickets.
type TicketRepo interface {
	Get(u *models.User, uid string) (models.Ticket, error)
	Search(u *models.User, query ast.AST) ([]models.Ticket, error)
	Update(u *models.User, uid string, updated models.Ticket) error
	Create(u *models.User, ticket models.Ticket) (models.Ticket, error)
	Delete(u *models.User, uid string) error

	AddComment(u *models.User, uid string, comment models.Comment) (models.Ticket, error)
	NextTicketKey(u *models.User, projectKey string) (string, error)
	LabelSearch(u *models.User, query string) ([]string, error)
}

// FieldSchemeRepo handles storing, retrieving, updating, and creating field schemes.
type FieldSchemeRepo interface {
	Get(u *models.User, uid string) (models.FieldScheme, error)
	Search(u *models.User, query string) ([]models.FieldScheme, error)
	Update(u *models.User, uid string, updated models.FieldScheme) error
	Create(u *models.User, fieldScheme models.FieldScheme) (models.FieldScheme, error)
	Delete(u *models.User, uid string) error
}

// ProjectRepo handles storing, retrieving, updating, and creating projects.
type ProjectRepo interface {
	Get(u *models.User, uid string) (models.Project, error)
	Search(u *models.User, query string) ([]models.Project, error)
	Update(u *models.User, uid string, updated models.Project) error
	Create(u *models.User, project models.Project) (models.Project, error)
	Delete(u *models.User, uid string) error

	HasLead(u *models.User, lead models.User) ([]models.Project, error)
}

// UserRepo handles storing, retrieving, updating, and creating users.
type UserRepo interface {
	Get(u *models.User, uid string) (models.User, error)
	Search(u *models.User, query string) ([]models.User, error)
	Update(u *models.User, uid string, updated models.User) error
	Create(u *models.User, user models.User) (models.User, error)
	Delete(u *models.User, uid string) error
}

// WorkflowRepo handles storing, retrieving, updating, and creating workflows.
type WorkflowRepo interface {
	Get(u *models.User, uid string) (models.Workflow, error)
	Search(u *models.User, query string) ([]models.Workflow, error)
	Update(u *models.User, uid string, updated models.Workflow) error
	Create(u *models.User, workflow models.Workflow) (models.Workflow, error)
	Delete(u *models.User, uid string) error
}

// NotificationRepo handles storing, retrieving, updating, and creating workflows.
type NotificationRepo interface {
	Create(u *models.User, notification models.Notification) (models.Notification, error)

	MarkRead(u *models.User, uid string) error

	ForProject(u *models.User, project models.Project, onlyUnread bool, last int) ([]models.Notification, error)
	ForUser(u *models.User, user models.User, onlyUnread bool, last int) ([]models.Notification, error)
}

// Repo is a container interface for combining all the other repos.
type Repo interface {
	Tickets() TicketRepo
	Projects() ProjectRepo
	Users() UserRepo
	Fields() FieldSchemeRepo
	Workflows() WorkflowRepo
	Notifications() NotificationRepo

	Clean() error
	Test() error
	Init() error
}

// Cache is used for storing temporary resources. Usually backed by Mongo, Bolt
// or Redis
type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Remove(key string) error
	GetSession(key string) (models.Session, error)
	SetSession(key string, user models.Session) error
	RemoveSession(key string) error
}

// GlobalRepo is used to store a global repo instance that can be accessed from
// anywhere in the application
var GlobalRepo Repo

// Tickets is an alias to the method of the same name on the global Repo
func Tickets() TicketRepo { return GlobalRepo.Tickets() }

// Projects is an alias to the method of the same name on the global Repo
func Projects() ProjectRepo { return GlobalRepo.Projects() }

// Users is an alias to the method of the same name on the global Repo
func Users() UserRepo { return GlobalRepo.Users() }

// Fields is an alias to the method of the same name on the global Repo
func Fields() FieldSchemeRepo { return GlobalRepo.Fields() }

// Workflows is an alias to the method of the same name on the global Repo
func Workflows() WorkflowRepo { return GlobalRepo.Workflows() }

// Notifications is an alias to the method of the same name on the global Repo
func Notifications() NotificationRepo { return GlobalRepo.Notifications() }

// Clean is an alias to the method of the same name on the global Repo
func Clean() error { return GlobalRepo.Clean() }

// Test is an alias to the method of the same name on the global Repo
func Test() error { return GlobalRepo.Test() }

// Init is an alias to the method of the same name on the global Repo
func Init() error { return GlobalRepo.Init() }
