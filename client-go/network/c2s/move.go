package c2s

import "grpg/data-go/gbuf"

type MovePacket struct {
	X int32
	Y int32
}

func (m *MovePacket) Opcode() byte {
	return 0x02
}

func (m *MovePacket) Handle(buf *gbuf.GBuf) {
	buf.WriteInt32(m.X)
	buf.WriteInt32(m.Y)
}
