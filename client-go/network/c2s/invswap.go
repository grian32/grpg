package c2s

import "grpg/data-go/gbuf"

type InvSwap struct {
	From byte
	To   byte
}

func (s *InvSwap) Opcode() byte {
	return 0x07
}

func (s *InvSwap) Handle(buf *gbuf.GBuf) {
	buf.WriteByte(s.From)
	buf.WriteByte(s.To)
}
