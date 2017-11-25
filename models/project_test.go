package models

import (
	"testing"

	"github.com/praelatus/praelatus/models/permission"
	"gopkg.in/mgo.v2/bson"
)

func permsAreEqual(perms1 permission.Permissions, perms2 permission.Permissions) bool {
	for _, p := range perms1 {
		if !perms2.Contains(permission.Permission(p)) {
			return false
		}
	}

	return true
}

func TestGetPermsForRoles(t *testing.T) {
	p := Project{
		Permissions: []RolePermission{
			{
				Role:       Role("Administrator"),
				Permission: permission.AdminProject,
			},
			{
				Role:       Role("Administrator"),
				Permission: permission.CreateTicket,
			},
		},
	}

	adminPerms := p.GetPermsForRoles(Role("Administrator"))
	anonPerms := p.GetPermsForRoles(Role(""))

	if len(anonPerms) != 0 {
		t.Errorf("Expected No Perms Got %v", anonPerms)
	}

	if adminPerms == nil {
		t.Errorf("Expected nil Got %v", anonPerms)
	}

	expectedPerms := permission.Permissions{}.Add(permission.CreateTicket,
		permission.AdminProject)

	if !permsAreEqual(adminPerms, expectedPerms) {
		t.Errorf("Expected %v Got %v", expectedPerms, adminPerms)
	}
}

func TestGetWorkflow(t *testing.T) {
	workflowID := bson.NewObjectId()
	defaultID := bson.NewObjectId()

	p := Project{
		WorkflowScheme: []WorkflowMapping{
			{
				Workflow:   workflowID,
				TicketType: "BUG",
			},
			{
				Workflow: defaultID,
			},
		},
	}

	if workflowID != p.GetWorkflow("BUG") {
		t.Errorf("Expected %s Got %s", workflowID, p.GetWorkflow("BUG"))
	}

	if defaultID != p.GetWorkflow("") {
		t.Errorf("Expected %s Got %s", workflowID, p.GetWorkflow(""))
	}

	if defaultID != p.GetWorkflow("NONE") {
		t.Errorf("Expected %s Got %s", workflowID, p.GetWorkflow("NONE"))
	}
}

func TestHasTicketType(t *testing.T) {
	p := Project{
		TicketTypes: []string{
			"TEST",
		},
	}

	if !p.HasTicketType("TEST") {
		t.Error("Expected True for TEST but got False")
	}

	if p.HasTicketType("NONE") {
		t.Error("Expected False for NONE but got True")
	}
}

func TestHashPermission(t *testing.T) {
	p := Project{
		Key:    "UNITTEST",
		Public: true,
		Permissions: []RolePermission{
			{
				Role:       Role("Administrator"),
				Permission: permission.AdminProject,
			},
			{
				Role:       Role("Administrator"),
				Permission: permission.CreateTicket,
			},
		},
	}

	if len(HasPermission(permission.AdminProject, User{IsAdmin: true}, p)) == 0 {
		t.Error("Expected SysAdmin to have all permissions. Got none instead.")
	}

	has := HasPermission(permission.AdminProject, User{}, p)
	if len(has) != 0 {
		t.Error("Expected anonymous user to not have admin permission.")
	}

	has = HasPermission(permission.AdminProject, User{
		Roles: []UserRole{
			{
				Project: "UNITTEST",
				Role:    Role("Administrator"),
			},
		},
	}, p)
	if len(has) == 0 {
		t.Error("Expected project admin to have admin permission.")
	}

	has = HasPermission(permission.ViewProject, User{}, p)
	if len(has) == 0 {
		t.Error("Expected anon user to have view permission to public project.")
	}
}
