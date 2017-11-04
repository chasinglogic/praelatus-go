// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package parser contains the parser / validator for PQL
package parser

import (
	"fmt"
	"strconv"

	"github.com/praelatus/praelatus/ql/ast"
	"github.com/praelatus/praelatus/ql/lexer"
	"github.com/praelatus/praelatus/ql/token"
)

const (
	_      int = iota
	LOWEST     // Lowest priority
	ANDOR
	COMPARISON // ==
)

var precedences = map[token.TokenType]int{
	token.EQ:   COMPARISON,
	token.NE:   COMPARISON,
	token.LT:   COMPARISON,
	token.GT:   COMPARISON,
	token.GTE:  COMPARISON,
	token.LTE:  COMPARISON,
	token.LIKE: COMPARISON,
	token.AND:  ANDOR,
	token.OR:   ANDOR,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser parses a PQL query and returns an AST for that query
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New returns a parser for the given lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = map[token.TokenType]prefixParseFn{
		token.IDENT:  p.parseFieldName,
		token.STRING: p.parseString,
		token.INT:    p.parseIntegerLiteral,
		token.LPAREN: p.parseGroupedExpression,
	}

	p.infixParseFns = map[token.TokenType]infixParseFn{
		token.EQ:  p.parseInfixExpression,
		token.NE:  p.parseInfixExpression,
		token.LT:  p.parseInfixExpression,
		token.GT:  p.parseInfixExpression,
		token.GTE: p.parseInfixExpression,
		token.LTE: p.parseInfixExpression,
		token.OR:  p.parseLogicExpression,
		token.AND: p.parseLogicExpression,
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("%s not allowed in comparison expression", t)
	p.errors = append(p.errors, msg)
}

// Parse will turn the given query into an ast.AST
// TODO: Write this
func (p *Parser) Parse() ast.AST {
	var a ast.AST
	a.Modifiers = make([]ast.ModifierStatement, 0)

	for !p.curTokenIs(token.EOF) {
		switch p.curToken.Type {
		case token.LIMIT:
			a.Modifiers = append(a.Modifiers, p.parseModifierStatement())
		case token.ORDER:
			a.Modifiers = append(a.Modifiers, p.parseModifierStatement())
		default:
			a.Query = p.parseExpressionStatement()
		}

		p.nextToken()
	}

	return a
}

func (p *Parser) parseExpressionStatement() ast.ExpressionStatement {
	stmt := ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.EOF) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseModifierStatement() ast.ModifierStatement {
	stmt := ast.ModifierStatement{
		Token:    p.curToken,
		Modifier: p.curToken.Literal,
	}

	if stmt.Token.Type == token.ORDER && !p.peekTokenIs(token.IDENT) {
		p.errors = append(p.errors, "ORDER_BY must be followed by a fieldname")
		return ast.ModifierStatement{}
	}

	if stmt.Token.Type == token.LIMIT && !p.peekTokenIs(token.INT) {
		p.errors = append(p.errors, "LIMIT must be followed by a number")
		return ast.ModifierStatement{}
	}

	p.nextToken()

	if stmt.Token.Type == token.ORDER {
		stmt.Value = p.parseFieldName()
	} else if stmt.Token.Type == token.LIMIT {
		stmt.Value = p.parseIntegerLiteral()
	} else {
		return ast.ModifierStatement{}
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if p.curToken.Type == token.EOF {
		p.errors = append(p.errors, "unexpected end of input")
		return nil
	}

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.EOF) &&
		!p.peekTokenIs(token.ORDER) &&
		!p.peekTokenIs(token.LIMIT) &&
		precedence < p.peekPrecedence() {

		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseString() ast.Expression {
	return ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseFieldName() ast.Expression {
	return ast.FieldLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseLogicExpression(left ast.Expression) ast.Expression {
	expression := ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	if _, ok := expression.Left.(ast.InfixExpression); !ok {
		p.errors = append(p.errors, "logic operators must (AND / OR) must be preceded by a comparison expression")
		return nil
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	if _, ok := expression.Right.(ast.InfixExpression); !ok {
		p.errors = append(p.errors, "logic operators must (AND / OR) must be followed by a comparison expression")
		return nil
	}

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	if _, ok := expression.Right.(ast.FieldLiteral); ok {
		p.errors = append(p.errors, "missing quotes around string: "+expression.Right.String())
		return nil
	}

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}
