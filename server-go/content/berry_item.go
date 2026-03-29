package content

import (
	"log"
	"server/constants"
	"server/scripts"
)

func init() {
	scripts.OnItemUse(constants.BERRIES, func(ctx *scripts.ItemUseCtx) {
		log.Println("used berries")
	})
}
