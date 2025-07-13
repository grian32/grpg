package lex_old

import (
	"fmt"
	"reflect"
)

type TokenType int

var tokenTypeToString = map[TokenType]string{
	LeftParen:  "LeftParen",
	RightParen: "RightParen",
	LeftBrace:  "LeftBrace",
	RightBrace: "RightBrace",
	Comma:      "Comma",
	Dot:        "Dot",
	Minus:      "Minus",
	Plus:       "Plus",
	Semicolon:  "Semicolon",
	Slash:      "Slash",
	Star:       "Star",

	Bang:         "Bang",
	BangEqual:    "BangEqual",
	Equal:        "Equal",
	EqualEqual:   "EqualEqual",
	Greater:      "Greater",
	GreaterEqual: "GreaterEqual",
	Less:         "Less",
	LessEqual:    "LessEqual",

	Identifier: "Identifier",
	String:     "String",
	Int:        "Int",

	And:    "And",
	Else:   "Else",
	False:  "False",
	True:   "True",
	Fnc:    "Fnc",
	If:     "If",
	Nil:    "Nil",
	Or:     "Or",
	Return: "Return",
	Var:    "Var",
	For:    "For",

	Eof: "Eof",
}

// String returns the string name of the TokenType
func (t TokenType) String() string {
	if str, ok := tokenTypeToString[t]; ok {
		return str
	}
	return fmt.Sprintf("TokenType(%d)", t)
}

const (
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	Identifier
	String
	Int

	And
	Else
	False
	True
	Fnc
	If
	Nil
	Or
	Return
	Var
	For

	Eof
)

type Token struct {
	Type    TokenType
	Repr    string // a.k.a lexeme but i think that name is dumb & im not a lang scientist or osmething they can cope and seethe
	Literal any
	Line    uint32
}

// Equal should only be used for testing as it uses reflect.deepequal
func (t *Token) Equal(other Token) bool {
	return t.Type == other.Type && t.Repr == other.Repr && t.Line == other.Line && reflect.DeepEqual(t.Literal, other.Literal)
}

// String should only be used for testing
func (t Token) String() string {
	var literalStr string
	if t.Literal != nil {
		literalStr = fmt.Sprintf("%v (%T)", t.Literal, t.Literal)
	} else {
		literalStr = "nil"
	}

	return fmt.Sprintf("{%v %q %s %d}", t.Type, t.Repr, literalStr, t.Line)
}

func TokenSliceString(slice []Token) string {
	str := "["
	for idx, token := range slice {
		str += token.String()
		if idx != len(slice)-1 {
			str += ", "
		}
	}
	str += "]"

	return str
}
