package lexer

import (
	"grpgscript/token"
	"testing"
)

// BenchmarkNextToken
// This isn't meant to be used as an actual benchmark due to how it appends, real usage has parser
// calling lexer.NextToken(), It's mainly meant to be used just to track any possible performance regressions.
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
	for b.Loop() {
		l := New(input)

		var tokenArr []token.Token

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			tokenArr = append(tokenArr, tok)
		}
	}
}
