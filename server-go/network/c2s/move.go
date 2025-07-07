package c2s

import (
	"cmp"
	"grpg/data-go/gbuf"
	"log"
	"server/shared"
)

type Move struct {
}

func (m *Move) Handle(buf *gbuf.GBuf, game *shared.Game, playerPos int) {
	newX, err1 := buf.ReadInt32()
	newY, err2 := buf.ReadInt32()

	if err := cmp.Or(err1, err2); err != nil {
		log.Printf("Failed to read move packet: %v\n", err)
		return
	}

	player := game.Players[playerPos]
	player.Pos.X = uint32(newX)
	player.Pos.Y = uint32(newY)
}
