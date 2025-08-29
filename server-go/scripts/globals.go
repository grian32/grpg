package scripts

import (
	"grpg/data-go/grpgnpc"
	"grpgscript/object"
	"log"
	"server/shared"
	"server/util"
)

func AddGlobals(env *object.Environment, game *shared.Game, npcs map[uint16]*grpgnpc.Npc) {
	env.Set("spawnNpc", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 3 {
				log.Printf("warn: script tried to call spawnNpc with less or more than 3 args")
			}

			npcId, ok := args[0].(*object.Integer)
			if !ok {
				log.Printf("warn: script tried to call spawnNpc with non integer first arg")
			}

			x, ok := args[1].(*object.Integer)
			if !ok {
				log.Printf("warn: script tried to call spawnNpc with non integer second arg")
			}

			y, ok := args[2].(*object.Integer)
			if !ok {
				log.Printf("warn: script tried to call spawnNpc with non integer third arg")
			}

			npcData, ok := npcs[uint16(npcId.Value)]
			if !ok {
				log.Printf("warn: script tried to call spawnNpc with invalid npc id %d", npcId.Value)
			}

			pos := util.Vector2I{X: uint32(x.Value), Y: uint32(y.Value)}

			game.TrackedNpcs[pos] = &shared.GameNpc{
				Pos:     pos,
				NpcData: npcData,
			}

			return nil
		},
	})
}
