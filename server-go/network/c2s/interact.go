package c2s

import (
	"cmp"
	"fmt"
	"grpg/data-go/gbuf"
	"log"
	"server/scripts"
	"server/shared"
	"server/util"
)

type Interact struct{}

func (i *Interact) Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player, scriptManager *scripts.ScriptManager) {
	objId, err1 := buf.ReadUint16()
	x, err2 := buf.ReadUint32()
	y, err3 := buf.ReadUint32()

	objPos := util.Vector2I{X: x, Y: y}

	playerFacingCooord := player.GetFacingCoord()
	if player.GetFacingCoord() != objPos {
		fmt.Printf("warn: player %s @ facing %d, %d, %s tried to interact with obj @ %d, %d that he isn't facing\n",
			player.Name,
			playerFacingCooord.X, playerFacingCooord.Y,
			shared.DirectionString(player.Facing),
			x, y)
		return
	}

	if _, ok := game.Objs[objPos]; !ok {
		fmt.Printf("warn: player %s tried to interact with obj that doesn't exist %d, %d\n", player.Name, x, y)
		return
	}

	if err := cmp.Or(err1, err2, err3); err != nil {
		log.Printf("failed to read interact packet: %v\n", err)
	}

	script := scriptManager.InteractScripts[scripts.ObjConstant(objId)]
	script(scripts.NewObjInteractCtx(game, player, objPos))
}
