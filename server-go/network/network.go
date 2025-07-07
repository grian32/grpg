package network

import (
	"fmt"
	"grpg/data-go/gbuf"
	"net"
	"server/shared"
)
import "server/util"
import "server/network/s2c"

func UpdatePlayersByChunk(chunkPos util.Vector2I, game *shared.Game) {
	packet := &s2c.PlayersUpdate{
		ChunkPos: chunkPos,
	}

	for _, player := range game.PlayersByChunk[chunkPos] {
		SendPacket(player.Conn, packet, game)
	}
}

func SendPacket(conn net.Conn, packet s2c.Packet, game *shared.Game) {
	buf := gbuf.NewEmptyGBuf()
	buf.WriteByte(packet.Opcode())
	packet.Handle(buf, game)
	fmt.Println(buf.Bytes())
	_, err := conn.Write(buf.Bytes())
	if err != nil {
		return
	}
}
