package s2c

import (
	"client/shared"
	"client/util"
	"cmp"
	"grpg/data-go/gbuf"
	"log"
)

type NpcMove struct {
	From util.Vector2I
	To   util.Vector2I
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
		fromX, err1 := buf.ReadUint32()
		fromY, err2 := buf.ReadUint32()
		toX, err3 := buf.ReadUint32()
		toY, err4 := buf.ReadUint32()
		if err := cmp.Or(err1, err2, err3, err4); err != nil {
			log.Printf("warning: failed to read npc move in npc moves packet")
		}
		moves = append(moves, NpcMove{
			From: util.Vector2I{X: int32(fromX), Y: int32(fromY)},
			To:   util.Vector2I{X: int32(toX), Y: int32(toY)},
		})
	}

	for _, move := range moves {
		// thereotically shouldn't need any checks as TrackedNpcs is kept up solely by the server, which checks this already
		npc, _ := game.TrackedNpcs[move.From]
		game.TrackedNpcs[move.To] = npc
		game.TrackedNpcs[move.To].Position = move.To
		delete(game.TrackedNpcs, move.From)
	}

}
