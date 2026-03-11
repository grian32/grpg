package content

import (
	"log"
	"server/scripts"
	"strings"
)

func init() {
	scripts.OnCommand("log", func(ctx *scripts.CommandCtx) {
		log.Printf("cmd log: %s", strings.Join(ctx.Args(), " "))
	})

	// scripts.OnCommand("wander", func(ctx *scripts.CommandCtx) {
	// 	x, err1 := ctx.GetIntArg()
	// 	y, err2 := ctx.GetIntArg()
	// 	if err := cmp.Or(err1, err2); err != nil {
	// 		log.Printf("failed to parse ints in wander cmd")
	// 		return
	// 	}
	// 	g := ctx.Game()
	// 	npc, ok := g.TrackedNpcs[util.Vector2I{X: uint32(x), Y: uint32(y)}]
	// 	if !ok {
	// 		log.Printf("failed to find npc for wander cmd")
	// 		return
	// 	}
	// 	npc.Wander(g)
	// })
}
