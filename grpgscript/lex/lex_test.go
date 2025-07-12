package lex

import (
	"cmp"
	"io"
	"os"
	"testing"
)

// testdata
var (
	symbolsData = []Token{
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
	doubleSymbolsData = []Token{
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
	symbolsCommentsData = []Token{
		{Type: BangEqual, Repr: "!=", Literal: nil, Line: 2},
		{Type: Bang, Repr: "!", Literal: nil, Line: 3},
		{Type: Bang, Repr: "!", Literal: nil, Line: 4},
		{Type: Eof, Repr: "", Literal: nil, Line: 4},
	}
	intsData = []Token{
		{Type: Int, Repr: "123456", Literal: 123456, Line: 1},
		{Type: Int, Repr: "99999999999", Literal: 99999999999, Line: 2},
		{Type: Int, Repr: "421124142", Literal: 421124142, Line: 3},
		{Type: Eof, Repr: "", Literal: nil, Line: 3},
	}
	stringsData = []Token{
		{Type: String, Repr: "\"hello this is a string\"", Literal: "hello this is a string", Line: 1},
		{Type: String, Repr: "\"hello\nthis\nis\na\nmultiline\nstring\"", Literal: "hello\nthis\nis\na\nmultiline\nstring", Line: 7},
		{Type: Eof, Repr: "", Literal: nil, Line: 8},
	}
	helloWorldData = []Token{
		{Type: Identifier, Repr: "printf", Literal: nil, Line: 1},
		{Type: LeftParen, Repr: "(", Literal: nil, Line: 1},
		{Type: String, Repr: "\"Hello, world!\\n\"", Literal: "Hello, world!\\n", Line: 1},
		{Type: RightParen, Repr: ")", Literal: nil, Line: 1},
		{Type: Eof, Repr: "", Literal: nil, Line: 1},
	}
	functionsData = []Token{
		{Type: Fnc, Repr: "fnc", Literal: nil, Line: 1},
		{Type: Identifier, Repr: "helloWorld", Literal: nil, Line: 1},
		{Type: LeftParen, Repr: "(", Literal: nil, Line: 1},
		{Type: RightParen, Repr: ")", Literal: nil, Line: 1},
		{Type: LeftBrace, Repr: "{", Literal: nil, Line: 1},
		{Type: Return, Repr: "return", Literal: nil, Line: 2},
		{Type: String, Repr: "\"Hello, world!\\n\"", Literal: "Hello, world!\\n", Line: 2},
		{Type: RightBrace, Repr: "}", Literal: nil, Line: 3},
		{Type: Identifier, Repr: "printf", Literal: nil, Line: 5},
		{Type: LeftParen, Repr: "(", Literal: nil, Line: 5},
		{Type: String, Repr: "\"%s\\n\"", Literal: "%s\\n", Line: 5},
		{Type: Comma, Repr: ",", Literal: nil, Line: 5},
		{Type: Identifier, Repr: "helloWorld", Line: 5},
		{Type: LeftParen, Repr: "(", Literal: nil, Line: 5},
		{Type: RightParen, Repr: ")", Literal: nil, Line: 5},
		{Type: RightParen, Repr: ")", Literal: nil, Line: 5},
		{Type: Eof, Repr: "", Literal: nil, Line: 5},
	}
)

func TestLex(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		expected  []Token
	}{
		{"ParseSymbols", "../testdata/symbols.grpgscript", symbolsData},
		{"ParseDoubleSymbols", "../testdata/doublesymbols.grpgscript", doubleSymbolsData},
		{"ParseSymbolsComments", "../testdata/symbolscomments.grpgscript", symbolsCommentsData},
		{"ParseInt", "../testdata/numbers.grpgscript", intsData},
		{"ParseStrings", "../testdata/strings.grpgscript", stringsData},
		{"ParseHelloWorld", "../testdata/helloworld.grpgscript", helloWorldData},
		{"ParseFunction", "../testdata/functions.grpgscript", functionsData},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, err1 := os.Open(test.inputFile)
			bytes, err2 := io.ReadAll(file)

			if err := cmp.Or(err1, err2); err != nil {
				t.Errorf("Error reading file %s, on test %s: %v", test.inputFile, test.name, err)
			}

			scanner := NewScanner(string(bytes))
			scanner.ScanTokens()

			output := scanner.Tokens

			if !tokenSliceEquals(test.expected, output) {
				t.Errorf("Wanted %v, got %v on test %s", test.expected, output, test.name)
			}
		})
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
