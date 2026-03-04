package c2s

import "grpg/data-go/gbuf"

type Command struct {
	Msg string
}

func (c *Command) Opcode() byte {
	return 0x06
}

func (c *Command) Handle(buf *gbuf.GBuf) {
	buf.WriteUint16(4 + uint16(len(c.Msg))) // writeString = 4x from uint32 len + string
	buf.WriteString(c.Msg)
}
