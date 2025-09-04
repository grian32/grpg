package scripts

import (
	"grpgscript/object"
	"log"
)

func AddListeners(env *object.Environment, scriptManager *ScriptManager) {
	// TODO: figure out some better way to do some of this, atleast some of the validation
	env.Set("onInteract", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 2 {
				log.Println("warn: script tried to call onInteract with non-2 arguments")
				return nil
			}

			id, ok := args[0].(*object.Integer)
			if !ok {
				log.Println("warn: script tried to call onInteract with non integer argument")
				return nil
			}

			fn, ok := args[1].(*object.Function)
			if !ok {
				log.Println("warn: script tried to call onInteract with non function argument")
				return nil
			}

			scriptManager.InteractScripts[uint16(id.Value)] = fn.Body

			return nil
		},
	})
	env.Set("onTalkNpc", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 2 {
				log.Printf("warn: script tried to call onTalkNpc with less or more than 2 arguments.")
				return nil
			}
			id, ok := args[0].(*object.Integer)
			if !ok {
				log.Printf("script tried to call onTalkNpc with non integer argument")
				return nil
			}

			fn, ok := args[1].(*object.Function)
			if !ok {
				log.Fatal("script tried to call onTalkNpc with non function argument")
				return nil
			}

			scriptManager.NpcTalkScripts[uint16(id.Value)] = fn.Body

			return nil
		},
	})
}
