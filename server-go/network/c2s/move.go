package c2s

import (
	"cmp"
	"grpg/data-go/gbuf"
	"log"
	"server/network"
	"server/shared"
	"server/util"
)

type Move struct {
}

func (m *Move) Handle(buf *gbuf.GBuf, game *shared.Game, playerPos int) {
	newX, err1 := buf.ReadUint32()
	newY, err2 := buf.ReadUint32()

	if err := cmp.Or(err1, err2); err != nil {
		log.Printf("Failed to read move packet: %v\n", err)
		return
	}

	_, exists := game.CollisionMap[util.Vector2I{X: newX, Y: newY}]
	if newX > game.MaxX || newX < 0 || newY > game.MaxY || newY < 0 || exists {
		return
	}

	chunkPos := util.Vector2I{X: newX / 16, Y: newY / 16}

	player := game.Players[playerPos]
	player.Pos.X = newX
	player.Pos.Y = newY
	player.ChunkPos = chunkPos

	network.UpdatePlayersByChunk(chunkPos, game)
}
