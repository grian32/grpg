package lex

import (
	"cmp"
	"io"
	"os"
	"testing"
)

func TestParseSymbols(t *testing.T) {
	file, err1 := os.Open("../testdata/symbols.grpgscript")
	bytes, err2 := io.ReadAll(file)
	if err := cmp.Or(err1, err2); err != nil {
		t.Errorf("Error reading symbols file: %v", err)
	}

	scanner := NewScanner(string(bytes))
	scanner.ScanTokens()

	expectedTokenTypes := []Token{
		{Type: LeftParen, Repr: "(", Literal: nil, Line: 1},
		{Type: LeftBrace, Repr: "{", Literal: nil, Line: 1},
		{Type: RightBrace, Repr: "}", Literal: nil, Line: 1},
		{Type: RightParen, Repr: ")", Literal: nil, Line: 1},
		{Type: Plus, Repr: "+", Literal: nil, Line: 1},
		{Type: Minus, Repr: "-", Literal: nil, Line: 1},
		{Type: Star, Repr: "*", Literal: nil, Line: 1},
		{Type: Dot, Repr: ".", Literal: nil, Line: 1},
		{Type: Comma, Repr: ",", Literal: nil, Line: 1},
		{Type: Semicolon, Repr: ";", Literal: nil, Line: 1},
		{Type: Eof, Repr: "", Literal: nil, Line: 1},
	}

	output := scanner.Tokens

	if !tokenSliceEquals(expectedTokenTypes, output) {
		t.Errorf("Wanted %v, got %v", expectedTokenTypes, output)
	}
}

func tokenSliceEquals(a, b []Token) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !a[i].Equal(b[i]) {
			return false
		}
	}

	return true
}
