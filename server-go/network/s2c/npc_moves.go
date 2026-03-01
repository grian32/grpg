package s2c

import (
	"grpg/data-go/gbuf"
	"server/shared"
)

type NpcMoves struct {
	Moves []shared.NpcMove
}

func (n *NpcMoves) Opcode() byte {
	return 0x09
}

func (n *NpcMoves) Handle(buf *gbuf.GBuf, game *shared.Game) {
	packetLen := 4 + len(n.Moves)*4*4 // so thats 4 for the length and the 4*4 for 4 uints32 * 4 for each move
	buf.WriteUint16(uint16(packetLen))
	buf.WriteUint32(uint32(len(n.Moves)))
	for _, m := range n.Moves {
		buf.WriteUint32(m.From.X)
		buf.WriteUint32(m.From.Y)
		buf.WriteUint32(m.To.X)
		buf.WriteUint32(m.To.Y)
	}
}
