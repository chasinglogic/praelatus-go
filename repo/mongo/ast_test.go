// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package mongo

import (
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/praelatus/praelatus/ql/lexer"
	"github.com/praelatus/praelatus/ql/parser"
)

func TestSimpleEval(t *testing.T) {
	query := "summary = \"test\""
	l := lexer.New(query)
	p := parser.New(l)
	a := p.Parse()

	b := evalAST(a)
	q := bson.M{"summary": "test"}

	if b["summary"] != q["summary"] {
		t.Errorf("Expected: %v Got: %v", bson.M{"summary": "test"}, b)
	}
}

func TestComplexEval(t *testing.T) {
	query := "summary = \"test this parser\" OR (project = \"TEST\" AND (key = \"TEST-1\"))"
	l := lexer.New(query)
	p := parser.New(l)
	a := p.Parse()

	b := evalAST(a)
	t.Log(b)
}

func TestMultipleCompounds(t *testing.T) {
	query := "summary = \"test this parser\" OR project = \"TEST\" AND key = \"TEST-1\""
	l := lexer.New(query)
	p := parser.New(l)
	a := p.Parse()

	b := evalAST(a)
	t.Log(b)
}
