package c2s

import (
	"cmp"
	"grpg/data-go/gbuf"
	"grpgscript/evaluator"
	"grpgscript/object"
	"log"
	"server/network"
	"server/network/s2c"
	"server/scripts"
	"server/shared"
	"server/util"
)

type Interact struct{}

func (i *Interact) Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player, scriptManager *scripts.ScriptManager) {
	objId, err1 := buf.ReadUint16()
	x, err2 := buf.ReadUint32()
	y, err3 := buf.ReadUint32()

	if err := cmp.Or(err1, err2, err3); err != nil {
		log.Printf("failed to read interact packet: %v\n", err)
	}

	objPos := util.Vector2I{X: x, Y: y}
	script := scriptManager.InteractScripts[objId]
	env := object.NewEnclosedEnvinronment(scriptManager.Env)
	addInteractBuiltins(env, game, player, objPos, scriptManager)

	evaluator.Eval(script, env)
}

func addInteractBuiltins(env *object.Environment, game *shared.Game, player *shared.Player, objPos util.Vector2I, scriptManager *scripts.ScriptManager) {
	env.Set("getObjState", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			return &object.Integer{Value: int64(game.TrackedObjs[objPos].State)}
		},
	})
	env.Set("setObjState", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				log.Printf("warn: script tries to call setObjState in onInteract ctx with non-1 arguments")
				return nil
			}

			newState, ok := args[0].(*object.Integer)
			if !ok {
				log.Printf("warn: script tries to call setObjState in onInteract ctx without int arg")
				return nil
			}

			trackedObj := game.TrackedObjs[objPos]
			trackedObj.State = byte(newState.Value)

			network.UpdatePlayersByChunk(trackedObj.ChunkPos, game, &s2c.ObjUpdate{ChunkPos: trackedObj.ChunkPos, Rebuild: false})

			return nil
		},
	})
	env.Set("playerInvAdd", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				log.Printf("warn: script tries to call playerInvAdd in onInteract ctx with non-1 arguments")
				return nil
			}

			itemId, ok := args[0].(*object.Integer)
			if !ok {
				log.Printf("warn: scriped tries to call playerInvAdd in onInteract ctx without int arg")
				return nil
			}

			player.Inventory.AddItem(uint16(itemId.Value))

			network.SendPacket(player.Conn, &s2c.InventoryUpdate{Player: player}, game)

			return nil
		},
	})
	env.Set("timer", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 2 {
				log.Printf("warn: script tries to call timer in onInteract ctx with non-2 arguments")
				return nil
			}

			tickCount, ok := args[0].(*object.Integer)
			if !ok {
				log.Printf("warn: script tries to call timer in onInteract ctx without int arg")
				return nil
			}
			fn, ok := args[1].(*object.Function)
			if !ok {
				log.Printf("warn: script tries to call timer in onInteract ctx without function arg")
				return nil
			}

			scriptManager.AddTimedScript(game.CurrentTick+uint32(tickCount.Value),
				scripts.TimedScript{
					Script: fn.Body,
					Env:    env,
				})

			return nil
		},
	})
}
