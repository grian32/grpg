package lexer

import (
	"grpgscript/token"
	"testing"
)

func BenchmarkNextToken(b *testing.B) {
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
`
	l := New(input)

	var tokenArr []token.Token

	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		tokenArr = append(tokenArr, tok)
	}
}
