package c2s

import "grpg/data-go/gbuf"

type TalkPacket struct {
	NpcId uint16
	Uid   uint32
}

func (t *TalkPacket) Opcode() byte {
	return 0x04
}

func (t *TalkPacket) Handle(buf *gbuf.GBuf) {
	buf.WriteUint16(t.NpcId)
	buf.WriteUint32(t.Uid)
}
