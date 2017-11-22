// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package lexer contains the lexer for PQL
package lexer

import (
	"github.com/praelatus/praelatus/ql/token"
)

// Lexer maintains document position and lexing the input
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

// New create a new lexer for the given input
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

type validator func(byte) bool

func (l *Lexer) read(valid validator) string {
	start := l.position

	for valid(l.ch) {
		l.readChar()
	}

	return l.input[start:l.position]
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(token.EQ, l.ch)
	case '<':
		if l.peekChar() == '=' {
			tok = token.Token{token.LTE, "<="}
			break
		}

		tok = newToken(token.LT, l.ch)
	case '>':
		if l.peekChar() == '=' {
			tok = token.Token{token.GTE, ">="}
			break
		}

		tok = newToken(token.GT, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{token.NE, "!="}
			break
		}

		tok = newToken(token.ILLEGAL, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '~':
		tok = newToken(token.LIKE, l.ch)
	case 0:
		tok = token.Token{token.EOF, ""}
	default:
		if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.read(isDigit)
			return tok
		} else if isLetter(l.ch) {
			tok.Literal = l.read(isLetter)
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if l.ch == '"' {
			// Skip opening quote
			l.readChar()
			tok.Type = token.STRING

			// When inside double quotes numbers and spaces are allows so treat them as such
			tok.Literal = l.read(func(ch byte) bool {
				return isLetter(ch) || ch == ' ' || isDigit(ch)
			})

			// Skip closing quote
			l.readChar()
			return tok
		}

		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func isWhitespace(ch byte) bool {
	return ' ' == ch || '\n' == ch || '\t' == ch
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'z' || ch == '_' || ch == '-' || ch == ','
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tt token.TokenType, ch byte) token.Token {
	return token.Token{Type: tt, Literal: string(ch)}
}
