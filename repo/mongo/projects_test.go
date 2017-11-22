// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package mongo_test

import (
	"testing"

	"github.com/praelatus/praelatus/models"
)

func TestProjectGet(t *testing.T) {
	project, err := r.Projects().Get(&admin, "TEST")
	if err != nil {
		t.Error(err)
		return
	}

	if project.Key == "" {
		t.Error("Expected a key got: ", project)
	}
}

func TestProjectSearch(t *testing.T) {
	ps, e := r.Projects().Search(&admin, "")
	if e != nil {
		t.Error(e)
		return
	}

	if ps == nil || len(ps) == 0 {
		t.Error("Expected to get projects instead got none.")
	}
}

func TestProjectUpdate(t *testing.T) {
	p, e := r.Projects().Get(&admin, "TEST")
	if e != nil {
		t.Error(e)
		return
	}

	p.Name = "Test project save"

	e = r.Projects().Update(&admin, p.Key, p)
	if e != nil {
		t.Error(e)
		return
	}

	p2, e := r.Projects().Get(&admin, "TEST")
	if e != nil {
		t.Error(e)
		return
	}

	if p2.Name != "Test project save" {
		t.Errorf("Expected: Test project save Got: %s\n", p.Name)
	}
}

func TestProjectDelete(t *testing.T) {
	e := r.Projects().Delete(&admin, "TEST2")
	if e != nil {
		t.Error(e)
		return
	}

	if _, e = r.Projects().Get(&admin, "TEST2"); e == nil {
		t.Errorf("Expected an error getting project but got none.")
	}
}

func TestProjectHasLead(t *testing.T) {
	p, e := r.Projects().HasLead(&admin, models.User{Username: "testadmin"})
	if e != nil {
		t.Error(e)
		return
	}

	if p == nil || len(p) == 0 {
		t.Error("Expected to get projects instead got none.")
	}

	if p[0].Lead != "testadmin" {
		t.Error("Expected to get a project with testadmin as a lead")
	}

}
