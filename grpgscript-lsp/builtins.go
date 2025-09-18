package grpgscript_lsp

import (
	"grpgscript/object"
)

type TypeTag byte

const (
	INT TypeTag = iota
	STRING
	FUNCTION
)

type BuiltinDefinition struct {
	Name          string
	ArgumentNames []string
	Types         []TypeTag
}

type NamedBuiltin struct {
	Name    string
	Builtin *object.Builtin
}

func MockBuiltin(def BuiltinDefinition, subBuiltins []NamedBuiltin) *object.Builtin {
	return &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != len(def.Types) {
				// todo: err
			}

			for i, want := range def.Types {
				if !matchType(args[i], want) {
					// todo: err
				}
			}

			for _, sub := range subBuiltins {
				env.Set(sub.Name, sub.Builtin)
			}

			return nil
		},
	}
}

func matchType(obj object.Object, want TypeTag) bool {
	switch want {
	case INT:
		_, ok := obj.(*object.Integer)
		return ok
	case STRING:
		_, ok := obj.(*object.String)
		return ok
	case FUNCTION:
		_, ok := obj.(*object.Function)
		return ok
	}

	return false
}
