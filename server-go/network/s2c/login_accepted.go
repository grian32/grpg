package s2c

import (
	"grpg/data-go/gbuf"
)

type LoginAccepted struct {
	InitialX int32
	InitialY int32
}

func (l *LoginAccepted) Opcode() byte {
	return 0x01
}

func (l *LoginAccepted) Handle(buf *gbuf.GBuf) {
	buf.WriteInt32(l.InitialX)
	buf.WriteInt32(l.InitialY)
}
