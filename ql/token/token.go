// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package token contains the tokens for PQL
package token

type TokenType string

// TokenTypes
const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF               = "EOF"

	// Identifiers and literals
	IDENT  = "IDENT"
	COMMA  = ","
	INT    = "INT"
	STRING = "STRING"

	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="

	EQ = "="
	NE = "!="

	LIKE = "~"

	LPAREN = "("
	RPAREN = ")"

	AND = "AND"
	OR  = "OR"

	ORDER = "ORDER_BY"
	LIMIT = "LIMIT"
)

var keywords = map[string]TokenType{
	"AND":      AND,
	"and":      AND,
	"OR":       OR,
	"or":       OR,
	"ORDER_BY": ORDER,
	"order_by": ORDER,
	"LIMIT":    LIMIT,
	"limit":    LIMIT,
}

// Token contains the string literal and the type for a given token
type Token struct {
	Type    TokenType
	Literal string
}

// LookupIdent will determine if the given word is a keyword or a identifier
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
