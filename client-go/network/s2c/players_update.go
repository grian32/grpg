package s2c

import (
	"client/shared"
	"fmt"
	"grpg/data-go/gbuf"
)

type PlayersUpdate struct{}

func (p PlayersUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	var lst []shared.Player

	for range 2 {
		_, _ = buf.ReadInt32()
		name, _ := buf.ReadBytes(8)
		x, _ := buf.ReadInt32()
		y, _ := buf.ReadInt32()

		if string(name) == game.Player.Name {
			fmt.Println(x, y)
			game.Player.Move(x, y, game)
		} else {
			lst = append(lst, shared.Player{
				X:      x,
				Y:      y,
				RealX:  (x % 16) * 64,
				RealY:  (y % 16) * 64,
				ChunkX: x / 16,
				ChunkY: y / 16,
				Name:   string(name),
			})
		}
	}
	game.OtherPlayers = lst
}
