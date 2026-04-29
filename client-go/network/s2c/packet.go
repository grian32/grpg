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
	LoginAcceptedData   = PacketData{Opcode: 0x01, Length: 0, Handler: &LoginAccepted{}}
	LoginRejectedData   = PacketData{Opcode: 0x02, Length: 0, Handler: &LoginRejected{}}
	PlayersUpdateData   = PacketData{Opcode: 0x03, Length: -1, Handler: &PlayersUpdate{}}
	ObjUpdateData       = PacketData{Opcode: 0x04, Length: -1, Handler: &ObjUpdate{}}
	InventoryUpdateData = PacketData{Opcode: 0x05, Length: -1, Handler: &InventoryUpdate{}}
	NpcUpdateData       = PacketData{Opcode: 0x06, Length: -1, Handler: &NpcUpdate{}}
	TalkboxData         = PacketData{Opcode: 0x07, Length: -1, Handler: &Talkbox{}}
	SkillUpdateData     = PacketData{Opcode: 0x08, Length: 7, Handler: &SkillUpdate{}}
	NpcMovesData        = PacketData{Opcode: 0x09, Length: -1, Handler: &NpcMoves{}}
	PlayerVarFullData   = PacketData{Opcode: 0x0A, Length: -1, Handler: &PlayerVarFull{}}
	PlayerVarIndivData  = PacketData{Opcode: 0x0B, Length: 4, Handler: &PlayerVarIndiv{}}
	StatUpdateData      = PacketData{Opcode: 0x0C, Length: 1, Handler: &StatUpdate{}}
)

var Packets = map[byte]PacketData{
	0x01: LoginAcceptedData,
	0x02: LoginRejectedData,
	0x03: PlayersUpdateData,
	0x04: ObjUpdateData,
	0x05: InventoryUpdateData,
	0x06: NpcUpdateData,
	0x07: TalkboxData,
	0x08: SkillUpdateData,
	0x09: NpcMovesData,
	0x0A: PlayerVarFullData,
	0x0B: PlayerVarIndivData,
	0x0C: StatUpdateData,
}
