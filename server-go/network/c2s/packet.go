package c2s

import (
	"grpg/data-go/gbuf"
	"server/scripts"
	"server/shared"
)

type Packet interface {
	Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player, scriptManager *scripts.ScriptManager)
}

type PacketData struct {
	Opcode  byte
	Length  int16
	Handler Packet
}

var (
	LoginData    = PacketData{Opcode: 0x01, Length: -1, Handler: nil}
	MoveData     = PacketData{Opcode: 0x02, Length: 9, Handler: &Move{}}
	InteractData = PacketData{Opcode: 0x03, Length: 10, Handler: &Interact{}}
)

var Packets = map[byte]PacketData{
	0x01: LoginData,
	0x02: MoveData,
	0x03: InteractData,
}
