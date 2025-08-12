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
	Length  int16
	Handler Packet
}

var (
	LoginAcceptedData = PacketData{Opcode: 0x01, Length: 0, Handler: &LoginAccepted{}}
	LoginRejectedData = PacketData{Opcode: 0x02, Length: 0, Handler: &LoginRejected{}}
	PlayersUpdateData = PacketData{Opcode: 0x03, Length: -1, Handler: &PlayersUpdate{}}
	ObjUpdateData     = PacketData{Opcode: 0x04, Length: -1, Handler: &ObjUpdate{}}
)

var Packets = map[byte]PacketData{
	0x01: LoginAcceptedData,
	0x02: LoginRejectedData,
	0x03: PlayersUpdateData,
	0x04: ObjUpdateData,
}
