package c2s

import (
	"grpg/data-go/gbuf"
	"server/scripts"
	"server/shared"
)

type Continue struct {
}

func (c *Continue) Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player, scriptManager *scripts.ScriptManager) {
	SendDialoguePacket(player, game, player.DialogueQueue.ActiveNpcId)
}
