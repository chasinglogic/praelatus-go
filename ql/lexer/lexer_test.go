package lexer

import "testing"

func TestLexer(t *testing.T) {
	l := New("summary = \"test\"")
	t.Log(l.NextToken())
}
