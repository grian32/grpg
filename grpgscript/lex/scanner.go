package lex

import (
	"fmt"
	"grpgscript/util"
	"log"
)

type Scanner struct {
	Source  string
	Tokens  []Token
	start   uint32
	current uint32
	line    uint32
}

var oneLengthTokenMap = map[rune]TokenType{
	'(': LeftParen,
	')': RightParen,
	'{': LeftBrace,
	'}': RightBrace,
	',': Comma,
	'.': Dot,
	'-': Minus,
	'+': Plus,
	';': Semicolon,
	'*': Star,
}

func NewScanner(src string) *Scanner {
	return &Scanner{
		Source: src,
		Tokens: make([]Token, 0),
	}
}

func (s *Scanner) ScanTokens() {
	srcLen := uint32(len(s.Source))
	s.line = 1

	for s.current < srcLen {
		s.start = s.current
		s.ScanToken()
	}

	s.Tokens = append(s.Tokens, Token{
		Type:    Eof,
		Repr:    "",
		Literal: nil,
		Line:    s.line,
	})
}

func (s *Scanner) ScanToken() {
	char := s.Advance()

	// quick lookup for single char tokens
	token, exists := oneLengthTokenMap[char]
	if exists {
		s.AddToken(token, nil)
		return
	}

	switch char {
	case '!':
		s.AddToken(s.IfNextIsT('=', BangEqual, Bang), nil)
	case '=':
		s.AddToken(s.IfNextIsT('=', EqualEqual, Equal), nil)
	case '>':
		s.AddToken(s.IfNextIsT('=', GreaterEqual, Greater), nil)
	case '<':
		s.AddToken(s.IfNextIsT('=', LessEqual, Less), nil)
	case '"':
		s.String()
	case '/':
		if s.NextIs('/') {
			for s.Peek() != '\n' && !s.IsAtEnd() {
				s.Advance()
			}
		} else {
			s.AddToken(Slash, nil)
		}
	case ' ', '\r', '\t':
	case '\n':
		s.line++
	default:
		if util.IsDigit(char) {
			s.Int()
		} else {
			log.Printf("Unrecognized char %c, %d", char, s.line)
		}
	}
}

func (s *Scanner) String() {
	for s.Peek() != '"' && !s.IsAtEnd() {
		if s.Peek() == '\n' {
			s.line++
		}
		s.Advance()
	}

	if s.IsAtEnd() {
		fmt.Printf("Unterminated string %d", s.line)
	}

	s.Advance()

	// TODO: escapes? maybe?
	s.AddToken(String, s.Source[s.start+1:s.current-1])
}

func (s *Scanner) Int() {
	for util.IsDigit(s.Peek()) {
		s.Advance()
	}

	s.AddToken(Int, s.Source[s.start:s.current])
}

func (s *Scanner) AddToken(token TokenType, literal any) {
	text := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, Token{
		Type:    token,
		Repr:    text,
		Literal: literal,
		Line:    s.line,
	})
}

func (s *Scanner) Advance() rune {
	char := s.Source[s.current]
	s.current += 1
	return rune(char)
}

func (s *Scanner) IfNextIsT(next rune, t TokenType, f TokenType) TokenType {
	if s.IsAtEnd() {
		return f
	}
	if rune(s.Source[s.current]) != next {
		return f
	}

	s.current += 1
	return t
}

func (s *Scanner) NextIs(next rune) bool {
	if s.IsAtEnd() {
		return false
	}
	if rune(s.Source[s.current]) != next {
		return false
	}

	s.current += 1

	return true
}

func (s *Scanner) Peek() rune {
	if s.IsAtEnd() {
		return '\000' // craft: = \0 java
	}
	return rune(s.Source[s.current])
}

func (s *Scanner) IsAtEnd() bool {
	return s.current >= uint32(len(s.Source))
}
