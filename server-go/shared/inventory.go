package shared

import (
	"cmp"
	"grpg/data-go/gbuf"
)

type Inventory struct {
	Items [24]InventoryItem
}

type InventoryItem struct {
	ItemId uint16
	Count  uint16
	Dirty  bool
}

func (i *Inventory) AddItem(itemId uint16) {
	firstEmptyIdx := -1

	for idx := range 24 {
		if i.Items[idx].ItemId == uint16(itemId) {
			i.Items[idx].Count++
			i.Items[idx].Dirty = true
			return
		}

		if i.Items[idx].ItemId == 0 && firstEmptyIdx == -1 {
			firstEmptyIdx = idx
		}
	}

	// if it finds a pre existing stack then it returns early anyway so np
	if firstEmptyIdx != -1 {
		i.Items[firstEmptyIdx].ItemId = uint16(itemId)
		i.Items[firstEmptyIdx].Count = 1
		i.Items[firstEmptyIdx].Dirty = true
	}

}

func (i *Inventory) EncodeToBlob() []byte {
	buf := gbuf.NewEmptyGBuf()

	for idx := range 24 {
		buf.WriteUint16(i.Items[idx].ItemId)
		buf.WriteUint16(i.Items[idx].Count)
	}

	return buf.Bytes()
}

func DecodeInventoryFromBlob(blob []byte) (Inventory, error) {
	buf := gbuf.NewGBuf(blob)
	inv := [24]InventoryItem{}

	for idx := range 24 {
		id, err1 := buf.ReadUint16()
		count, err2 := buf.ReadUint16()
		if err := cmp.Or(err1, err2); err != nil {
			return Inventory{}, err
		}

		dirty := false
		if id != 0 {
			dirty = true
		}

		inv[idx] = InventoryItem{
			ItemId: id,
			Count:  count,
			Dirty:  dirty,
		}
	}

	return Inventory{Items: inv}, nil
}
