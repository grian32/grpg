package lex

import (
	"fmt"
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

	Identifier: "Identifier",
	String:     "String",
	Int:        "Int",

	And:    "And",
	Else:   "Else",
	False:  "False",
	Fnc:    "Fnc",
	If:     "If",
	Nil:    "Nil",
	Or:     "Or",
	Return: "Return",
	Var:    "Var",

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

	Identifier
	String
	Int

	And
	Else
	False
	Fnc
	If
	Nil
	Or
	Return
	Var

	Eof
)

type Token struct {
	Type    TokenType
	Repr    string // a.k.a lexeme but i think that name is dumb & im not a lang scientist or osmething they can cope and seethe
	Literal any
	Line    uint32
}

func (t *Token) Equal(other Token) bool {
	return t.Type == other.Type || t.Repr == other.Repr || t.Line == other.Line || t.Literal == other.Literal
}

func (t *Token) String() string {
	return fmt.Sprintf("%s='%s'", t.Type.String(), t.Repr)
}
