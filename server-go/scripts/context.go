package scripts

import (
	"server/constants"
	"server/network"
	"server/network/s2c"
	"server/shared"
	"server/util"
	"strconv"
)

type TimerFunc func()

// GenericCtx possibly bad name but this consists of functions that should be on every context
type GenericCtx struct {
	game   *shared.Game
	player *shared.Player
}

func (g *GenericCtx) InventoryAdd(itemId constants.ItemConstant) {
	err := g.player.Inventory.AddItem(g.game.Items[itemId])
	if err != nil {
		return
	}
	network.SendPacket(g.player.Conn, &s2c.InventoryUpdate{
		Player: g.player,
	}, g.game)
}

func (g *GenericCtx) AddXp(skill shared.Skill, xpAmount uint32) {
	g.player.AddXp(skill, xpAmount)
	network.SendPacket(g.player.Conn, &s2c.SkillUpdate{
		SkillIds: []shared.Skill{skill},
		Player:   g.player,
	}, g.game)
}

func (g *GenericCtx) AddTimer(ticks uint32, fn TimerFunc) {
	endTick := g.game.CurrentTick + ticks
	_, ok := g.game.TimedScripts[endTick]
	if !ok {
		g.game.TimedScripts[endTick] = []func(){fn}
	} else {
		g.game.TimedScripts[endTick] = append(g.game.TimedScripts[endTick], fn)
	}
}

func (g *GenericCtx) GetPlayerVar(varId constants.PlayerVarId) uint16 {
	return g.player.PlayerVars[varId]
}

func (g *GenericCtx) SetPlayerVar(varId constants.PlayerVarId, newValue uint16) {
	g.player.PlayerVars[varId] = newValue
	network.SendPacket(g.player.Conn, &s2c.PlayerVarIndiv{
		VarId:    uint16(varId),
		VarValue: newValue,
	}, g.game)
}

func (g *GenericCtx) Player() *shared.Player {
	return g.player
}

func (g *GenericCtx) EquipItem(invIdx byte, slot int) {
	if g.player.Equipment.Items[slot].ItemId != 0 {
		return
	}
	item := g.player.Inventory.Items[invIdx]
	if item.ItemId == 0 {
		return
	}
	g.player.Inventory.Items[invIdx] = shared.InventoryItem{Dirty: true}
	g.player.Equipment.Items[slot] = shared.EquipmentItem{
		ItemId: item.ItemId,
		Dirty:  true,
	}
	network.SendPacket(g.player.Conn, &s2c.InventoryUpdate{
		Player: g.player,
	}, g.game)
}

func (g *GenericCtx) UnequipItem(slot int) {
	if g.player.Equipment.Items[slot].ItemId == 0 {
		return
	}
	item := g.player.Equipment.Items[slot]
	err := g.player.Inventory.AddItem(g.game.Items[constants.ItemConstant(item.ItemId)])
	if err != nil {
		return
	}
	g.player.Equipment.Items[slot] = shared.EquipmentItem{Dirty: true}
	network.SendPacket(g.player.Conn, &s2c.InventoryUpdate{
		Player: g.player,
	}, g.game)
}

type ObjInteractCtx struct {
	game   *shared.Game
	player *shared.Player
	objPos util.Vector2I
	GenericCtx
}

func NewObjInteractCtx(game *shared.Game, player *shared.Player, objPos util.Vector2I) *ObjInteractCtx {
	return &ObjInteractCtx{
		game:   game,
		player: player,
		objPos: objPos,
		GenericCtx: GenericCtx{
			game:   game,
			player: player,
		},
	}
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

type NpcTalkCtx struct {
	player *shared.Player
	game   *shared.Game
	npcId  constants.NpcConstant
	GenericCtx
}

func NewNpcTalkCtx(player *shared.Player, game *shared.Game, npcId constants.NpcConstant) *NpcTalkCtx {
	return &NpcTalkCtx{
		player: player,
		game:   game,
		npcId:  npcId,
		GenericCtx: GenericCtx{
			game:   game,
			player: player,
		},
	}
}

func (n *NpcTalkCtx) TalkPlayer(msg string) {
	// TODO: maybe make this append a function on dq
	n.player.DialogueQueue.Dialogues = append(n.player.DialogueQueue.Dialogues, shared.Dialogue{
		Type:    shared.PLAYER,
		Content: msg,
	})
	n.player.DialogueQueue.MaxIndex++
}

func (n *NpcTalkCtx) TalkNpc(msg string) {
	// TODO: maybe make this append a function on dq
	n.player.DialogueQueue.Dialogues = append(n.player.DialogueQueue.Dialogues, shared.Dialogue{
		Type:    shared.NPC,
		Content: msg,
	})
	n.player.DialogueQueue.MaxIndex++
	n.player.DialogueQueue.ActiveNpcId = uint16(n.npcId)
}

func (n *NpcTalkCtx) ClearDialogueQueue() {
	n.player.DialogueQueue.Clear()
	n.sendDialoguePacket()
}

func (n *NpcTalkCtx) StartDialogue() {
	n.sendDialoguePacket()
}

func (n *NpcTalkCtx) sendDialoguePacket() {
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
		NpcId: uint16(n.npcId),
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

type CommandCtx struct {
	args       []string
	currArgIdx uint
	game       *shared.Game
	GenericCtx
}

func NewCommandCtx(args []string, game *shared.Game, player *shared.Player) *CommandCtx {
	return &CommandCtx{
		args:       args,
		currArgIdx: 0,
		game:       game,
		GenericCtx: GenericCtx{
			game:   game,
			player: player,
		},
	}
}

// Args basically getter as i dont want users modifiying these, & i wanna add automatic parsing of some sort l8r
func (c *CommandCtx) Args() []string {
	return c.args
}
func (c *CommandCtx) Game() *shared.Game { return c.game }

func (c *CommandCtx) GetIntArg() (int64, error) {
	arg := c.args[c.currArgIdx]
	p, err := strconv.ParseInt(arg, 0, 64)
	if err != nil {
		return -1, err
	}
	c.currArgIdx++
	return p, nil
}

type ItemUseCtx struct {
	GenericCtx
	invIdx byte
	equip  bool
}

func NewItemUseCtx(game *shared.Game, player *shared.Player, invIdx byte, equip bool) *ItemUseCtx {
	i := &ItemUseCtx{
		GenericCtx: GenericCtx{
			game:   game,
			player: player,
		},
		equip: equip,
	}
	if i.equip {
		i.invIdx = invIdx - 24
	} else {
		i.invIdx = invIdx
	}

	return i
}

func (c *ItemUseCtx) InventoryIndex() byte {
	return c.invIdx
}

func (c *ItemUseCtx) IsEquipmentSlot() bool {
	return c.equip
}
