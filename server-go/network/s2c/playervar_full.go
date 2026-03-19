package s2c

import (
	"grpg/data-go/gbuf"
	"server/shared"
)

type PlayerVarFull struct {
	Player *shared.Player
}

func (p *PlayerVarFull) Opcode() byte {
	return 0x0A
}

func (p *PlayerVarFull) Handle(buf *gbuf.GBuf, game *shared.Game) {
	packetLen := 4 + len(p.Player.PlayerVars)*2 // 4 bytes len + 2*pv size
	buf.WriteUint16(uint16(packetLen))
	buf.WriteUint32(uint32(len(p.Player.PlayerVars)))
	for _, pv := range p.Player.PlayerVars {
		buf.WriteUint16(pv)
	}
}
