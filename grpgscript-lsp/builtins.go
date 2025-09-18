package grpgscript_lsp

import (
	"grpgscript/ast"
	"grpgscript/evaluator"
	"grpgscript/object"
	"strings"
)

type TypeTag byte

const (
	INT TypeTag = iota
	STRING
	FUNCTION
	NULL
)

func (tt TypeTag) String() string {
	switch tt {
	case INT:
		return "INT"
	case STRING:
		return "STRING"
	case FUNCTION:
		return "FUNCTION"
	case NULL:
		return "NULL"
	}
	return ""
}

type BuiltinDefinition struct {
	Name          string
	ArgumentNames []string
	Types         []TypeTag
	Label         string
	ReturnType    TypeTag
}

func NewBuiltinDefinition(name string, argNames []string, types []TypeTag, returnType TypeTag) BuiltinDefinition {
	if len(argNames) != len(types) {
		panic("builtin definition cannot have mismatched argname and type len")
	}

	var b strings.Builder

	b.WriteString(name)
	b.WriteByte('(')

	for i, s := range argNames {
		b.WriteString(s)
		b.WriteByte(' ')
		b.WriteString(types[i].String())

		if i < len(argNames)-1 {
			b.WriteString(", ")
		}
	}

	b.WriteByte(')')

	return BuiltinDefinition{
		Name:          name,
		ArgumentNames: argNames,
		Types:         types,
		Label:         b.String(),
		ReturnType:    returnType,
	}
}

type NamedBuiltin struct {
	Name    string
	Builtin *object.Builtin
}

func NewNamedBuiltin(name string, builtin *object.Builtin) NamedBuiltin {
	return NamedBuiltin{
		Name:    name,
		Builtin: builtin,
	}
}

func MockBuiltin(def BuiltinDefinition, subBuiltins []NamedBuiltin) *object.Builtin {
	return &object.Builtin{
		Fn: func(env *object.Environment, pos ast.Position, errorStore *object.ErrorStore, args ...object.Object) object.Object {
			if len(args) != len(def.Types) {
				errorStore.NewError(pos, "got %d arguments for %s, want %d arguments", len(args), def.Name, len(def.Types))
				return getReturn(def.ReturnType)
			}

			if subBuiltins != nil {
				for _, sub := range subBuiltins {
					env.Set(sub.Name, sub.Builtin)
				}
			}

			for i, want := range def.Types {
				eval := &evaluator.Evaluator{ErrorStore: errorStore}
				if !matchType(args[i], want, env, eval) {
					errorStore.NewError(pos, "arg %s for %s is not of type %s", def.ArgumentNames[i], def.Name, def.Types[i].String())
					return getReturn(def.ReturnType)
				}
			}

			return getReturn(def.ReturnType)
		},
	}
}

func matchType(obj object.Object, want TypeTag, env *object.Environment, eval *evaluator.Evaluator) bool {
	switch want {
	case INT:
		_, ok := obj.(*object.Integer)
		return ok
	case STRING:
		_, ok := obj.(*object.String)
		return ok
	case FUNCTION:
		fnc, ok := obj.(*object.Function)
		execAll(fnc.Body, eval, env)
		return ok
	case NULL:
		return true
	}

	return false
}

func execAll(block *ast.BlockStatement, eval *evaluator.Evaluator, env *object.Environment) {
	for _, stmt := range block.Statements {
		expr, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			return
		}
		if ifExpr, ok := expr.Expression.(*ast.IfExpression); ok {
			execAll(ifExpr.Consequence, eval, env)
			if ifExpr.Alternative != nil {
				execAll(ifExpr.Alternative, eval, env)
			}
		} else {
			eval.Eval(expr, env)
		}
	}
}

func getReturn(want TypeTag) object.Object {
	switch want {
	case INT:
		return &object.Integer{Value: -1}
	case STRING:
		return &object.String{Value: ""}
	case FUNCTION:
		return nil
	case NULL:
		return nil
	}
	return nil
}
