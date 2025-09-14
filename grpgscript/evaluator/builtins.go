package evaluator

import (
	"fmt"
	"grpgscript/object"
	"slices"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"println": {
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
	"push": {
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			return pushUnshift(PUSH, args...)
		},
	},
	"unshift": {
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			return pushUnshift(UNSHIFT, args...)
		},
	},
	"concat": {
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of argument, got=%d, want=2", len(args))
			}

			firstArrArg, ok1 := args[0].(*object.Array)
			secondArrArg, ok2 := args[1].(*object.Array)

			if len(firstArrArg.Elements) == 0 && len(secondArrArg.Elements) == 0 {
				return &object.Array{Elements: []object.Object{}}
			}

			if len(firstArrArg.Elements) == 0 {
				return secondArrArg
			}

			if len(secondArrArg.Elements) == 0 {
				return firstArrArg
			}

			if !(ok1 && ok2) {
				return newError("one or both of the arguments to concat are not arrays")
			}

			// it's fine to check only first elem since arrays are guaranteed to be of the same type on all elems due to eval
			if firstArrArg.Elements[0].Type() != secondArrArg.Elements[0].Type() {
				return newError("both arrays passed to concat must be of the same element type")
			}

			concattedElems := slices.Concat(firstArrArg.Elements, secondArrArg.Elements)

			newArr := &object.Array{Elements: concattedElems}

			return newArr
		},
	},
}

type PushUnshift byte

const (
	PUSH PushUnshift = iota
	UNSHIFT
)

func pushUnshift(use PushUnshift, args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}

	arrayArg, ok := args[0].(*object.Array)
	if !ok {
		return newError("first arg is not arr, got=%T(%+v)", args[0], args[0])
	}
	itemArg := args[1]

	if len(arrayArg.Elements) == 0 {
		arrayArg.Elements = []object.Object{itemArg}
		return &object.Integer{Value: 1}
	}

	arrType := arrayArg.Elements[0].Type()
	if itemArg.Type() != arrType {
		return newError("cannot add element of type %s to array of type %s", itemArg.Type(), arrType)
	}

	if use == PUSH {
		arrayArg.Elements = append(arrayArg.Elements, itemArg)
	} else {
		arrayArg.Elements = append([]object.Object{itemArg}, arrayArg.Elements...)
	}

	return &object.Integer{Value: int64(len(arrayArg.Elements))}
}
