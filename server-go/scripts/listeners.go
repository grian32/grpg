package scripts

import "server/constants"

type ObjInteractFunc func(ctx *ObjInteractCtx)
type NpcTalkFunc func(ctx *NpcTalkCtx)
type CommandFunc func(ctx *CommandCtx)
type ItemUseFunc func(ctx *ItemUseCtx)

type PendingObjInteract struct {
	id constants.ObjConstant
	fn ObjInteractFunc
}

type pendingNpcTalk struct {
	id constants.NpcConstant
	fn NpcTalkFunc
}

type pendingNpcSpawn struct {
	npcId       constants.NpcConstant
	x           uint32
	y           uint32
	wanderRange uint8
}

type pendingCmd struct {
	name string
	fn   CommandFunc
}

type pendingItemUse struct {
	itemId constants.ItemConstant
	fn     ItemUseFunc
}

var pendingObjInteracts []PendingObjInteract

var pendingNpcTalks []pendingNpcTalk

var pendingNpcSpawns []pendingNpcSpawn

var pendingCmds []pendingCmd

var pendingItemUses []pendingItemUse

func OnObjInteract(objId constants.ObjConstant, fnc ObjInteractFunc) {
	pendingObjInteracts = append(pendingObjInteracts, PendingObjInteract{
		id: objId,
		fn: fnc,
	})
}

func OnTalkNpc(npcId constants.NpcConstant, fnc NpcTalkFunc) {
	pendingNpcTalks = append(pendingNpcTalks, pendingNpcTalk{
		id: npcId,
		fn: fnc,
	})
}

func SpawnNpc(npcId constants.NpcConstant, x uint32, y uint32, wanderRange uint8) {
	pendingNpcSpawns = append(pendingNpcSpawns, pendingNpcSpawn{
		npcId:       npcId,
		x:           x,
		y:           y,
		wanderRange: wanderRange,
	})
}

func OnCommand(name string, fnc CommandFunc) {
	pendingCmds = append(pendingCmds, pendingCmd{
		name: name,
		fn:   fnc,
	})
}

func OnItemUse(itemId constants.ItemConstant, fnc ItemUseFunc) {
	pendingItemUses = append(pendingItemUses, pendingItemUse{
		itemId: itemId,
		fn:     fnc,
	})
}
