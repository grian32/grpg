package scripts

import (
	"grpgscript/object"
	"log"
)

func AddListeners(env *object.Environment, scriptManager *ScriptManager) {
	// TODO: figure out some better way to do some of this, atleast some of the validation
	env.Set("onInteract", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			id, ok := args[0].(*object.Integer)
			if !ok {
				log.Fatal("script tried to call onInteract with non integer argument")
			}

			fn, ok := args[1].(*object.Function)
			if !ok {
				log.Fatal("script tried to call onInteract with non function argument")
			}

			scriptManager.InteractScripts[uint16(id.Value)] = fn.Body

			return nil
		},
	})
}
