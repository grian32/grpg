package content

import (
	"log"
	"server/constants"
	"server/network"
	"server/network/s2c"
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
		ctx.InventoryAdd(constants.ItemConstant(id))
	})

	scripts.OnCommand("sethealth", func(ctx *scripts.CommandCtx) {
		newHp, err := ctx.GetIntArg()
		if err != nil {
			log.Printf("failed to parse int in sethealth cmd")
			return
		}

		ctx.Player().Health = uint8(newHp)

		network.SendPacket(ctx.Player().Conn, &s2c.StatUpdate{
			Player: ctx.Player(),
		}, ctx.Game())
	})
}
