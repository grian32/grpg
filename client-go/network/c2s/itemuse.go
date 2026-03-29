package c2s

import "grpg/data-go/gbuf"

type ItemUse struct {
	InvIdx uint8
}

func (i *ItemUse) Opcode() byte {
	return 0x08
}

func (i *ItemUse) Handle(buf *gbuf.GBuf) {
	buf.WriteByte(i.InvIdx)
}
