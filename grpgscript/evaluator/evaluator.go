package evaluator

import (
	"fmt"
	"grpgscript/ast"
	"grpgscript/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)
var boolLookup = map[bool]*object.Boolean{
	true:  TRUE,
	false: FALSE,
}

type EvalError struct {
	Msg      string
	Position ast.Position
}

type ErrorStore struct {
	Errors []EvalError
}

type Evaluator struct {
	ErrorStore *ErrorStore
}

func NewEvaluator() *Evaluator {
	return &Evaluator{ErrorStore: &ErrorStore{Errors: make([]EvalError, 0)}}
}

func (e *Evaluator) Eval(node ast.Node, env *object.Environment) object.Object {
	// TODO: error returning in this on eval is kinda meh

	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return boolLookup[node.Value]
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.PrefixExpression:
		right := e.Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return e.evalPrefixExpression(node.Operator, right, node.Pos())
	case *ast.InfixExpression:
		left := e.Eval(node.Left, env)
		right := e.Eval(node.Right, env)
		if isError(left) {
			return left
		}
		if isError(right) {
			return right
		}
		return e.evalInfixExpression(node.Operator, left, right, node.Pos())
	case *ast.BlockStatement:
		return e.evalBlockStatement(node, env)
	case *ast.IfExpression:
		return e.evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := e.Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.VarStatement:
		val := e.Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return e.evalIdentifier(node, env, node.Pos())
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpresion:
		function := e.Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := e.evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return e.applyFunction(function, args, env, node.Pos())
	case *ast.ArrayLiteral:
		elements := e.evalArrayExpressions(node.Elements, env, node.Pos())

		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := e.Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := e.Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return e.evalIndexExpression(left, index, node.Pos())
	case *ast.HashLiteral:
		return e.evalHashLiteral(node, env, node.Pos())
	default:
		panic(fmt.Sprintf("unexpected ast.Node: %#v", node))
	}

	return nil
}

func (e *Evaluator) evalHashLiteral(node *ast.HashLiteral, env *object.Environment, pos ast.Position) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := e.Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashable, ok := key.(object.Hashable)
		if !ok {
			return e.ErrorStore.NewError(pos, "unusable as hash key: %s", key.Type())
		}

		val := e.Eval(valueNode, env)
		if isError(val) {
			return key
		}

		hashKey := hashable.HashKey()

		pairs[hashKey] = object.HashPair{
			Key:   key,
			Value: val,
		}
	}

	return &object.Hash{Pairs: pairs}
}

func (e *Evaluator) evalIndexExpression(left, index object.Object, pos ast.Position) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return e.evalArrayIndexExpression(left, index, pos)
	case left.Type() == object.HASH_OBJ:
		return e.evalHashIndexExpression(left, index, pos)
	default:
		return e.ErrorStore.NewError(pos, "index operator not supported: %s", left.Type())
	}
}

func (e *Evaluator) evalHashIndexExpression(left, index object.Object, pos ast.Position) object.Object {
	hash := left.(*object.Hash)

	hashableIdx, ok := index.(object.Hashable)
	if !ok {
		return e.ErrorStore.NewError(pos, "unusable as hash key: %s", index.Type())
	}

	val, ok := hash.Pairs[hashableIdx.HashKey()]
	if !ok {
		return e.ErrorStore.NewError(pos, "unknown hash key")
	}

	return val.Value
}

func (e *Evaluator) evalArrayIndexExpression(left, index object.Object, pos ast.Position) object.Object {
	arr := left.(*object.Array)
	idx := index.(*object.Integer).Value

	max := int64(len(arr.Elements) - 1)

	if idx < 0 || idx > max {
		return e.ErrorStore.NewError(pos, "index %d out of bounds on array of size %d", idx, len(arr.Elements))
	}

	return arr.Elements[idx]
}

func (e *Evaluator) applyFunction(fn object.Object, args []object.Object, env *object.Environment, pos ast.Position) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := e.extendFunctionEnv(fn, args)
		evaluated := e.Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(env, args...)
	default:
		return e.ErrorStore.NewError(pos, "not a function, %s", fn.Type())
	}
}

func (e *Evaluator) extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvinronment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func (e *Evaluator) evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expr := range exps {
		evaluated := e.Eval(expr, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func (e *Evaluator) evalArrayExpressions(exps []ast.Expression, env *object.Environment, pos ast.Position) []object.Object {
	if len(exps) == 0 {
		return []object.Object{}
	}

	first := e.Eval(exps[0], env)

	result := []object.Object{first}
	firstType := first.Type()

	for _, expr := range exps[1:] {
		evaluated := e.Eval(expr, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		// todo: see if this has significant perf impact or not.. if it does i'd rather just write code with multi type arrays and just not use it.
		if evaluated.Type() != firstType {
			return []object.Object{e.ErrorStore.NewError(pos, "arrays can only be of the same type, found type %s, in array of type %s", evaluated.Type(), firstType)}
		}
		result = append(result, evaluated)
	}

	return result
}

func (e *Evaluator) evalIdentifier(node *ast.Identifier, env *object.Environment, pos ast.Position) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return e.ErrorStore.NewError(pos, "identifier not found: %s", node.Value)
}

func (e *Evaluator) evalBlockStatement(node *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range node.Statements {
		result = e.Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func (e *Evaluator) evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := e.Eval(node.Condition, env)

	if isError(condition) {
		return condition
	}

	if condition == TRUE {
		return e.Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return e.Eval(node.Alternative, env)
	} else {
		return NULL
	}
}

func (e *Evaluator) evalInfixExpression(operator string, left, right object.Object, pos ast.Position) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return e.evalIntegerInfixExpression(operator, left, right, pos)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return e.evalStringInfixExpression(operator, left, right, pos)
	case operator == "==":
		return boolLookup[left == right]
	case operator == "!=":
		return boolLookup[left != right]
	case left.Type() != right.Type():
		return e.ErrorStore.NewError(pos, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return e.ErrorStore.NewError(pos, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (e *Evaluator) evalStringInfixExpression(operator string, left, right object.Object, pos ast.Position) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	case "==":
		return &object.Boolean{Value: leftValue == rightValue}
	case "!=":
		return &object.Boolean{Value: leftValue != rightValue}
	default:
		return e.ErrorStore.NewError(pos, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (e *Evaluator) evalIntegerInfixExpression(operator string, left, right object.Object, pos ast.Position) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "<":
		return boolLookup[leftValue < rightValue]
	case ">":
		return boolLookup[leftValue > rightValue]
	case "!=":
		return boolLookup[leftValue != rightValue]
	case "==":
		return boolLookup[leftValue == rightValue]
	default:
		return e.ErrorStore.NewError(pos, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (e *Evaluator) evalPrefixExpression(operator string, right object.Object, pos ast.Position) object.Object {
	switch operator {
	case "!":
		return e.evalBangOperatorExpression(right, pos)
	case "-":
		return e.evalMinusPrefixOperatorExpression(right, pos)
	default:
		return e.ErrorStore.NewError(pos, "unknown operator: %s%s", operator, right.Type())
	}
}

func (e *Evaluator) evalBangOperatorExpression(right object.Object, pos ast.Position) object.Object {
	if right.Type() != object.BOOLEAN_OBJ {
		return e.ErrorStore.NewError(pos, "unknown operator: !%s", right.Type())
	}

	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		// unreacheable theoretically
		return NULL
	}
}

func (e *Evaluator) evalMinusPrefixOperatorExpression(right object.Object, pos ast.Position) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return e.ErrorStore.NewError(pos, "unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func (e *Evaluator) evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = e.Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func (e *ErrorStore) NewError(pos ast.Position, format string, a ...any) *object.Error {
	msg := fmt.Sprintf(format, a...)
	e.Errors = append(e.Errors, EvalError{
		Msg:      msg,
		Position: pos,
	})

	return &object.Error{Message: msg}
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
