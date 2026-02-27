package scripts

import (
	"grpg/data-go/grpgnpc"
	"log"
	"server/shared"
	"server/util"
)

type ScriptManager struct {
	InteractScripts map[ObjConstant]ObjInteractFunc
	NpcTalkScripts  map[NpcConstant]NpcTalkFunc
}

func NewScriptManager(game *shared.Game, npcs map[uint16]*grpgnpc.Npc) *ScriptManager {
	s := &ScriptManager{
		InteractScripts: make(map[ObjConstant]ObjInteractFunc),
		NpcTalkScripts:  make(map[NpcConstant]NpcTalkFunc),
	}

	for _, reg := range pendingObjInteracts {
		s.InteractScripts[reg.id] = reg.fn
	}

	for _, reg := range pendingNpcTalks {
		s.NpcTalkScripts[reg.id] = reg.fn
	}

	for _, reg := range pendingNpcSpawns {
		npcData, ok := npcs[uint16(reg.npcId)]
		if !ok {
			log.Printf("unknown npc %d for npcSpawn", reg.npcId)
			continue
		}

		pos := util.Vector2I{X: reg.x, Y: reg.y}
		chunkPos := util.Vector2I{X: pos.X / 16, Y: pos.Y / 16}

		game.TrackedNpcs[pos] = &shared.GameNpc{
			Pos:      pos,
			NpcData:  npcData,
			ChunkPos: chunkPos,
		}
	}

	pendingObjInteracts = nil
	pendingNpcTalks = nil
	pendingNpcSpawns = nil

	return s
}
