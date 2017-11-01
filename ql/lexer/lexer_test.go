package lexer

import (
	"testing"

	"github.com/praelatus/praelatus/ql/token"
)

// TODO: Write more lexer tests

func TestLexer(t *testing.T) {
	lexerTests := []struct {
		Inp    string
		Tokens []token.Token
	}{
		{
			Inp: "summary = \"test\"",
			Tokens: []token.Token{
				{
					Type:    token.IDENT,
					Literal: "summary",
				},
				{
					Type:    token.EQ,
					Literal: "=",
				},
				{
					Type:    token.STRING,
					Literal: "test",
				},
			},
		},
	}

	for _, test := range lexerTests {
		l := New(test.Inp)

		for _, tok := range test.Tokens {
			curToken := l.NextToken()

			if curToken.Type != tok.Type {
				t.Errorf("Unexpected Token Type Expected: %s Got: %s",
					tok.Type, curToken.Type)
			}

			if curToken.Literal != tok.Literal {
				t.Errorf("Unexpected Token Literal Expected: %s Got: %s",
					tok.Literal, curToken.Literal)
			}
		}
	}
}
