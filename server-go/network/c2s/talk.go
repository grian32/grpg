package c2s

import (
	"cmp"
	"grpg/data-go/gbuf"
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
	_ = scriptManager.NpcTalkScripts[npcId]
	env := object.NewEnclosedEnvinronment(scriptManager.Env)
	addTalkBuiltins(env)

	//evaluator.Eval(script, env)
}

func addTalkBuiltins(env *object.Environment) {

}
