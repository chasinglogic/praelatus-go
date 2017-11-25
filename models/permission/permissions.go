// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package permission is used to store the various permission in
// praelatus, it acts as a pseudo-enumeration
package permission

// Permission is an alias type to make it's use more clear inside of models.
type Permission string

// Permissions is used to add some convenience functions for checking permisisons
type Permissions []string

func (p Permissions) Contains(permName Permission) bool {
	for _, perm := range p {
		if Permission(perm) == permName {
			return true
		}
	}

	return false
}

func (p Permissions) Add(perms ...Permission) Permissions {
	for _, perm := range perms {
		p = append(p, string(perm))
	}

	return p
}

// These are the permissions available in Praelatus
const (
	ViewProject      Permission = "VIEW_PROJECT"
	AdminProject                = "ADMIN_PROJECT"
	CreateTicket                = "CREATE_TICKET"
	CommentTicket               = "COMMENT_TICKET"
	RemoveComment               = "REMOVE_COMMENT"
	RemoveOwnComment            = "REMOVE_OWN_COMMENT"
	EditOwnComment              = "EDIT_OWN_COMMENT"
	EditComment                 = "EDIT_COMMENT"
	TransitionTicket            = "TRANSITION_TICKET"
	EditTicket                  = "EDIT_TICKET"
	RemoveTicket                = "REMOVE_TICKET"
)

// ListOfPermissions holds available permissions in a slice. This is valuable for
// various areas where we need to return all permissions or iterate
// permissions.
var ListOfPermissions = [...]Permission{
	ViewProject,
	AdminProject,
	CreateTicket,
	CommentTicket,
	RemoveComment,
	RemoveOwnComment,
	EditOwnComment,
	EditComment,
	TransitionTicket,
	EditTicket,
	RemoveTicket,
}

// ValidPermission will verify that a given permission string is valid.
func ValidPermission(permName Permission) bool {
	for _, p := range ListOfPermissions {
		if p == permName {
			return true
		}
	}

	return false
}
