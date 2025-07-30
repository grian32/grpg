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

		leftInt, leftOk := left.(*ast.IntegerLiteral)
		rightInt, rightOk := right.(*ast.IntegerLiteral)

		if leftOk && rightOk {
			switch e.Operator {
			case "+":
				return &ast.IntegerLiteral{Value: leftInt.Value + rightInt.Value}
			case "*":
				return &ast.IntegerLiteral{Value: leftInt.Value * rightInt.Value}
			case "/":
				return &ast.IntegerLiteral{Value: leftInt.Value / rightInt.Value}
			case "-":
				return &ast.IntegerLiteral{Value: leftInt.Value - rightInt.Value}
			}
		}

	case *ast.PrefixExpression:
		right := foldExpr(e.Right)
		e.Right = right

		switch e.Operator {
		case "-":
			if intLit, ok := right.(*ast.IntegerLiteral); ok {
				return &ast.IntegerLiteral{Value: -intLit.Value}
			} else if boolLit, ok := right.(*ast.Boolean); ok {
				return &ast.Boolean{Value: !boolLit.Value}
			}
		case "!":

		}
	default:
		return expr
	}

	return expr
}
