package c2s

import (
	"cmp"
	"fmt"
	"grpg/data-go/gbuf"
	"grpgscript/evaluator"
	"grpgscript/object"
	"log"
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

	script := game.ScriptManager.InteractScripts[objId]
	env := object.NewEnclosedEnvinronment(game.ScriptManager.Env)
	addInteractBuiltins(env, game, util.Vector2I{X: x, Y: y})

	evaluator.Eval(script, env)

	for pos, obj := range game.TrackedObjs {
		if pos.Y == 0 {
			fmt.Printf("%v: %d\n", pos, obj.State)
		}
	}
}

func addInteractBuiltins(env *object.Environment, game *shared.Game, objPos util.Vector2I) {
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

			game.TrackedObjs[objPos].State = uint16(newState.Value)

			return nil
		},
	})
	env.Set("playerInvAdd", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			// TODO
			return nil
		},
	})
}
