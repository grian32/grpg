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
	Type TalkboxType
	Msg  string
}

func (t *Talkbox) Opcode() byte {
	return 0x07
}

func (t *Talkbox) Handle(buf *gbuf.GBuf, game *shared.Game) {
	packetLen := 1 // type byte

	if t.Type != CLEAR {
		packetLen += 4 + len(t.Msg) // uint32 len + string len
	}

	buf.WriteUint16(uint16(packetLen))
	buf.WriteByte(byte(t.Type))
	if t.Type != CLEAR {
		buf.WriteString(t.Msg)
	}
}
