package evaluator

import (
	"grpgscript/lexer"
	"grpgscript/object"
	"grpgscript/parser"
	"testing"
)

func BenchmarkEval(b *testing.B) {
	input := `
var fib = fnc(n) {
		if (n < 2) {
			n
		} else {
			fib(n - 1) + fib(n - 2)
		}
};

var makeAdder = fnc(x) {
	fnc(y) { x + y; };
};

var addTwo = makeAdder(2);
var addThree = makeAdder(3);

var x = addTwo(1);

fib(10) + x;

var arr = [1,2,3,4];

fib(arr[3]) + 3;

var map = {"hey": 5, "bye": 6}l

var getMap = fnc(bool) {
	if (bool) {
		return map["hey"]
	} else {
	    return map["bye"]
	}
}

fib(getMap(true)) + fib(getMap(false))
	`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()

	for b.Loop() {
		env := object.NewEnvironment()
		result := Eval(program, env)
		_ = result
	}
}

// BenchmarkEvalBuiltinCalls is practically meant to simulate calling a bunch of builtins, as you might do in actual game scripts, will extend further with actual game builtins once they exist.
func BenchmarkEvalBuiltinCalls(b *testing.B) {
	input := `
		var x = len("Hello, World!");
		var y = len([1,2,3,4,5]);
		var z = len("Bye, World!");
		var a = len("Welcome Back, World!");
		var b = len(["here", "we", "are"]);
		var c = len([true, false, true, true, true, false]);
	`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()

	for b.Loop() {
		env := object.NewEnvironment()
		result := Eval(program, env)
		_ = result
	}
}
