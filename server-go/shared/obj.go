package shared

import (
	"grpg/data-go/grpgobj"
	"server/util"
)

type GameObj struct {
	ObjData grpgobj.Obj
	// this is mainly so i don't have to compute this in packet handling, it's not strictly necessary info to have here
	// & i already have the chunk pos of the object from the zone header when adding it to the tracked map
	ChunkPos util.Vector2I
	State    byte
}
