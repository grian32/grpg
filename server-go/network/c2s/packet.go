package c2s

import (
	"grpg/data-go/gbuf"
	"server/shared"
	"server/util"
)

type Packet interface {
	Handle(buf *gbuf.GBuf, game *shared.Game, playerPos util.Vector2I)
}

type PacketData struct {
	Opcode  byte
	Length  int16
	Handler Packet
}

var (
	LoginData = PacketData{Opcode: 0x01, Length: -1, Handler: nil}
	MoveData  = PacketData{Opcode: 0x01, Length: 8, Handler: &Move{}}
)

var Packets = map[byte]PacketData{
	0x01: LoginData,
	0x02: MoveData,
}
