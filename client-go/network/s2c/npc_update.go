package s2c

import (
	"client/shared"
	"client/util"
	"cmp"
	"grpg/data-go/gbuf"
	"log"
)

type NpcUpdate struct {
}

func (n *NpcUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	npcLen, err := buf.ReadUint16()
	if err != nil {
		log.Printf("failed to read uint16 npc len in npc update")
		return
	}
	npcMap := make(map[util.Vector2I]*shared.GameNpc)

	for range npcLen {
		x, err1 := buf.ReadUint32()
		y, err2 := buf.ReadUint32()
		id, err3 := buf.ReadUint16()

		if err := cmp.Or(err1, err2, err3); err != nil {
			log.Printf("failed to read npc in npc update %v\n", err)
			return
		}

		pos := util.Vector2I{X: int32(x), Y: int32(y)}

		npcMap[pos] = &shared.GameNpc{
			Position: pos,
			NpcData:  game.Npcs[id],
		}
	}

	game.TrackedNpcs = npcMap
}
