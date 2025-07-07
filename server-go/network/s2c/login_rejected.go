package s2c

import (
	"grpg/data-go/gbuf"
	"server/shared"
)

type LoginRejected struct {
}

func (l *LoginRejected) Opcode() byte {
	return 0x02
}

func (l *LoginRejected) Handle(buf *gbuf.GBuf, game *shared.Game) {
	// noop
}
