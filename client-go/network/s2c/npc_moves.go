package s2c

import (
	"client/shared"
	"client/util"
	"grpg/data-go/gbuf"
	"log"

)

type NpcMove struct {
	NpcUid uint32
	Move util.Vector2I
}

type NpcMoves struct {
}

func (n *NpcMoves) Handle(buf *gbuf.GBuf, game *shared.Game) {
	length, err := buf.ReadUint32()
	if err != nil {
		log.Printf("warning: failed to read length in npc moves packet")
	}
	// NOTE: it might be worth just doing the moves inline? :shrug:
	moves := make([]NpcMove, 0, length)
	for _ = range length {
		npcUid, err := buf.ReadUint32()
		if err != nil {
			log.Printf("warning: failed to read npc uid in npc moves packet\n")
			continue
		}
		toX, err := buf.ReadUint32()
		if err != nil {
			log.Printf("warning: failed to read npc move x in npc moves packet\n")
			continue
		}
		toY, err := buf.ReadUint32()
		if err != nil {
			log.Printf("warning: failed to read npc move y in npc moves packet\n")
			continue
		}
		moves = append(moves, NpcMove{
			NpcUid: npcUid,
			Move:   util.Vector2I{X: int32(toX), Y: int32(toY)},
		})
	}

	for _, move := range moves {
		// thereotically shouldn't need any checks as TrackedNpcs is kept up solely by the server, which checks this already
		npc, _ := game.TrackedNpcs[move.NpcUid]
		delete(game.NpcsByPos, npc.Position)
		game.TrackedNpcs[npc.Uid].Position = move.Move
		game.NpcsByPos[move.Move] = npc
	}
}
