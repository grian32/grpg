package s2c

import "grpg/data-go/gbuf"

type Packet interface {
	Opcode() byte
	Handle(buf *gbuf.GBuf)
}
