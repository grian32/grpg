package scripts

import (
	"grpgscript/ast"
	"grpgscript/object"
	"server/util"
)

type UpdateType byte

const (
	OBJECT UpdateType = iota
)

// TimedScript TODO: improve the updating here, maybe tie the update to the builtin? rather than having it done in the interact packet
type TimedScript struct {
	Script *ast.BlockStatement
	Env    *object.Environment
	Update UpdateType
	// this is not used unless relevant to the update
	ChunkPos util.Vector2I
}
