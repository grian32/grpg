package content

import (
	"log"
	"server/constants"
	"server/scripts"
	"strings"
)

func init() {
	scripts.OnCommand("log", func(ctx *scripts.CommandCtx) {
		log.Printf("cmd log: %s", strings.Join(ctx.Args(), " "))
	})

	scripts.OnCommand("wander", func(ctx *scripts.CommandCtx) {
		uid, err := ctx.GetIntArg()
		if err != nil {
			log.Printf("failed to parse ints in wander cmd")
			return
		}
		g := ctx.Game()
		npc, ok := g.TrackedNpcs[uint32(uid)]
		if !ok {
			log.Printf("failed to find npc for wander cmd")
			return
		}
		npc.Wander(g)
	})

	scripts.OnCommand("logequipment", func(ctx *scripts.CommandCtx) {
		log.Printf("equipment for %s: %v", ctx.Player().Name, ctx.Player().Equipment)
	})

	scripts.OnCommand("additem", func(ctx *scripts.CommandCtx) {
		id, err := ctx.GetIntArg()
		if err != nil {
			log.Printf("failed to parse int in additem cmd")
			return
		}
		ctx.PlayerInvAdd(constants.ItemConstant(id))
	})
}
