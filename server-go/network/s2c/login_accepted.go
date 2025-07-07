package s2c

import (
	"grpg/data-go/gbuf"
)

type LoginAccepted struct {
}

func (l *LoginAccepted) Opcode() byte {
	return 0x01
}

func (l *LoginAccepted) Handle(buf *gbuf.GBuf) {
	// noop
}
