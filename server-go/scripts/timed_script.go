package scripts

import (
	"grpgscript/ast"
	"grpgscript/object"
)

type TimedScript struct {
	Script *ast.BlockStatement
	Env    *object.Environment
}
