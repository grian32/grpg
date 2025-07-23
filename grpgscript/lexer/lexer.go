package lexer

import (
	"grpgscript/token"
	"grpgscript/util"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

var singleCharTokens = map[byte]token.TokenType{
	';': token.SEMICOLON,
	'(': token.LPAREN,
	')': token.RPAREN,
	'{': token.LBRACE,
	'}': token.RBRACE,
	',': token.COMMA,
	'*': token.ASTERISK,
}

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
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	singleCharToken, exists := singleCharTokens[l.ch]
	if exists {
		tok = newToken(singleCharToken, l.ch)
		// need to advance to the next char over, this behaviour is replicated below before returning.
		// not necessary for ints/literals as those advance before returning on readint/literal
		l.readChar()
		return tok
	}

	switch l.ch {
	case '=':
		tok = l.ifNextIsDoubleLen('=', token.EQ, token.ASSIGN)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		tok = l.ifNextIsDoubleLen('=', token.NOT_EQ, token.BANG)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if util.IsAlpha(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if util.IsDigit(l.ch) {
			tok.Literal = l.readInt()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// TODO: pass in func to make read* generic?
func (l *Lexer) readIdentifier() string {
	startPos := l.position

	for util.IsAlpha(l.ch) {
		l.readChar()
	}

	return l.input[startPos:l.position]
}

func (l *Lexer) readInt() string {
	startPos := l.position

	for util.IsDigit(l.ch) {
		l.readChar()
	}

	return l.input[startPos:l.position]
}

func (l *Lexer) readString() string {
	startPos := l.position + 1 // +1 cuz "
	for {
		l.readChar()

		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[startPos:l.position]
}

func (l *Lexer) ifNextIsDoubleLen(char byte, t, f token.TokenType) token.Token {
	if l.peekChar() == char {
		ch := l.ch
		l.readChar()
		return token.Token{Type: t, Literal: string(ch) + string(l.ch)}
	} else {
		return newToken(f, l.ch)
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
