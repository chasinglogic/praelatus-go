package mongo_test

import (
	"testing"
)

func TestTicketGet(t *testing.T) {
	ticket, err := r.Tickets().Get(&admin, "TEST-1")
	if err != nil {
		t.Error(err)
	}

	if ticket.Key == "" {
		t.Error("Expected a key got: ", ticket)
	}
}

func TestTicketSearch(t *testing.T) {
	tks, e := r.Tickets().Search(&admin, "")
	if e != nil {
		t.Error(e)
		return
	}

	if tks == nil || len(tks) == 0 {
		t.Error("Expected to get tickets instead got none.")
	}
}

func TestTicketUpdate(t *testing.T) {
	t.Skip("Ticket Update is unimplemented")

	tk, e := r.Tickets().Get(&admin, "TEST-4")
	if e != nil {
		t.Error(e)
		return
	}

	tk.Summary = "Test ticket save"

	e = r.Tickets().Update(&admin, tk.Key, tk)
	if e != nil {
		t.Error(e)
		return
	}

	tk2, e := r.Tickets().Get(&admin, "TEST-4")
	if e != nil {
		t.Error(e)
		return
	}

	if tk2.Summary != "Test ticket save" {
		t.Errorf("Expected: Test ticket save Got: %s\n", tk.Summary)
	}
}

func TestTicketDelete(t *testing.T) {
	e := r.Tickets().Delete(&admin, "TEST-3")
	if e != nil {
		t.Error(e)
		return
	}

	if _, e = r.Tickets().Get(&admin, "TEST-3"); e == nil {
		t.Errorf("Expected an error getting ticket but got none.")
	}
}
