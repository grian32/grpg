package shared

import (
	"grpg/data-go/gbuf"
	"net"
	"server/network/s2c"
	"server/util"
)

type Game struct {
	Players        []*Player
	PlayersByChunk map[util.Vector2I][]*Player
	// these will be dynamic once map loading is done and as such will be needed
	// for bounds checks.
	MaxX uint32
	MaxY uint32
}

func SendPacket(conn net.Conn, packet s2c.Packet) {
	buf := gbuf.NewEmptyGBuf()
	buf.WriteByte(packet.Opcode())
	packet.Handle(buf)
	_, err := conn.Write(buf.Bytes())
	if err != nil {
		return
	}
}

func (g *Game) FindFirstAvailablePosition() util.Vector2I {
	for y := range g.MaxY + 1 {
		for x := range g.MaxX + 1 {
			pos := util.Vector2I{X: x, Y: y}
			if _, exists := g.Players[pos]; !exists {
				return pos
			}
		}
	}

	return util.Vector2I{X: 0, Y: 0}
}
