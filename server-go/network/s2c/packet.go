package s2c

import (
	"grpg/data-go/gbuf"
	"server/shared"
)

type Packet interface {
	Opcode() byte
	Handle(buf *gbuf.GBuf, game *shared.Game)
}
