package s2c

import (
	"client/constants"
	"client/shared"
	"fmt"
	"grpg/data-go/gbuf"
)

type PlayerVarFull struct {}

func (p *PlayerVarFull) Handle(buf *gbuf.GBuf, game *shared.Game) {
	len, err := buf.ReadUint32()
	if err != nil {
		fmt.Printf("warn: error reading playervarfull: %v\n", err)
		return
	}
	for i := range len {
		value, err := buf.ReadUint16()
		if err != nil {
			fmt.Printf("warn: error reading playervarfull: %v\n", err)
			return
		}
		game.PlayerVars[constants.PlayerVarId(i)] = value
	}
	fmt.Printf("pv(full): %v", game.PlayerVars)
}
