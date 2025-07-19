package parser

import (
	"grpgscript/lexer"
	"testing"
)

func BenchmarkParseProgram(b *testing.B) {
	input := `var five = 5;
var ten = 10;

var add = fnc(x, y) {
    x + y;
};

var result = add(five, ten);

if (5 < 10) {
	add(ten, five);
} else {
    add(five, five);
}
`

	for b.Loop() {
		l := lexer.New(input)
		p := New(l)

		result := p.ParseProgram()
		_ = result
	}
}
