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

type Talk struct{}

func (t *Talk) Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player, scriptManager *scripts.ScriptManager) {
	npcId, err1 := buf.ReadUint16()
	x, err2 := buf.ReadUint32()
	y, err3 := buf.ReadUint32()

	npcPos := util.Vector2I{X: x, Y: y}

	if player.GetFacingCoord() != npcPos {
		fmt.Printf("warn: player %s tried to talk with npc that he isn't facing %d, %d", player.Name, x, y)
		return
	}

	if _, ok := game.TrackedNpcs[npcPos]; !ok {
		fmt.Printf("warn: player %s tried to talk with npc that doesn't exist %d, %d", player.Name, x, y)
		return
	}

	if err := cmp.Or(err1, err2, err3); err != nil {
		log.Printf("failed reading npc in talk packet")
		return
	}

	script := scriptManager.NpcTalkScripts[scripts.NpcConstant(npcId)]
	script(scripts.NewNpcTalkCtx(player, game, scripts.NpcConstant(npcId)))
}
