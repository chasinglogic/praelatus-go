package store

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/models/permission"
)

// DefaultWorkflow should be given when the /api/workflows/default endpoint is
// queried
var DefaultWorkflow = models.Workflow{
	Name: "Simple Workflow",
	Transitions: map[string][]models.Transition{
		"Backlog": {
			{
				Name:     "In Progress",
				ToStatus: models.Status{ID: 2},
				Hooks:    []models.Hook{},
			},
		},
		"In Progress": {
			{
				Name:     "Done",
				ToStatus: models.Status{ID: 3},
				Hooks:    []models.Hook{},
			},
			{
				Name:     "Backlog",
				ToStatus: models.Status{ID: 1},
				Hooks:    []models.Hook{},
			},
		},
		"Done": {
			{
				Name:     "ReOpen",
				ToStatus: models.Status{ID: 1},
				Hooks:    []models.Hook{},
			},
		},
	},
}

// DefaultPermissionScheme is the default permission scheme for a new instance
var DefaultPermissionScheme = models.PermissionScheme{
	Name:        "Default Permission Scheme",
	Description: "The recommended defaults for permissions.",
	Permissions: map[string][]permission.Permission{
		"Administrator": []permission.Permission{
			permission.VIEWPROJECT,
			permission.ADMINPROJECT,
			permission.CREATETICKET,
			permission.COMMENTTICKET,
			permission.REMOVECOMMENT,
			permission.REMOVEOWNCOMMENT,
			permission.EDITOWNCOMMENT,
			permission.EDITCOMMENT,
			permission.TRANSITIONTICKET,
			permission.EDITTICKET,
			permission.REMOVETICKET,
		},
		"Contributor": []permission.Permission{
			permission.VIEWPROJECT,
			permission.CREATETICKET,
			permission.COMMENTTICKET,
			permission.REMOVEOWNCOMMENT,
			permission.EDITOWNCOMMENT,
			permission.TRANSITIONTICKET,
			permission.EDITTICKET,
		},
		"User": []permission.Permission{
			permission.VIEWPROJECT,
			permission.CREATETICKET,
			permission.COMMENTTICKET,
			permission.REMOVEOWNCOMMENT,
			permission.EDITOWNCOMMENT,
		},
		"Anonymous": []permission.Permission{
			permission.VIEWPROJECT,
		},
	},
}

var defaults = []func(s Store) error{
	SeedTicketTypes,
	SeedFields,
	SeedStatuses,
	SeedWorkflows,
	SeedPermissions,
}

var seedFuncs = []func(s Store) error{
	SeedUsers,
	SeedTeams,
	SeedProjects,
	SeedTicketTypes,
	SeedFields,
	SeedStatuses,
	SeedLabels,
	SeedWorkflows,
	SeedTickets,
	SeedComments,
	SeedPermissions,
}

// SeedDefaults will seed the database with the basics needed to use Praelatus
func SeedDefaults(s Store) error {
	log.Println("Seeding database with defaults...")
	for _, f := range defaults {
		e := f(s)
		if e != nil {
			return e
		}
	}

	return nil
}

var dev bool

// SeedAll will run all of the seed functions
func SeedAll(s Store) error {
	fmt.Println("Seeding All")
	dev = true
	for _, f := range seedFuncs {
		if s == nil {
			fmt.Println("WTF")
		}

		e := f(s)
		if e != nil {
			return e
		}
	}

	return nil
}

// SeedLabels will add some test labesl to the database
func SeedLabels(s Store) error {
	labels := []models.Label{
		{
			Name: "test",
		},
		{
			Name: "duplicate",
		},
		{
			Name: "wontfix",
		},
	}

	for _, l := range labels {
		e := s.Labels().New(&l)
		if e != nil && e != ErrDuplicateEntry {
			return e
		}
	}

	return nil
}

// SeedTickets will add some test tickets to the database
func SeedTickets(s Store) error {
	fmt.Println("Seeding tickets")
	priorities := &models.FieldOption{
		Selected: []string{"HIGH", "MEDIUM", "LOW"}[rand.Intn(3)],
		Options:  []string{"HIGH", "MEDIUM", "LOW"},
	}

	for i := 0; i < 50; i++ {
		t := models.Ticket{
			Key:         s.Tickets().NextTicketKey(models.Project{ID: 1, Key: "TEST"}),
			Summary:     "This is a test ticket. #" + strconv.Itoa(i),
			Description: "No really, this is just a test",
			WorkflowID:  1,
			Reporter:    models.User{ID: 1},
			Assignee:    models.User{ID: 1},
			Status:      models.Status{ID: 1},
			Labels: []models.Label{
				{
					ID:   1,
					Name: "test",
				},
			},
			Fields: []models.FieldValue{
				{
					Name:  "Story Points",
					Value: rand.Intn(100),
				},
				{
					Name:  "Priority",
					Value: priorities,
				},
			},
			Type: models.TicketType{ID: 1},
		}

		e := s.Tickets().New(models.Project{ID: 1, Key: "TEST"}, &t)
		if e != nil && e != ErrDuplicateEntry {
			return e
		}
	}

	return nil
}

// SeedStatuses will add some ticket statuses to the database
func SeedStatuses(s Store) error {
	statuses := []models.Status{
		{
			Name: "Backlog",
		},
		{
			Name: "In Progress",
		},
		{
			Name: "Done",
		},
		{
			Name: "For Saving",
		},
		{
			Name: "For Removing",
		},
	}

	fmt.Println("Seeding statuses")
	for _, st := range statuses {
		e := s.Statuses().New(&st)
		if e != nil && e != ErrDuplicateEntry {
			return e
		}
	}

	return nil
}

// SeedComments will add some comments to all tickets
func SeedComments(s Store) error {
	fmt.Println("Seeding comments")

	t, se := s.Tickets().GetAll(models.User{ID: 1})
	if se != nil {
		return se
	}

	for _, tk := range t {
		noOfComments := rand.Intn(30) + 1

		for x := 1; x < noOfComments; x++ {
			c := &models.Comment{
				Body: fmt.Sprintf(`This is the %d th comment
# Yo Dawg
**I** *heard* you
> like markdown
so I put markdown in your comments`,
					x),
				Author: models.User{ID: 1},
			}

			e := s.Tickets().NewComment(tk, c)
			if e != nil && e != ErrDuplicateEntry {
				return e
			}
		}

	}

	return nil
}

// SeedFields will seed the given store with some test Fields.
func SeedFields(s Store) error {
	priorities := &models.FieldOption{
		Selected: "LOW",
		Options:  []string{"HIGH", "MEDIUM", "LOW"},
	}

	fields := []models.Field{
		{
			Name:     "Story Points",
			DataType: "INT",
		},
		{
			Name:     "TestField2",
			DataType: "FLOAT",
		},
		{
			Name:     "TestField3",
			DataType: "INT",
		},
		{
			Name:     "TestField4",
			DataType: "DATE",
		},
		{
			Name:     "Priority",
			DataType: "OPT",
			Options:  priorities,
		},
	}

	fmt.Println("Seeding fields")
	for _, f := range fields {
		e := s.Fields().New(&f)
		if e != nil && e != ErrDuplicateEntry {
			return e
		}

		if e == ErrDuplicateEntry {
			return nil
		}

		e = s.Fields().AddToProject(models.User{ID: 1}, models.Project{ID: 1}, &f)
		if e != nil && e != ErrDuplicateEntry {
			return e
		}

		if e == ErrDuplicateEntry {
			return nil
		}
	}

	return nil
}

// SeedProjects will seed the given store with some test projects.
func SeedProjects(s Store) error {
	projects := []models.Project{
		{
			Name: "TEST Project",
			Key:  "TEST",
			Lead: models.User{ID: 1},
		},
		{
			Name: "TEST Project 2",
			Key:  "TEST2",
			Lead: models.User{ID: 2},
		},
		{
			Name: "TEST Project 3",
			Key:  "TEST3",
			Lead: models.User{ID: 2},
		},
	}

	fmt.Println("Seeding projects")
	for _, p := range projects {
		e := s.Projects().New(&p)
		if e != nil && e != ErrDuplicateEntry {
			return e
		}

		if e == ErrDuplicateEntry {
			return nil
		}
	}

	return nil
}

// SeedTeams will seed the database with some test Teams.
func SeedTeams(s Store) error {
	teams := []models.Team{
		{
			Name: "The A Team",
			Lead: models.User{
				ID: 1,
			},
			Members: []models.User{
				{ID: 1},
				{ID: 2},
			},
		},
		{
			Name: "The B Team",
			Lead: models.User{
				ID: 2,
			},
			Members: []models.User{
				{ID: 1},
				{ID: 2},
			},
		},
	}

	fmt.Println("Seeding teams")
	for _, team := range teams {
		team.Lead = models.User{ID: 1}

		e := s.Teams().New(&team)
		if e != nil && e != ErrDuplicateEntry {
			return e
		}

		if e == ErrDuplicateEntry {
			return nil
		}
	}

	return nil
}

// SeedTicketTypes will seed the database with some test TicketTypes.
func SeedTicketTypes(s Store) error {
	types := []models.TicketType{
		{
			Name: "Bug",
		},
		{
			Name: "Epic",
		},
		{
			Name: "Story",
		},
		{
			Name: "Feature",
		},
		{
			Name: "Question",
		},
	}

	fmt.Println("Seeding ticket types")
	for _, t := range types {
		e := s.Types().New(&t)
		if e != nil && e != ErrDuplicateEntry {
			return e
		}

		if e == ErrDuplicateEntry {
			return nil
		}
	}

	return nil
}

// SeedUsers will seed the database with some test users.
func SeedUsers(s Store) error {
	t1, be := models.NewUser("testuser", "test", "Test Testerson",
		"test@example.com", false)
	if be != nil {
		return be
	}

	t2, be := models.NewUser("testadmin", "test", "Test Testerson II",
		"test1@example.com", true)
	if be != nil {
		return be
	}

	t3, be := models.NewUser("testadmin2", "test", "Test Testerson III",
		"test1@example.com", true)
	if be != nil {
		return be
	}

	t4, be := models.NewUser("saveusertest", "test", "Test Testerson IV",
		"test4@example.com", true)
	if be != nil {
		return be
	}

	users := []models.User{
		*t2,
		*t1,
		*t3,
		*t4,
	}

	fmt.Println("Seeding users")
	for _, u := range users {
		e := s.Users().New(&u)
		if e != nil && e != ErrDuplicateEntry {
			return e
		}

		if e == ErrDuplicateEntry {
			return nil
		}
	}

	return nil
}

// SeedWorkflows will seed the database with some workflows
func SeedWorkflows(s Store) error {
	p1 := models.Project{ID: 1}
	p2 := models.Project{ID: 2}
	wk1 := DefaultWorkflow

	fmt.Println("Seeding workflows")
	e := s.Workflows().New(p1, &wk1)
	if e != nil && e != ErrDuplicateEntry {
		return e
	}

	if !dev {
		return nil
	}

	wk1.Name = wk1.Name + "-TEST"
	e = s.Workflows().New(p1, &wk1)
	if e != nil && e != ErrDuplicateEntry {
		return e
	}

	wk1.Name = wk1.Name + "-TEST1"
	e = s.Workflows().New(p2, &wk1)
	if e != nil && e != ErrDuplicateEntry {
		return e
	}

	wk1.Name = wk1.Name + "-TEST2"
	e = s.Workflows().New(p2, &wk1)
	if e != nil && e != ErrDuplicateEntry {
		return e
	}

	return nil
}

// SeedPermissions will put the default permission scheme in the
// database
func SeedPermissions(s Store) error {
	fmt.Println("Seeding permissions")

	s1 := DefaultPermissionScheme
	s2 := s1
	s2.Name = "Copy of Default Permission Scheme"
	s3 := s1
	s3.Name = "Save Me"
	s4 := s1
	s4.Name = "Delete Me"

	schemes := []models.PermissionScheme{s1, s2, s3, s4}

	for _, p := range schemes {
		err := s.Permissions().New(&p)
		if err != nil && err != ErrDuplicateEntry {
			return err
		}
	}

	return nil
}
