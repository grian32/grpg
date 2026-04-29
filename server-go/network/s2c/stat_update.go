package s2c

import (
	"grpg/data-go/gbuf"
	"server/shared"
)

type StatUpdate struct {
	Player *shared.Player
}

func (s *StatUpdate) Opcode() byte {
	return 0x0C
}

func (s *StatUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	buf.WriteByte(s.Player.Health)
}
