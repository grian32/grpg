package lexer

import (
	"grpgscript/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5;
var ten = 10;

var add = fnc(x, y) {
    x + y;
};

var result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
    return true;
} else {
    return false;
}

10 == 10;
10 != 9;

"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    uint64
		expectedCol     uint64
		expectedEnd     uint64
	}{
		{token.VAR, "var", 0, 0, 3},
		{token.IDENT, "five", 0, 4, 8},
		{token.ASSIGN, "=", 0, 9, 10},
		{token.INT, "5", 0, 11, 12},
		{token.SEMICOLON, ";", 0, 12, 13},
		{token.VAR, "var", 1, 0, 3},
		{token.IDENT, "ten", 1, 4, 7},
		{token.ASSIGN, "=", 1, 8, 9},
		{token.INT, "10", 1, 10, 12},
		{token.SEMICOLON, ";", 1, 12, 13},
		{token.VAR, "var", 3, 0, 3},
		{token.IDENT, "add", 3, 4, 7},
		{token.ASSIGN, "=", 3, 8, 9},
		{token.FUNCTION, "fnc", 3, 10, 13},
		{token.LPAREN, "(", 3, 13, 14},
		{token.IDENT, "x", 3, 14, 15},
		{token.COMMA, ",", 3, 15, 16},
		{token.IDENT, "y", 3, 17, 18},
		{token.RPAREN, ")", 3, 18, 19},
		{token.LBRACE, "{", 3, 20, 21},
		{token.IDENT, "x", 4, 4, 5},
		{token.PLUS, "+", 4, 6, 7},
		{token.IDENT, "y", 4, 8, 9},
		{token.SEMICOLON, ";", 4, 9, 10},
		{token.RBRACE, "}", 5, 0, 1},
		{token.SEMICOLON, ";", 5, 1, 2},
		{token.VAR, "var", 7, 0, 3},
		{token.IDENT, "result", 7, 4, 10},
		{token.ASSIGN, "=", 7, 11, 12},
		{token.IDENT, "add", 7, 13, 16},
		{token.LPAREN, "(", 7, 16, 17},
		{token.IDENT, "five", 7, 17, 21},
		{token.COMMA, ",", 7, 21, 22},
		{token.IDENT, "ten", 7, 23, 26},
		{token.RPAREN, ")", 7, 26, 27},
		{token.SEMICOLON, ";", 7, 27, 28},
		{token.BANG, "!", 8, 0, 1},
		{token.MINUS, "-", 8, 1, 2},
		{token.SLASH, "/", 8, 2, 3},
		{token.ASTERISK, "*", 8, 3, 4},
		{token.INT, "5", 8, 4, 5},
		{token.SEMICOLON, ";", 8, 5, 6},
		{token.INT, "5", 9, 0, 1},
		{token.LT, "<", 9, 2, 3},
		{token.INT, "10", 9, 4, 6},
		{token.GT, ">", 9, 7, 8},
		{token.INT, "5", 9, 9, 10},
		{token.SEMICOLON, ";", 9, 10, 11},
		{token.IF, "if", 11, 0, 2},
		{token.LPAREN, "(", 11, 3, 4},
		{token.INT, "5", 11, 4, 5},
		{token.LT, "<", 11, 6, 7},
		{token.INT, "10", 11, 8, 10},
		{token.RPAREN, ")", 11, 10, 11},
		{token.LBRACE, "{", 11, 12, 13},
		{token.RETURN, "return", 12, 4, 10},
		{token.TRUE, "true", 12, 11, 15},
		{token.SEMICOLON, ";", 12, 15, 16},
		{token.RBRACE, "}", 13, 0, 1},
		{token.ELSE, "else", 13, 2, 6},
		{token.LBRACE, "{", 13, 7, 8},
		{token.RETURN, "return", 14, 4, 10},
		{token.FALSE, "false", 14, 11, 16},
		{token.SEMICOLON, ";", 14, 16, 17},
		{token.RBRACE, "}", 15, 0, 1},
		{token.INT, "10", 17, 0, 2},
		{token.EQ, "==", 17, 3, 5},
		{token.INT, "10", 17, 6, 8},
		{token.SEMICOLON, ";", 17, 8, 9},
		{token.INT, "10", 18, 0, 2},
		{token.NOT_EQ, "!=", 18, 3, 5},
		{token.INT, "9", 18, 6, 7},
		{token.SEMICOLON, ";", 18, 7, 8},
		{token.STRING, "foobar", 20, 0, 8},
		{token.STRING, "foo bar", 21, 0, 9},
		{token.LBRACKET, "[", 22, 0, 1},
		{token.INT, "1", 22, 1, 2},
		{token.COMMA, ",", 22, 2, 3},
		{token.INT, "2", 22, 4, 5},
		{token.RBRACKET, "]", 22, 5, 6},
		{token.SEMICOLON, ";", 22, 6, 7},
		{token.LBRACE, "{", 23, 0, 1},
		{token.STRING, "foo", 23, 1, 6},
		{token.COLON, ":", 23, 6, 7},
		{token.STRING, "bar", 23, 8, 13},
		{token.RBRACE, "}", 23, 13, 14},
		{token.EOF, "", 24, 0, 0},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong, expected=%d, got=%d; @%v", i, tt.expectedLine, tok.Line, tok)
		}

		if tok.Col != tt.expectedCol {
			t.Fatalf("tests[%d] - col wrong, expected=%d, got=%d", i, tt.expectedCol, tok.Col)
		}

		if tok.End != tt.expectedEnd {
			t.Fatalf("tests[%d] - end wrong, expected=%d, got=%d", i, tt.expectedEnd, tok.End)
		}
	}
}
