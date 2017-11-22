// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package mongo_test

import (
	"strings"
	"testing"

	"github.com/praelatus/praelatus/ql/ast"
	"github.com/praelatus/praelatus/ql/lexer"
	"github.com/praelatus/praelatus/ql/parser"
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
	tks, e := r.Tickets().Search(&admin, ast.AST{})
	if e != nil {
		t.Error(e)
		return
	}

	if tks == nil || len(tks) == 0 {
		t.Error("Expected to get tickets instead got none.")
	}
}

func TestTicketSearchLimit(t *testing.T) {
	l := lexer.New("LIMIT 5")
	p := parser.New(l)
	a := p.Parse()

	if p.Errors() != nil {
		t.Errorf("Parsing error: %v", p.Errors())
		return
	}

	tks, e := r.Tickets().Search(&admin, a)
	if e != nil {
		t.Error(e)
		return
	}

	if tks == nil || len(tks) == 0 {
		t.Error("Expected to get tickets instead got none.")
	}

	if len(tks) != 5 {
		t.Errorf("Expected 5 tickets Got %d", len(tks))
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

func TestLabelSearch(t *testing.T) {
	lbls, e := r.Tickets().LabelSearch(&admin, "example")
	if e != nil {
		t.Error(e)
		return
	}

	if len(lbls) == 0 {
		t.Error("Expected 2 Labels Got 0")
		return
	}

	for i := range lbls {
		if !strings.HasPrefix(lbls[i], "example") {
			t.Errorf("Expected labels to start with example- Got %s\n", lbls[i])
			return
		}
	}
}
