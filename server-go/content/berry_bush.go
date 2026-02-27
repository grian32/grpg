package content

import (
	"server/scripts"
	"server/shared"
)

func init() {
	scripts.OnObjInteract(scripts.BERRY_BUSH, func(ctx *scripts.ObjInteractCtx) {
		if ctx.GetObjState() == 0 {
			ctx.SetObjState(1)
			ctx.PlayerInvAdd(scripts.BERRIES)
			ctx.PlayerAddXp(shared.Foraging, 100)

			ctx.AddTimer(100, func() {
				ctx.SetObjState(0)
			})
		}
	})
}
