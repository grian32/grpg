package s2c

import (
	"grpg/data-go/gbuf"
	"server/shared"
)

type TalkboxType byte

const (
	PLAYER TalkboxType = iota
	NPC
	CLEAR
)

type Talkbox struct {
	Type  TalkboxType
	Msg   string
	NpcId uint16
}

func (t *Talkbox) Opcode() byte {
	return 0x07
}

func (t *Talkbox) Handle(buf *gbuf.GBuf, game *shared.Game) {
	packetLen := 1 // type byte

	if t.Type == PLAYER {
		packetLen += 4 + len(t.Msg) // uint32 len + string len
	} else if t.Type == NPC {
		packetLen += 2 + 4 + len(t.Msg) // uint16 npc id + uint32 len + string len
	}

	buf.WriteUint16(uint16(packetLen))
	buf.WriteByte(byte(t.Type))
	if t.Type == PLAYER {
		buf.WriteString(t.Msg)
	} else if t.Type == NPC {
		buf.WriteUint16(t.NpcId)
		buf.WriteString(t.Msg)
	}
}
