// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package ast describes an abstract syntax tree for PQL
package ast

import "github.com/praelatus/praelatus/ql/token"

type AST struct {
	Root Expression
}

type Expression interface {
	Node
	expressionNode()
}

// Node is implemented by all nodes in the AST
type Node interface {
	TokenLiteral() string
	String() string
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie InfixExpression) expressionNode()      {}
func (ie InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie InfixExpression) String() string {
	return ie.Left.String() + " " + ie.Operator + " " + ie.Right.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il IntegerLiteral) expressionNode()      {}
func (il IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il IntegerLiteral) String() string       { return il.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl StringLiteral) expressionNode()      {}
func (sl StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl StringLiteral) String() string       { return sl.Token.Literal }

type DateLiteral struct {
	Token token.Token
	Value int64
}

func (dl DateLiteral) expressionNode()      {}
func (dl DateLiteral) TokenLiteral() string { return dl.Token.Literal }
func (dl DateLiteral) String() string       { return dl.Token.Literal }

type FieldLiteral struct {
	Token token.Token
	Value string
}

func (fl FieldLiteral) expressionNode()      {}
func (fl FieldLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl FieldLiteral) String() string       { return fl.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (ident Identifier) expressionNode()      {}
func (ident Identifier) TokenLiteral() string { return ident.Token.Literal }
func (ident Identifier) String() string       { return ident.Token.Literal }
