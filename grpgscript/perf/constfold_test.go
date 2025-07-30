package perf

import (
	"grpgscript/ast"
	"grpgscript/lexer"
	"grpgscript/parser"
	"testing"
)

func TestConstFold(t *testing.T) {
	tests := map[string]any{
		"(1 + 2) * 3": int64(9),
		"-(1+1)":      int64(-2),
		"!true":       false,
		"!(!true)":    true,
	}

	for input, expected := range tests {
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()
		ConstFold(program)

		exprStmt := program.Statements[0].(*ast.ExpressionStatement)
		switch e := expected.(type) {
		case int64:
			lit := exprStmt.Expression.(*ast.IntegerLiteral)
			if lit.Value != expected {
				t.Errorf("lit.Value=%d, want=%d, input=%s", lit.Value, e, input)
			}
		case bool:
			lit := exprStmt.Expression.(*ast.Boolean)
			if lit.Value != expected {
				t.Errorf("lit.Value=%t, want=%t, input=%s", lit.Value, e, input)
			}
		}
	}
}
