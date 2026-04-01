package content

import (
	"server/constants"
	"server/scripts"
	"server/shared"
)

func init() {
	scripts.OnObjInteract(constants.BERRY_BUSH, func(ctx *scripts.ObjInteractCtx) {
		if ctx.GetObjState() == 0 {
			ctx.SetObjState(1)
			ctx.InventoryAdd(constants.BERRIES)
			ctx.AddXp(shared.Foraging, 100)

			ctx.AddTimer(100, func() {
				ctx.SetObjState(0)
			})
		}
	})
}
