package perf

import (
	"fmt"
	"grpgscript/ast"
)

func ConstFold(program *ast.Program) {
	for _, stmt := range program.Statements {
		foldStmt(stmt)
	}
}

func foldStmt(stmt ast.Statement) ast.Statement {
	switch s := stmt.(type) {
	case *ast.BlockStatement:
		for idx, bStmt := range s.Statements {
			s.Statements[idx] = foldStmt(bStmt)
		}
	case *ast.ExpressionStatement:
		s.Expression = foldExpr(s.Expression)
		return s
	case *ast.ReturnStatement:
		s.ReturnValue = foldExpr(s.ReturnValue)
		return s
	case *ast.VarStatement:
		s.Value = foldExpr(s.Value)
		return s
	default:
		panic(fmt.Sprintf("unexpected ast.Statement: %#v", s))
	}

	return nil
}

func foldExpr(expr ast.Expression) ast.Expression {
	switch e := expr.(type) {
	case *ast.InfixExpression:
		left := foldExpr(e.Left)
		right := foldExpr(e.Right)
		e.Left = left
		e.Right = right

		leftInt, leftIntOk := left.(*ast.IntegerLiteral)
		rightInt, rightIntOk := right.(*ast.IntegerLiteral)

		leftBool, leftBoolOk := left.(*ast.Boolean)
		rightBool, rightBoolOk := right.(*ast.Boolean)

		if leftIntOk && rightIntOk {
			switch e.Operator {
			case "+":
				return &ast.IntegerLiteral{Value: leftInt.Value + rightInt.Value}
			case "*":
				return &ast.IntegerLiteral{Value: leftInt.Value * rightInt.Value}
			case "/":
				if rightInt.Value != 0 {
					return &ast.IntegerLiteral{Value: leftInt.Value / rightInt.Value}
				}
			case "-":
				return &ast.IntegerLiteral{Value: leftInt.Value - rightInt.Value}
			case "<":
				return &ast.Boolean{Value: leftInt.Value < rightInt.Value}
			case ">":
				return &ast.Boolean{Value: leftInt.Value > rightInt.Value}
			case "==":
				return &ast.Boolean{Value: leftInt.Value == rightInt.Value}
			case "!=":
				return &ast.Boolean{Value: leftInt.Value != rightInt.Value}
			}
		}

		// this should also handle stuff like (2 < 3) != false,
		// since 2 < 3 gets folded down to true, and then true != false gets handled by this
		if leftBoolOk && rightBoolOk {
			switch e.Operator {
			case "==":
				return &ast.Boolean{Value: leftBool.Value == rightBool.Value}
			case "!=":
				return &ast.Boolean{Value: leftBool.Value != rightBool.Value}
			}
		}
	case *ast.PrefixExpression:
		right := foldExpr(e.Right)
		e.Right = right

		switch e.Operator {
		case "-":
			if intLit, ok := right.(*ast.IntegerLiteral); ok {
				return &ast.IntegerLiteral{Value: -intLit.Value}
			}
		case "!":
			if boolLit, ok := right.(*ast.Boolean); ok {
				return &ast.Boolean{Value: !boolLit.Value}
			}
		}
	case *ast.CallExpresion:
		for idx, arg := range e.Arguments {
			e.Arguments[idx] = foldExpr(arg)
		}

		if ident, ok := e.Function.(*ast.Identifier); ok && len(e.Arguments) == 1 {
			switch ident.Value {
			case "len":
				switch arg := e.Arguments[0].(type) {
				case *ast.StringLiteral:
					return &ast.IntegerLiteral{Value: int64(len(arg.Value))}
				case *ast.ArrayLiteral:
					return &ast.IntegerLiteral{Value: int64(len(arg.Elements))}
				}
			}
		}

		return e
	default:
		return expr
	}

	return expr
}
