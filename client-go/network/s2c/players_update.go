package s2c

import (
	"client/shared"
	"cmp"
	"fmt"
	"grpg/data-go/gbuf"
)

type PlayersUpdate struct{}

func (p PlayersUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	var lst []shared.RemotePlayer

	playersLen, err := buf.ReadUint16()

	if err != nil {
		fmt.Printf("Failed to read players update\n")
		return
	}

	for _ = range playersLen {
		name, err1 := buf.ReadString()
		x, err2 := buf.ReadUint32()
		y, err3 := buf.ReadUint32()

		if err := cmp.Or(err1, err2, err3); err != nil {
			fmt.Printf("Failed to read player struct %v\n", err)
			return
		}

		newX := int32(x)
		newY := int32(y)

		if name == game.Player.Name {
			game.Player.Move(newX, newY)
		} else {
			lst = append(lst, shared.NewRemotePlayer(newX, newY, shared.DOWN, name, game))
		}
	}

	game.OtherPlayers = lst
}
