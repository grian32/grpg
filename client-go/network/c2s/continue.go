package c2s

import "grpg/data-go/gbuf"

type Continue struct {
}

func (c *Continue) Opcode() byte {
	return 0x05
}

func (c *Continue) Handle(buf *gbuf.GBuf) {
}
