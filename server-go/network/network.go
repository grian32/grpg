package network

import (
	"grpg/data-go/gbuf"
	"net"
	"server/network/s2c"
	"server/shared"
	"server/util"
)

func UpdatePlayersByChunk(chunkPos util.Vector2I, game *shared.Game) {
	packet := &s2c.PlayersUpdate{
		ChunkPos: chunkPos,
	}

	for player, _ := range game.Players {
		if player.ChunkPos == chunkPos {
			SendPacket(player.Conn, packet, game)
		}
	}
}

func SendPacket(conn net.Conn, packet s2c.Packet, game *shared.Game) {
	buf := gbuf.NewEmptyGBuf()
	buf.WriteByte(packet.Opcode())
	packet.Handle(buf, game)
	_, err := conn.Write(buf.Bytes())
	if err != nil {
		return
	}
}
