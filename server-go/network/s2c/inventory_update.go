package s2c

import (
	"grpg/data-go/gbuf"
	"server/shared"
)

type InventoryUpdate struct {
	Player *shared.Player
}

func (i *InventoryUpdate) Opcode() byte {
	return 0x05
}

func (i *InventoryUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	var packetLen uint16 = 3 // 3 mask bytes
	var dirtyIndexes []int

	var firstByte, secondByte, thirdByte byte // = 0

	for idx, item := range i.Player.Inventory {
		if item.Dirty {
			packetLen += 2 + 2 // id, count
			dirtyIndexes = append(dirtyIndexes, idx)

			if idx < 8 {
				firstByte |= 1 << idx
			} else if idx < 16 {
				secondByte |= 1 << (idx - 8)
			} else {
				thirdByte |= 1 << (idx - 16)
			}
		}
	}

	// TODO
	if len(dirtyIndexes) == 0 {
		return
	}

	buf.WriteUint16(packetLen)

	buf.WriteByte(firstByte)
	buf.WriteByte(secondByte)
	buf.WriteByte(thirdByte)

	for _, idx := range dirtyIndexes {
		buf.WriteUint16(i.Player.Inventory[idx].ItemId)
		buf.WriteUint16(i.Player.Inventory[idx].Count)
		i.Player.Inventory[idx].Dirty = false
	}
}
