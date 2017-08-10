package pg_test

import (
	"testing"

	"github.com/praelatus/backend/models"
)

func TestTypeGet(t *testing.T) {
	tp := models.TicketType{ID: 1}
	e := s.Types().Get(&tp)
	failIfErr("Type Get", t, e)

	if tp.Name == "" {
		t.Errorf("Expected a name got: %s\n", tp.Name)
	}
}

func TestTypeGetAll(t *testing.T) {
	tp, e := s.Types().GetAll()
	failIfErr("Type Get All", t, e)

	if len(tp) == 0 {
		t.Error("Expected to get more than 0 types.")
	}
}

func TestTypeSave(t *testing.T) {
	tp := models.TicketType{ID: 2}
	e := s.Types().Get(&tp)
	failIfErr("Type Save", t, e)

	tp.Name = "Test Save Type"

	e = s.Types().Save(tp)
	failIfErr("Type Save", t, e)

	tp = models.TicketType{ID: 2}
	e = s.Types().Get(&tp)
	failIfErr("Type Save", t, e)

	if tp.Name != "Test Save Type" {
		t.Errorf("Expected: Test Save Type Got: %s\n", tp.Name)
	}
}

func TestTypeRemove(t *testing.T) {
	tp := models.TicketType{ID: 3}
	e := s.Types().Remove(tp)
	failIfErr("Type Remove", t, e)
}
