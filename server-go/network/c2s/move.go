package c2s

import (
	"cmp"
	"grpg/data-go/gbuf"
	"log"
	"server/network"
	"server/network/s2c"
	"server/scripts"
	"server/shared"
	"server/util"
)

type Move struct {
}

func (m *Move) Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player, scriptManagers *scripts.ScriptManager) {
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

	prevChunkPos := player.ChunkPos

	chunkPos := util.Vector2I{X: newX / 16, Y: newY / 16}

	crossedZone := chunkPos.X != player.ChunkPos.X || chunkPos.Y != player.ChunkPos.Y

	player.Pos.X = newX
	player.Pos.Y = newY
	player.ChunkPos = chunkPos
	player.Facing = shared.Direction(facing)

	network.UpdatePlayersByChunk(chunkPos, game, &s2c.PlayersUpdate{ChunkPos: chunkPos})
	if player.DialogueQueue.MaxIndex > 0 {
		player.DialogueQueue.Clear()
		// rest doesnt matter just let default init here
	}
	network.SendPacket(player.Conn, &s2c.Talkbox{Type: s2c.CLEAR}, game)
	if crossedZone {
		network.UpdatePlayersByChunk(prevChunkPos, game, &s2c.PlayersUpdate{ChunkPos: prevChunkPos})
		network.SendPacket(player.Conn, &s2c.ObjUpdate{ChunkPos: chunkPos, Rebuild: true}, game)
		network.SendPacket(player.Conn, &s2c.NpcUpdate{ChunkPos: chunkPos}, game)
	}
}
