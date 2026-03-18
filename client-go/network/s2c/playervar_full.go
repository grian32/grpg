package s2c

import (
	"client/constants"
	"client/shared"
	"grpg/data-go/gbuf"
	"log"
)

type PlayerVarFull struct {}

func (p *PlayerVarFull) Handle(buf *gbuf.GBuf, game *shared.Game) {
	len, err := buf.ReadUint32()
	if err != nil {
		log.Printf("warn: error reading playervarfull: %v\n", err)
		return
	}
	for i := range len {
		value, err := buf.ReadUint16()
		if err != nil {
			log.Printf("warn: error reading playervarfull: %v\n", err)
			return
		}
		game.PlayerVars[constants.PlayerVarId(i)] = value
		handler, ok := game.PlayerVarHandlers[constants.PlayerVarId(i)]
		if !ok {
			log.Printf("warn: no handler for playervar %d\n", i)
			continue
		}
		log.Printf("pvfull: %d=%d", i, value)
		handler(game, value)
	}
}
