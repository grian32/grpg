package evaluator

import (
	"grpgscript/lexer"
	"grpgscript/object"
	"grpgscript/parser"
	"grpgscript/perf"
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

var map = {"hey": 5, "bye": 6}

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

	perf.ConstFold(program)

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

	perf.ConstFold(program)

	b.ResetTimer()

	for b.Loop() {
		env := object.NewEnvironment()
		result := Eval(program, env)
		_ = result
	}
}

const foldTestInput = `
var x = (1 + 2) * (3 + 4);
var y = x + 10;
var z = y * 2;

var a = 5 * 5 + 5 * (2 + 3);
var b = a + 100;
var c = -x;

var notTrue = !true;
var notFalse = !false;
var doubleNot = !!true;
var notComparison = !(1 < 2);
var notEquality = !(1 == 2);
var notExpression = !(5 + 5 == 10);

var arrLen = len([1, 2, 3, 4]);
var strLen = len("hello");
var mapLen = len({"a": 1, "b": 2, "c": 3});

var foldedArr = [1 + 2, 3 * 4, !(false)];
var foldedMap = {"a": 5 * 5, "b": 1 + 1};

b + z - c;
`
const foldTestInputLarge = `
var a = (1 + 2) * (3 + 4);
var b = -(5 * 6);
var c = !false;
var d = !!true;
var e = !(1 < 2);
var f = !!((2 + 2) == 4);
var g = !((10 - 5) > 2);
var h = -(3 + 2) * -(1 + 1);
var i = -(-(1 + 1));
var j = !(!(!(false)));

var k = 100 + 200 - 50 * 2;
var l = (3 * 3 + 3 * (2 + 3)) / 3;
var m = -(100 - 50);
var n = !((5 + 5) == 10);
var o = ((3 < 4) == false);
var p = ((10 > 5) == true);
var q = ((true == !false) == true);
var r = -(-(-(10)));
var s = !!(!(!!true));
var t = (2 + 2 + 2 + 2 + 2 + 2 + 2) * 0;

a + b + h + i + k + l + m + r + t;
`

func BenchmarkEval_NoConstFold(b *testing.B) {
	runBenchWithoutFolding(b, "NormalInput", foldTestInput)
	runBenchWithoutFolding(b, "LargeInput", foldTestInputLarge)
}

func BenchmarkEval_ConstFold(b *testing.B) {
	runBenchWithFolding(b, "NormalInput", foldTestInput)
	runBenchWithFolding(b, "LargeInput", foldTestInputLarge)
}

func BenchmarkBooleanFoldingMicro(b *testing.B) {
	input := "!(!true)"

	runBenchWithoutFolding(b, "NoFold", input)
	runBenchWithFolding(b, "Fold", input)
}

func runBenchWithoutFolding(b *testing.B, name, input string) {
	b.Run(name, func(b *testing.B) {
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		b.ResetTimer()

		for b.Loop() {
			env := object.NewEnvironment()
			result := Eval(program, env)
			_ = result
		}
	})
}

func runBenchWithFolding(b *testing.B, name, input string) {
	b.Run(name, func(b *testing.B) {
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		perf.ConstFold(program)

		b.ResetTimer()

		for b.Loop() {
			env := object.NewEnvironment()
			result := Eval(program, env)
			_ = result
		}
	})
}
