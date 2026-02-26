package scripts

import (
	"grpg/data-go/grpgnpc"
	"log"
	"server/shared"
	"server/util"
)

type ScriptManager struct {
	InteractScripts map[uint16]ObjInteractFunc
	NpcTalkScripts  map[uint16]NpcTalkFunc
	TimedScripts    map[uint32][]TimerFunc
}

func (s *ScriptManager) AddTimedScript(tick uint32, script TimerFunc) {
	_, ok := s.TimedScripts[tick]
	if !ok {
		s.TimedScripts[tick] = []TimerFunc{script}
	} else {
		s.TimedScripts[tick] = append(s.TimedScripts[tick], script)
	}
}

func NewScriptManager(game *shared.Game, npcs map[uint16]*grpgnpc.Npc) *ScriptManager {
	s := &ScriptManager{
		InteractScripts: make(map[uint16]ObjInteractFunc),
		NpcTalkScripts:  make(map[uint16]NpcTalkFunc),
		TimedScripts:    make(map[uint32][]TimerFunc),
	}

	for _, reg := range pendingObjInteracts {
		s.InteractScripts[reg.id] = reg.fn
	}

	for _, reg := range pendingNpcTalks {
		s.NpcTalkScripts[reg.id] = reg.fn
	}

	for _, reg := range pendingNpcSpawns {
		npcData, ok := npcs[reg.npcId]
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
