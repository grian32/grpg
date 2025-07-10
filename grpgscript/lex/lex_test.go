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

func TestParseDoubleSymbols(t *testing.T) {
	file, err1 := os.Open("../testdata/doublesymbols.grpgscript")
	bytes, err2 := io.ReadAll(file)
	if err := cmp.Or(err1, err2); err != nil {
		t.Errorf("Error reading symbols file: %v", err)
	}

	scanner := NewScanner(string(bytes))
	scanner.ScanTokens()

	expectedTokenTypes := []Token{
		{Type: BangEqual, Repr: "!=", Literal: nil, Line: 1},
		{Type: EqualEqual, Repr: "==", Literal: nil, Line: 1},
		{Type: GreaterEqual, Repr: ">=", Literal: nil, Line: 1},
		{Type: LessEqual, Repr: "<=", Literal: nil, Line: 1},
		{Type: Equal, Repr: "=", Literal: nil, Line: 1},
		{Type: Less, Repr: "<", Literal: nil, Line: 1},
		{Type: Greater, Repr: ">", Literal: nil, Line: 1},
		{Type: Bang, Repr: "!", Literal: nil, Line: 1},
		{Type: Eof, Repr: "", Literal: nil, Line: 1},
	}

	output := scanner.Tokens

	if !tokenSliceEquals(expectedTokenTypes, output) {
		t.Errorf("Wanted %v, got %v", expectedTokenTypes, output)
	}
}

func TestParseSymbolsComments(t *testing.T) {
	file, err1 := os.Open("../testdata/symbolscomments.grpgscript")
	bytes, err2 := io.ReadAll(file)
	if err := cmp.Or(err1, err2); err != nil {
		t.Errorf("Error reading symbols file: %v", err)
	}

	scanner := NewScanner(string(bytes))
	scanner.ScanTokens()

	expectedTokenTypes := []Token{
		{Type: BangEqual, Repr: "!=", Literal: nil, Line: 2},
		{Type: Bang, Repr: "!", Literal: nil, Line: 3},
		{Type: Bang, Repr: "!", Literal: nil, Line: 4},
		{Type: Eof, Repr: "", Literal: nil, Line: 4},
	}

	output := scanner.Tokens

	if !tokenSliceEquals(expectedTokenTypes, output) {
		t.Errorf("Wanted %v, got %v", expectedTokenTypes, output)
	}
}

func TestParseInt(t *testing.T) {
	file, err1 := os.Open("../testdata/numbers.grpgscript")
	bytes, err2 := io.ReadAll(file)
	if err := cmp.Or(err1, err2); err != nil {
		t.Errorf("Error reading symbols file: %v", err)
	}

	scanner := NewScanner(string(bytes))
	scanner.ScanTokens()

	expectedTokenTypes := []Token{
		{Type: Int, Repr: "123456", Literal: 123456, Line: 1},
		{Type: Int, Repr: "99999999999", Literal: 99999999999, Line: 2},
		{Type: Int, Repr: "421124142", Literal: 421124142, Line: 3},
		{Type: Eof, Repr: "", Literal: nil, Line: 3},
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
