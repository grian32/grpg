package s2c

import (
	"client/shared"
	"cmp"
	"fmt"
	"grpg/data-go/gbuf"
)

type PlayersUpdate struct{}

func (p PlayersUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	playersLen, err := buf.ReadUint16()

	if err != nil {
		fmt.Printf("Failed to read players update\n")
		return
	}

	sentNames := map[string]struct{}{}

	for _ = range playersLen {
		name, err1 := buf.ReadString()
		x, err2 := buf.ReadUint32()
		y, err3 := buf.ReadUint32()
		facing, err4 := buf.ReadByte()

		if err := cmp.Or(err1, err2, err3, err4); err != nil {
			fmt.Printf("Failed to read player struct %v\n", err)
			return
		}

		newX := int32(x)
		newY := int32(y)

		if name == game.Player.Name {
			game.Player.Move(newX, newY, shared.Direction(facing))
		} else {
			// TODO: could store in something, to avoid realloc
			sentNames[name] = struct{}{}
			player, exists := game.OtherPlayers[name]
			if exists {
				player.Move(newX, newY, shared.Direction(facing))
			} else {
				game.OtherPlayers[name] = shared.NewRemotePlayer(newX, newY, shared.Direction(facing), name)
			}
		}
	}

	for existingName, _ := range game.OtherPlayers {
		_, nameExists := sentNames[existingName]

		if !nameExists {
			delete(game.OtherPlayers, existingName)
		}
	}
}
