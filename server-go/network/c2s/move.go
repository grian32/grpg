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

func (m *Move) Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player) {
	newX, err1 := buf.ReadUint32()
	newY, err2 := buf.ReadUint32()
	facing, err3 := buf.ReadByte()

	if err := cmp.Or(err1, err2, err3); err != nil {
		log.Printf("Failed to read move packet: %v\n", err)
		return
	}

	_, exists := game.CollisionMap[util.Vector2I{X: newX, Y: newY}]
	if newX > game.MaxX || newY > game.MaxY || exists || facing > 3 {
		return
	}

	chunkPos := util.Vector2I{X: newX / 16, Y: newY / 16}

	player.Pos.X = newX
	player.Pos.Y = newY
	player.ChunkPos = chunkPos
	player.Facing = shared.Direction(facing)

	network.UpdatePlayersByChunk(chunkPos, game)
}
