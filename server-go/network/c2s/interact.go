package c2s

import (
	"cmp"
	"grpg/data-go/gbuf"
	"grpgscript/evaluator"
	"grpgscript/object"
	"log"
	"server/network"
	"server/network/s2c"
	"server/shared"
	"server/util"
)

type Interact struct{}

func (i *Interact) Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player) {
	objId, err1 := buf.ReadUint16()
	x, err2 := buf.ReadUint32()
	y, err3 := buf.ReadUint32()

	if err := cmp.Or(err1, err2, err3); err != nil {
		log.Printf("failed to read interact packet: %v\n", err)
	}

	objPos := util.Vector2I{X: x, Y: y}
	script := game.ScriptManager.InteractScripts[objId]
	env := object.NewEnclosedEnvinronment(game.ScriptManager.Env)
	addInteractBuiltins(env, game, player, objPos)

	evaluator.Eval(script, env)

	chunkPos := game.TrackedObjs[objPos].ChunkPos

	network.UpdatePlayersByChunk(chunkPos, game, &s2c.ObjUpdate{ChunkPos: chunkPos, Rebuild: false})
}

func addInteractBuiltins(env *object.Environment, game *shared.Game, player *shared.Player, objPos util.Vector2I) {
	env.Set("getObjState", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			return &object.Integer{Value: int64(game.TrackedObjs[objPos].State)}
		},
	})
	env.Set("setObjState", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			newState, ok := args[0].(*object.Integer)
			if !ok {
				log.Printf("warn: script tries to call setObjState in onInteract ctx without int arg")
				return nil
			}

			game.TrackedObjs[objPos].State = byte(newState.Value)

			return nil
		},
	})
	env.Set("playerInvAdd", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			itemId, ok := args[0].(*object.Integer)
			if !ok {
				log.Printf("warn: scriped tries to call playerInvAdd in onInteract ctx without int arg")
				return nil
			}

			firstEmptyIdx := -1

			for idx := range 24 {
				if player.Inventory[idx].ItemId == uint16(itemId.Value) {
					player.Inventory[idx].Count++
					player.Inventory[idx].Dirty = true
					network.SendPacket(player.Conn, &s2c.InventoryUpdate{Player: player}, game)
					return nil
				}

				if player.Inventory[idx].ItemId == 0 && firstEmptyIdx == -1 {
					firstEmptyIdx = idx
				}
			}

			// if it finds a pre existing stack then it returns early anyway so np
			if firstEmptyIdx != -1 {
				player.Inventory[firstEmptyIdx].ItemId = uint16(itemId.Value)
				player.Inventory[firstEmptyIdx].Count = 1
				player.Inventory[firstEmptyIdx].Dirty = true
			}

			network.SendPacket(player.Conn, &s2c.InventoryUpdate{Player: player}, game)

			return nil
		},
	})
}
