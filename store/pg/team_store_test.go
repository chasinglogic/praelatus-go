package pg_test

import (
	"testing"

	"github.com/praelatus/backend/models"
)

func TestTeamGet(t *testing.T) {
	team := &models.Team{ID: 1}
	e := s.Teams().Get(team)
	failIfErr("Team Get", t, e)

	if len(team.Members) == 0 || team.Name == "" {
		t.Error("Team should have members", team)
	}
}

func TestTeamGetAll(t *testing.T) {
	teams, e := s.Teams().GetAll()
	failIfErr("Team Get All", t, e)

	if len(teams) == 0 || teams == nil {
		t.Error("No teams found.")
	}
}

func TestTeamGetForUser(t *testing.T) {
	u := models.User{ID: 1}
	teams, e := s.Teams().GetForUser(u)
	failIfErr("Team Get For User", t, e)

	if len(teams) == 0 || teams == nil {
		t.Error("No teams found.")
	}
}

func TestTeamSave(t *testing.T) {
	team := &models.Team{ID: 1}
	e := s.Teams().Get(team)
	failIfErr("Team Save", t, e)

	team.Name = "Test Team Save"

	e = s.Teams().Save(*team)
	failIfErr("Team Save", t, e)

	team = &models.Team{ID: 1}
	e = s.Teams().Get(team)
	failIfErr("Team Save", t, e)

	if team.Name != "Test Team Save" {
		t.Errorf("Expected: Test Team Save Got: %s\n", team.Name)
	}
}

func TestTeamRemove(t *testing.T) {
	e := s.Teams().Remove(models.Team{ID: 2})
	failIfErr("Team Remove", t, e)
}
