package s2c

import (
	"client/constants"
	"client/shared"
	"fmt"
	"grpg/data-go/gbuf"
)

type PlayerVarIndiv struct {}

func (p *PlayerVarIndiv) Handle(buf *gbuf.GBuf, game *shared.Game) {
	id, err := buf.ReadUint16()
	if err != nil {
		fmt.Printf("warn: error reading playervarindiv: %v\n", err)
		return
	}
	val, err := buf.ReadUint16()
	if err != nil {
		fmt.Printf("warn: error reading playervarindiv: %v\n", err)
		return
	}
	game.PlayerVars[constants.PlayerVarId(id)] = val
	handler, ok := game.PlayerVarHandlers[constants.PlayerVarId(id)]
	if !ok {
		fmt.Printf("warn: no handler for playervar %d\n", id)
		return
	}
	handler(game, val)
}
