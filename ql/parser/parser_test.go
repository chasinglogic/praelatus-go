// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package parser

import (
	"testing"

	"github.com/praelatus/praelatus/ql/ast"
	"github.com/praelatus/praelatus/ql/lexer"
)

func TestParse(t *testing.T) {
	l := lexer.New("summary = \"test\"")
	p := New(l)
	tree := p.Parse()

	inf, ok := tree.Root.(ast.InfixExpression)
	if !ok {
		t.Errorf("Expected an ast.InfixExpression Got %T", tree.Root)
		return
	}

	_, ok = inf.Left.(ast.FieldLiteral)
	if !ok {
		t.Errorf("Expected an ast.FieldLiteral Got %T", inf.Left)
		return
	}

	_, ok = inf.Right.(ast.StringLiteral)
	if !ok {
		t.Errorf("Expected an ast.StringLiteral Got %T", inf.Right)
		return
	}
}

func TestParseOR(t *testing.T) {
	l := lexer.New("summary = \"test this parser\" OR project = \"TEST\"")
	p := New(l)
	tree := p.Parse()
	t.Log(tree)
}

func TestParseAND(t *testing.T) {
	l := lexer.New("summary = \"test this parser\" AND project = \"TEST\"")
	p := New(l)
	tree := p.Parse()
	t.Log(tree)
}

func TestParseComplex(t *testing.T) {
	l := lexer.New("summary = \"test this parser\" OR (project = \"TEST\" AND (key = \"TEST-1\"))")
	p := New(l)
	tree := p.Parse()
	t.Log(tree)
}
