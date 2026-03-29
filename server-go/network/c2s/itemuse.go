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
	if invIdx > 23 {
		log.Printf("warn: invalid idx passed for itemuse packet: %d", invIdx)
		return
	}
	item := player.Inventory.Items[invIdx]
	if item.ItemId == 0 {
		return
	}
	script, exists := scriptManager.ItemUseScripts[constants.ItemConstant(item.ItemId)]
	if !exists {
		return
	}
	script(scripts.NewItemUseCtx(game, player))
}
