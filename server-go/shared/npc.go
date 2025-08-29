package shared

import (
	"grpg/data-go/grpgnpc"
	"server/util"
)

type GameNpc struct {
	Pos      util.Vector2I
	NpcData  *grpgnpc.Npc
	ChunkPos util.Vector2I
}
