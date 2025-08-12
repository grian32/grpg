package c2s

import (
	"cmp"
	"grpg/data-go/gbuf"
	"grpgscript/evaluator"
	"grpgscript/object"
	"log"
	"server/shared"
)

type Interact struct{}

func (i *Interact) Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player) {
	objId, err1 := buf.ReadUint16()
	_, err2 := buf.ReadUint32()
	_, err3 := buf.ReadUint32()

	if err := cmp.Or(err1, err2, err3); err != nil {
		log.Printf("failed to read interact packet: %v\n", err)
	}

	script := game.ScriptManager.InteractScripts[objId]
	env := object.NewEnvironment()

	evaluator.Eval(script, env)
}
