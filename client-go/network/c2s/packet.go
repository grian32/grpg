package c2s

import "grpg/data-go/gbuf"

type Packet interface {
	Opcode() byte
	Handle(buf *gbuf.GBuf)
}
