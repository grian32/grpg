package evaluator

import (
	"fmt"
	"grpgscript/ast"
	"grpgscript/object"
)

var NULL = &object.Null{}

var boolLookup = map[bool]*object.Boolean{
	true:  &object.Boolean{Value: true},
	false: &object.Boolean{Value: false},
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return boolLookup[node.Value]
	default:
		panic(fmt.Sprintf("unexpected ast.Node: %#v", node))
	}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}
