package lex

import "log"

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
	char := s.Source[s.current]
	s.current += 1

	token, exists := oneLengthTokenMap[rune(char)]
	if !exists {
		log.Printf("Invalid Character: %c", char)
	} else {
		s.AddToken(token, nil)
	}
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
