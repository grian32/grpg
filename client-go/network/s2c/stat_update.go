package s2c

import (
	"client/shared"
	"grpg/data-go/gbuf"
	"log"
)

type StatUpdate struct {
}

func (s *StatUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	newHealth, err := buf.ReadByte()
	if err != nil {
		log.Printf("errored reading byte in stat update: %v", err)
		return
	}
	game.Player.Health = newHealth
}
