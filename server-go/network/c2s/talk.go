package c2s

import (
	"grpg/data-go/gbuf"
	"log"
	"server/scripts"
	"server/constants"
	"server/shared"
)

type Talk struct{}

func (t *Talk) Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player, scriptManager *scripts.ScriptManager) {
	npcId, err := buf.ReadUint16()
	if err != nil {
		log.Printf("warn: failed to read uint32 in talk packet\n")
		return
	}
	uid, err := buf.ReadUint32()
	if err != nil {
		log.Printf("warn: failed to read uint32 in talk packet\n")
		return
	}

	npc, ok := game.TrackedNpcs[uid]
	if !ok {
		log.Printf("warn: player %s tried to talk with npc that doesn't exist with uid %d\n", player.Name, uid)
		return
	}

	npcPos := npc.Pos

	if player.GetFacingCoord() != npcPos {
		log.Printf("warn: player %s, facing [%d, %d] tried to talk with npc that he isn't facing, uid : %d\n", player.Name, player.GetFacingCoord().X, player.GetFacingCoord().Y, uid)
		return
	}

	script := scriptManager.NpcTalkScripts[constants.NpcConstant(npcId)]
	script(scripts.NewNpcTalkCtx(player, game, constants.NpcConstant(npcId)))
}
