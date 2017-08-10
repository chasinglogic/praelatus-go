package pg_test

import (
	"testing"

	"github.com/praelatus/backend/models"
)

func TestStatusGet(t *testing.T) {
	l := &models.Status{ID: 1}
	e := s.Statuses().Get(l)
	failIfErr("Status Get", t, e)

	if l == nil {
		t.Error("Expected a status and got nil instead.")
	}

	if l.Name == "" {
		t.Errorf("Expected status to have name got %s\n", l.Name)
	}
}

func TestStatusGetAll(t *testing.T) {
	l, e := s.Statuses().GetAll()
	failIfErr("Status Get All", t, e)

	if l == nil {
		t.Error("Expected to get some stores and got nil instead.")
	}

	if len(l) == 0 {
		t.Errorf("Expected stores to be returned got %d stores instead\n", len(l))
	}
}

func TestStatusSave(t *testing.T) {
	l := &models.Status{ID: 4}
	e := s.Statuses().Get(l)
	failIfErr("Status Save", t, e)

	l.Name = "SAVE_TEST_STATUS"

	e = s.Statuses().Save(*l)
	failIfErr("Status Save", t, e)

	l = &models.Status{ID: 4}
	e = s.Statuses().Get(l)
	failIfErr("Status Save", t, e)

	if l.Name != "SAVE_TEST_STATUS" {
		t.Errorf("Expected: SAVE_TEST_STATUS Got: %s\n", l.Name)
	}
}

func TestStatusRemove(t *testing.T) {
	l := models.Status{ID: 5}
	e := s.Statuses().Remove(l)
	failIfErr("Status Remove", t, e)
}
