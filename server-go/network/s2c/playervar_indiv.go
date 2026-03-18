package s2c

import (
	"grpg/data-go/gbuf"
	"server/shared"
)

type PlayerVarIndiv struct {
	VarId uint16
	VarValue uint16
}

func (p *PlayerVarIndiv) Opcode() byte {
	return 0x0B
}

func (p *PlayerVarIndiv) Handle(buf *gbuf.GBuf, game *shared.Game) {
	buf.WriteUint16(p.VarId)
	buf.WriteUint16(p.VarValue)
}
