package c2s

import (
	"cmp"
	"fmt"
	"grpg/data-go/gbuf"
	"grpgscript/evaluator"
	"grpgscript/object"
	"log"
	"server/scripts"
	"server/shared"
	"server/util"
)

type Talk struct{}

func (t *Talk) Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player, scriptManager *scripts.ScriptManager) {
	npcId, err1 := buf.ReadUint16()
	x, err2 := buf.ReadUint32()
	y, err3 := buf.ReadUint32()

	if err := cmp.Or(err1, err2, err3); err != nil {
		log.Printf("failed reading npc in talk packet")
		return
	}

	_ = util.Vector2I{X: x, Y: y}
	script := scriptManager.NpcTalkScripts[npcId]
	env := object.NewEnclosedEnvinronment(scriptManager.Env)
	addTalkBuiltins(env, player)

	evaluator.Eval(script, env)
	fmt.Println(player.DialogueQueue)
}

func addTalkBuiltins(env *object.Environment, player *shared.Player) {
	// i could probably make this cleaner by making it generic but it's only 2 functions
	env.Set("talkPlayer", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				log.Printf("warn: script tried to call talkPlayer with less or more than 1 arg\n")
				return nil
			}
			talk, ok := args[0].(*object.String)
			if !ok {
				log.Printf("warn: script tried to call talkPlayer with non string argument\n")
				return nil
			}

			player.DialogueQueue.Dialogues = append(player.DialogueQueue.Dialogues, shared.Dialogue{
				Type:    shared.PLAYER,
				Content: talk.Value,
			})

			return nil
		},
	})
	env.Set("talkNpc", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				log.Printf("warn: script tried to call talkNpc with less or more than 1 arg\n")
				return nil
			}
			talk, ok := args[0].(*object.String)
			if !ok {
				log.Printf("warn: script tried to call talkNpc with non string argument\n")
				return nil
			}

			player.DialogueQueue.Dialogues = append(player.DialogueQueue.Dialogues, shared.Dialogue{
				Type:    shared.NPC,
				Content: talk.Value,
			})

			return nil
		},
	})
	env.Set("clearDialogueQueue", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 0 {
				log.Printf("warn: script tried to call clearDialogueQueue with non zero args\n")
			}

			player.DialogueQueue.Clear()
			return nil
		},
	})
	env.Set("startDialogue", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 0 {
				log.Printf("warn: script tried to call startDialogue with non zero args\n")
			}

			// TODO: send first talk packet

			return nil
		},
	})
}
