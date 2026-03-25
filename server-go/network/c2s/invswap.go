package c2s

import (
	"grpg/data-go/gbuf"
	"log"
	"server/network"
	"server/network/s2c"
	"server/scripts"
	"server/shared"
)

type InvSwap struct {
}

func (i *InvSwap) Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player, scriptManager *scripts.ScriptManager) {
	from, err := buf.ReadByte()
	if err != nil {
		log.Printf("failed to read from byte in inv swap packet: %v\n", err)
		return
	}
	to, err := buf.ReadByte()
	if err != nil {
		log.Printf("failed to read to byte in inv swap packet: %v\n", err)
		return
	}

	// dont need to check < 0 since a byte is unsigned
	if from > 23 || to > 23 {
		log.Printf("invalid inventory slot in inv swap packet: from=%d to=%d\n", from, to)
		return
	}

	item := player.Inventory.Items[from]
	if item.ItemId == 0 {
		log.Printf("attempted to swap from empty inventory slot: from=%d to=%d\n", from, to)
		return
	}
	if player.Inventory.Items[to].ItemId != 0 {
		log.Printf("attempted to swap into non-empty inventory slot: from=%d to=%d\n", from, to)
		return
	}
	player.Inventory.Items[from] = shared.InventoryItem{}
	player.Inventory.Items[from].Dirty = true
	player.Inventory.Items[to] = item
	player.Inventory.Items[to].Dirty = true

	network.SendPacket(player.Conn, &s2c.InventoryUpdate{
		Player: player,
	}, game)
}
