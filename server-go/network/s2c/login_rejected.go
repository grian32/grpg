package s2c

import "grpg/data-go/gbuf"

type LoginRejected struct {
}

func (l *LoginRejected) Opcode() byte {
	return 0x02
}

func (l *LoginRejected) Handle(buf *gbuf.GBuf) {
	// noop
}
