package s2c

import (
	"client/shared"
	"cmp"
	"grpg/data-go/gbuf"
	"log"
)

type InventoryUpdate struct{}

func (i *InventoryUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	firstMask, err := buf.ReadByte()
	if err != nil {
		log.Printf("failed to read first mask byte in inv update: %v", err)
		return
	}
	secondMask, err := buf.ReadByte()
	if err != nil {
		log.Printf("failed to read second mask byte in inv update: %v", err)
		return
	}
	thirdMask, err := buf.ReadByte()
	if err != nil {
		log.Printf("failed to read third mask byte in inv update: %v", err)
		return
	}
	fourthMask, err := buf.ReadByte()
	if err != nil {
		log.Printf("failed to read fourth mask byte in inv update: %v", err)
		return
	}

	if firstMask == 0 && secondMask == 0 && thirdMask == 0 && fourthMask == 0 {
		// no dirty indexes, happens on a login for a brand new acc
		return
	}

	for idx := range 24 {
		var mask byte
		offset := 0
		if idx < 8 {
			mask = firstMask
			offset = 0
		} else if idx < 16 {
			mask = secondMask
			offset = 8
		} else {
			mask = thirdMask
			offset = 16
		}

		if mask&(1<<(idx-offset)) != 0 {
			item, err1 := buf.ReadUint16()
			count, err2 := buf.ReadUint16()
			if err := cmp.Or(err1, err2); err != nil {
				log.Printf("failed to read item & count in inv update: %v", err)
				return
			}

			game.Player.Inventory[idx] = shared.InventoryItem{
				ItemId: item,
				Count:  count,
			}
		}
	}

	if game.Player.Equipment == nil {
		game.Player.Equipment = make(map[shared.EquipmentType]uint16)
	}
	if fourthMask != 0 {
		for idx := range 5 {
			if (fourthMask & (1 << idx)) != 0 {
				item, err := buf.ReadUint16()
				if err != nil {
					log.Printf("failed to read equipment item in inv update: %v", err)
					return
				}
				game.Player.Equipment[shared.EquipmentType(idx)] = item
			}
		}
	}
}
