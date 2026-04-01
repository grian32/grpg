package c2s

import (
	"grpg/data-go/gbuf"
	"log"
	"server/constants"
	"server/scripts"
	"server/shared"
)

type ItemUse struct{}

func (i *ItemUse) Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player, scriptManager *scripts.ScriptManager) {
	invIdx, err := buf.ReadByte()
	if err != nil {
		log.Printf("warn: failed to read item id in item use packet: %v", err)
		return
	}
	if invIdx > 28 {
		log.Printf("warn: invalid idx passed for itemuse packet: %d", invIdx)
		return
	}
	equip := invIdx > 23
	var itemId uint16
	if !equip {
		item := player.Inventory.Items[invIdx]
		if item.ItemId == 0 {
			return
		}
		itemId = item.ItemId
	} else {
		// looks similar but equipment item also uses a diff struct so not much can do
		item := player.Equipment.Items[invIdx-24]
		if item.ItemId == 0 {
			return
		}
		itemId = item.ItemId
	}
	script, exists := scriptManager.ItemUseScripts[constants.ItemConstant(itemId)]
	if !exists {
		return
	}
	script(scripts.NewItemUseCtx(game, player, invIdx, equip))
}
