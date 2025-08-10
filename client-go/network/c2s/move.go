package c2s

import "grpg/data-go/gbuf"

type MovePacket struct {
	X      uint32
	Y      uint32
	Facing byte
}

func (m *MovePacket) Opcode() byte {
	return 0x02
}

func (m *MovePacket) Handle(buf *gbuf.GBuf) {
	buf.WriteUint32(m.X)
	buf.WriteUint32(m.Y)
	buf.WriteByte(m.Facing)
}
