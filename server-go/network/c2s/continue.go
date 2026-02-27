package c2s

import (
	"grpg/data-go/gbuf"
	"server/network"
	"server/network/s2c"
	"server/scripts"
	"server/shared"
)

type Continue struct {
}

func (c *Continue) Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player, scriptManager *scripts.ScriptManager) {
	SendDialoguePacket(player, game, player.DialogueQueue.ActiveNpcId)
}

// TODO: figure out some better way of doing this? just ripped from npcctx atm, maybe refactor out into dq? nfc.
func SendDialoguePacket(player *shared.Player, game *shared.Game, npcId uint16) {
	if player.DialogueQueue.Index >= player.DialogueQueue.MaxIndex {
		network.SendPacket(player.Conn, &s2c.Talkbox{
			Type: s2c.CLEAR,
			Msg:  "",
		}, game)
		return
	}

	pktType := dqTypeToPacketType(player.DialogueQueue.Dialogues[player.DialogueQueue.Index].Type)

	network.SendPacket(player.Conn, &s2c.Talkbox{
		Type:  pktType,
		NpcId: npcId,
		Msg:   player.DialogueQueue.Dialogues[player.DialogueQueue.Index].Content,
	}, game)
	player.DialogueQueue.Index++
}

func dqTypeToPacketType(t shared.DialogueType) s2c.TalkboxType {
	if t == shared.NPC {
		return s2c.NPC
	}

	return s2c.PLAYER
}
