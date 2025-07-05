package c2s

import "grpg/data-go/gbuf"

type LoginPacket struct {
	PlayerName string
}

func (l *LoginPacket) Opcode() byte {
	return 0x01
}

func (l *LoginPacket) Handle(buf *gbuf.GBuf) {
	buf.WriteString(l.PlayerName)
}
