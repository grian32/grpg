package scripts

import (
	"server/network"
	"server/network/s2c"
	"server/shared"
	"server/util"
)

type TimerFunc func()

type ObjInteractCtx struct {
	game   *shared.Game
	player *shared.Player
	objPos util.Vector2I
}

func (o *ObjInteractCtx) GetObjState() uint8 {
	return o.game.TrackedObjs[o.objPos].State
}

func (o *ObjInteractCtx) SetObjState(new uint8) {
	trackedObj := o.game.TrackedObjs[o.objPos]
	trackedObj.State = new

	network.UpdatePlayersByChunk(trackedObj.ChunkPos, o.game, &s2c.ObjUpdate{
		ChunkPos: trackedObj.ChunkPos,
		Rebuild:  false,
	})
}

func (o *ObjInteractCtx) PlayerInvAdd(itemId uint16) {
	o.player.Inventory.AddItem(itemId)
	network.SendPacket(o.player.Conn, &s2c.InventoryUpdate{
		Player: o.player,
	}, o.game)
}

func (o *ObjInteractCtx) PlayerAddXp(skill shared.Skill, xpAmount uint32) {
	o.player.AddXp(skill, xpAmount)
	network.SendPacket(o.player.Conn, &s2c.SkillUpdate{
		SkillIds: []shared.Skill{skill},
		Player:   o.player,
	}, o.game)
}

func (o *ObjInteractCtx) AddTimer(ticks uint64, fn TimerFunc) {
	// TODO: need to update timed scripts for this to work..
}

type NpcTalkContext struct {
	player *shared.Player
	game   *shared.Game
	npcId  uint16
}

func (n *NpcTalkContext) TalkPlayer(msg string) {
	// TODO: maybe make this append a function on dq
	n.player.DialogueQueue.Dialogues = append(n.player.DialogueQueue.Dialogues, shared.Dialogue{
		Type:    shared.PLAYER,
		Content: msg,
	})
	n.player.DialogueQueue.MaxIndex++
}

func (n *NpcTalkContext) TalkNpc(msg string) {
	// TODO: maybe make this append a function on dq
	n.player.DialogueQueue.Dialogues = append(n.player.DialogueQueue.Dialogues, shared.Dialogue{
		Type:    shared.NPC,
		Content: msg,
	})
	n.player.DialogueQueue.MaxIndex++
	n.player.DialogueQueue.ActiveNpcId = n.npcId
}

func (n *NpcTalkContext) ClearDialogueQueue() {
	n.player.DialogueQueue.Clear()
	n.sendDialoguePacket()
}

func (n *NpcTalkContext) StartDialogue() {
	n.sendDialoguePacket()
}

func (n *NpcTalkContext) sendDialoguePacket() {
	if n.player.DialogueQueue.Index >= n.player.DialogueQueue.MaxIndex {
		network.SendPacket(n.player.Conn, &s2c.Talkbox{
			Type: s2c.CLEAR,
			Msg:  "",
		}, n.game)
		return
	}

	pktType := dqTypeToPacketType(n.player.DialogueQueue.Dialogues[n.player.DialogueQueue.Index].Type)

	network.SendPacket(n.player.Conn, &s2c.Talkbox{
		Type:  pktType,
		NpcId: n.npcId,
		Msg:   n.player.DialogueQueue.Dialogues[n.player.DialogueQueue.Index].Content,
	}, n.game)
	n.player.DialogueQueue.Index++
}

func dqTypeToPacketType(t shared.DialogueType) s2c.TalkboxType {
	if t == shared.NPC {
		return s2c.NPC
	}

	return s2c.PLAYER
}
