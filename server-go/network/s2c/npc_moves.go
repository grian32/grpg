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
	packetLen := 4 + len(n.Moves)*4*3 // so thats 4 for the length and the 4*3 for 3 uints32, 2 move pos, one uid
	buf.WriteUint16(uint16(packetLen))
	buf.WriteUint32(uint32(len(n.Moves)))
	for _, m := range n.Moves {
		buf.WriteUint32(m.NpcUid)
		buf.WriteUint32(m.Move.X)
		buf.WriteUint32(m.Move.Y)
	}
}
