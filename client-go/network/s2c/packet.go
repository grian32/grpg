package s2c

import (
	"client/shared"
	"grpg/data-go/gbuf"
)

type Packet interface {
	Handle(buf *gbuf.GBuf, game *shared.Game)
}

type PacketData struct {
	Opcode  byte
	Length  uint16
	Handler Packet
}

var (
	LoginAcceptedData = PacketData{Opcode: 0x01, Length: 8, Handler: &LoginAccepted{}}
)

var Packets = map[byte]PacketData{
	0x01: LoginAcceptedData,
}
