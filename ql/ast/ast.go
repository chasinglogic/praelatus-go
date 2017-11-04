// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package ast describes an abstract syntax tree for PQL
package ast

import (
	"strings"
	"time"

	"github.com/praelatus/praelatus/ql/token"
)

// AST is the parsed abstract syntax tree of a query
type AST struct {
	Query     ExpressionStatement
	Modifiers []ModifierStatement
}

func (a AST) String() string {
	q := a.Query.String()

	for _, mod := range a.Modifiers {
		q = q + " " + mod.String()
	}

	return q
}

// TODO: Implement methods to make traversing the tree simpler.

// Expression represents an AST node that evaluates to a value
type Expression interface {
	Node
	expressionNode()
}

// Statement represents an AST node that performs some action
type Statement interface {
	Node
	statementNode()
}

// Literal is any literal value
type Literal interface {
	GetValue() interface{}
}

// Node is implemented by all nodes in the AST
type Node interface {
	TokenLiteral() string
	String() string
}

// InfixExpression is an expression with a left and right side, as well as an
// operator to perform on those sides
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie InfixExpression) expressionNode() {}

// TokenLiteral implements AST node
func (ie InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

// IntegerLiteral is a node whose value is an Integer
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il IntegerLiteral) expressionNode() {}

// TokenLiteral implements AST node
func (il IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il IntegerLiteral) String() string       { return il.Token.Literal }

func (il IntegerLiteral) GetValue() interface{} { return il.Value }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl StringLiteral) expressionNode() {}

// TokenLiteral implements AST node
func (sl StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl StringLiteral) String() string       { return "\"" + sl.Token.Literal + "\"" }

// Value impelements literal
func (sl StringLiteral) GetValue() interface{} { return sl.Value }

// DateLiteral is a string representing a date
type DateLiteral struct {
	Token token.Token
	Value time.Time
}

func (dl DateLiteral) expressionNode() {}

// TokenLiteral implements AST node
func (dl DateLiteral) TokenLiteral() string { return dl.Token.Literal }
func (dl DateLiteral) String() string       { return dl.Token.Literal }

// Value impelements literal
func (dl DateLiteral) GetValue() interface{} { return dl.Value }

// FieldLiteral is a field name either a custom field or otherwise
type FieldLiteral struct {
	Token token.Token
	Value string
}

func (fl FieldLiteral) expressionNode() {}

// TokenLiteral implements AST node
func (fl FieldLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl FieldLiteral) String() string       { return fl.Token.Literal }

// GetValue implements literal
func (fl FieldLiteral) GetValue() string { return fl.Value }

// IsCustomField will return whether or not the given field is a custom field
func (fl FieldLiteral) IsCustomField() bool {
	defaultFields := []string{
		"createddate",
		"updateddate",
		"key",
		"summary",
		"description",
		"status",
		"reporter",
		"assignee",
		"type",
		"labels",
		"project",
	}

	val := strings.ToLower(fl.Value)

	for _, defaultField := range defaultFields {
		if val == defaultField {
			return false
		}
	}

	return true
}

// ExpressionStatement is the root node of a query
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral implements AST node
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

// ModifierStatement is a statement which modifies the behavior of a query, i.e. LIMIT or ORDER BY
type ModifierStatement struct {
	Token    token.Token
	Modifier string
	Value    Expression
}

func (es *ModifierStatement) statementNode() {}

// TokenLiteral implements AST node
func (es *ModifierStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ModifierStatement) String() string {
	if es.Modifier != "" {
		return es.Modifier + " " + es.Value.String()
	}

	return ""
}
