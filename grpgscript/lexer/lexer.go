package lexer

import (
	"grpgscript/token"
	"grpgscript/util"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	currLine     int
	// this is literally just == position/readPosition but resets on newline
	currCol int
	readCol int
	ch      byte
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
	l.currCol = l.readCol
	l.readPosition += 1
	l.readCol += 1
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
		tok = newToken(singleCharToken, l.ch, l.currLine, l.currCol, l.readCol)
		// need to advance to the next char over, this behaviour is replicated below before returning.
		// not necessary for ints/literals as those advance before returning on readint/literal
		l.readChar()
		return tok
	}

	switch l.ch {
	case '=':
		tok = l.ifNextIsDoubleLen('=', token.EQ, token.ASSIGN)
	case '+':
		tok = newToken(token.PLUS, l.ch, l.currLine, l.currCol, l.readCol)
	case '-':
		tok = newToken(token.MINUS, l.ch, l.currLine, l.currCol, l.readCol)
	case '!':
		tok = l.ifNextIsDoubleLen('=', token.NOT_EQ, token.BANG)
	case '/':
		tok = newToken(token.SLASH, l.ch, l.currLine, l.currCol, l.readCol)
	case '<':
		tok = newToken(token.LT, l.ch, l.currLine, l.currCol, l.readCol)
	case '>':
		tok = newToken(token.GT, l.ch, l.currLine, l.currCol, l.readCol)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, l.currLine, l.currCol, l.readCol)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, l.currLine, l.currCol, l.readCol)
	case ':':
		tok = newToken(token.COLON, l.ch, l.currLine, l.currCol, l.readCol)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
		tok.Line = uint64(l.currLine)
		tok.Col = uint64(l.currCol - len(tok.Literal) - 1) // -1 for "
		tok.End = uint64(l.readCol)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = uint64(l.currLine)
		tok.Col = uint64(0)
		tok.End = uint64(0)
	default:
		if util.IsAlpha(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = uint64(l.currLine)
			tok.Col = uint64(l.currCol - len(tok.Literal))
			tok.End = uint64(l.currCol)
			return tok
		} else if util.IsDigit(l.ch) {
			tok.Literal = l.readInt()
			tok.Type = token.INT
			tok.Line = uint64(l.currLine)
			tok.Col = uint64(l.currCol - len(tok.Literal))
			tok.End = uint64(l.currCol)
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.currLine, l.currCol, l.readCol)
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
		return token.Token{Type: t, Literal: string(ch) + string(l.ch), Line: uint64(l.currLine), Col: uint64(l.currCol - 1), End: uint64(l.readCol)}
	} else {
		return newToken(f, l.ch, l.currLine, l.currCol, l.readCol)
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.currLine++
			l.readCol = 0
		}
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte, line, col, end int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: uint64(line), Col: uint64(col), End: uint64(end)}
}
