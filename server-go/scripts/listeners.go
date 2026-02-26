package scripts

type ObjInteractFunc func(ctx *ObjInteractCtx)
type NpcTalkFunc func(ctx *NpcTalkContext)

type PendingObjInteract struct {
	id uint16
	fn ObjInteractFunc
}

type pendingNpcTalk struct {
	id uint16
	fn NpcTalkFunc
}

type pendingNpcSpawn struct {
	npcId uint16
	x     uint32
	y     uint32
}

var pendingObjInteracts []PendingObjInteract

var pendingNpcTalks []pendingNpcTalk

var pendingNpcSpawns []pendingNpcSpawn

func OnObjInteract(objId uint16, fnc ObjInteractFunc) {
	pendingObjInteracts = append(pendingObjInteracts, PendingObjInteract{
		id: objId,
		fn: fnc,
	})
}

func OnTalkNpc(npcId uint16, fnc NpcTalkFunc) {
	pendingNpcTalks = append(pendingNpcTalks, pendingNpcTalk{
		id: npcId,
		fn: fnc,
	})
}

func SpawnNpc(npcId uint16, x uint32, y uint32) {
	pendingNpcSpawns = append(pendingNpcSpawns, pendingNpcSpawn{
		npcId: npcId,
		x:     x,
		y:     y,
	})
}
