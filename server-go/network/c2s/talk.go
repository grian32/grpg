package c2s

import (
	"cmp"
	"fmt"
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
	addTalkBuiltins(env, player, game)

	evaluator.Eval(script, env)
	fmt.Println(player.DialogueQueue)
}

func addTalkBuiltins(env *object.Environment, player *shared.Player, game *shared.Game) {
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
			player.DialogueQueue.MaxIndex++

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
			player.DialogueQueue.MaxIndex++

			return nil
		},
	})
	env.Set("clearDialogueQueue", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 0 {
				log.Printf("warn: script tried to call clearDialogueQueue with non zero args\n")
			}

			player.DialogueQueue.Clear()
			SendDialoguePacket(player, game)
			return nil
		},
	})
	env.Set("startDialogue", &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 0 {
				log.Printf("warn: script tried to call startDialogue with non zero args\n")
			}

			SendDialoguePacket(player, game)

			return nil
		},
	})
}

func SendDialoguePacket(player *shared.Player, game *shared.Game) {
	if player.DialogueQueue.Index >= player.DialogueQueue.MaxIndex {
		network.SendPacket(player.Conn, &s2c.Talkbox{
			Type: s2c.CLEAR,
			Msg:  "",
		}, game)
		return
	}

	pktType := dqTypeToPacketType(player.DialogueQueue.Dialogues[player.DialogueQueue.Index].Type)

	network.SendPacket(player.Conn, &s2c.Talkbox{
		Type: pktType,
		Msg:  player.DialogueQueue.Dialogues[player.DialogueQueue.Index].Content,
	}, game)
	player.DialogueQueue.Index++
}

func dqTypeToPacketType(t shared.DialogueType) s2c.TalkboxType {
	if t == shared.NPC {
		return s2c.NPC
	}

	return s2c.PLAYER
}
