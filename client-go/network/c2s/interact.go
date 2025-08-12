package c2s

import "grpg/data-go/gbuf"

type InteractPacket struct {
	ObjId uint16
	X     uint32
	Y     uint32
}

func (i *InteractPacket) Opcode() byte {
	return 0x03
}

func (i *InteractPacket) Handle(buf *gbuf.GBuf) {
	buf.WriteUint16(i.ObjId)
	buf.WriteUint32(i.X)
	buf.WriteUint32(i.Y)
}
