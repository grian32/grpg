package scripts

type ObjInteractFunc func(ctx *ObjInteractCtx)
type NpcTalkFunc func(ctx *NpcTalkCtx)

type PendingObjInteract struct {
	id ObjConstant
	fn ObjInteractFunc
}

type pendingNpcTalk struct {
	id NpcConstant
	fn NpcTalkFunc
}

type pendingNpcSpawn struct {
	npcId NpcConstant
	x     uint32
	y     uint32
}

var pendingObjInteracts []PendingObjInteract

var pendingNpcTalks []pendingNpcTalk

var pendingNpcSpawns []pendingNpcSpawn

func OnObjInteract(objId ObjConstant, fnc ObjInteractFunc) {
	pendingObjInteracts = append(pendingObjInteracts, PendingObjInteract{
		id: objId,
		fn: fnc,
	})
}

func OnTalkNpc(npcId NpcConstant, fnc NpcTalkFunc) {
	pendingNpcTalks = append(pendingNpcTalks, pendingNpcTalk{
		id: npcId,
		fn: fnc,
	})
}

func SpawnNpc(npcId NpcConstant, x uint32, y uint32) {
	pendingNpcSpawns = append(pendingNpcSpawns, pendingNpcSpawn{
		npcId: npcId,
		x:     x,
		y:     y,
	})
}
