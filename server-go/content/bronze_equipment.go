package content

import (
	"server/constants"
	"server/scripts"
)

func equipScript(itemId constants.ItemConstant, slot int) {
	scripts.OnItemUse(itemId, func(ctx *scripts.ItemUseCtx) {
		ctx.EquipItem(ctx.InventoryIndex(), slot)
	})
}

func init() {
	equipScript(constants.BRONZE_HELM, 0)
	equipScript(constants.BRONZE_CHESTPLATE, 1)
	equipScript(constants.BRONZE_LEGS, 2)
	equipScript(constants.BRONZE_DAGGER, 3)
	equipScript(constants.BRONZE_RING, 4)
}
