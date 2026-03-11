package shared

import (
	"client/util"
	"grpg/data-go/grpgnpc"
)

type GameNpc struct {
	Position util.Vector2I
	Uid      uint32
	NpcData  *grpgnpc.Npc
}
