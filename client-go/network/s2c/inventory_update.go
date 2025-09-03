package s2c

import (
	"client/shared"
	"cmp"
	"grpg/data-go/gbuf"
	"log"
)

type InventoryUpdate struct{}

func (i *InventoryUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	firstMask, err1 := buf.ReadByte()
	secondMask, err2 := buf.ReadByte()
	thirdMask, err3 := buf.ReadByte()

	if firstMask == 0 && secondMask == 0 && thirdMask == 0 {
		// no dirty indexes, happens on a login for a brand new acc
		return
	}

	if err := cmp.Or(err1, err2, err3); err != nil {
		log.Printf("failed to read mask bytes in inv update: %v", err)
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
}
