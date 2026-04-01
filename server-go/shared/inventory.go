package shared

import (
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgitem"
)

type Inventory struct {
	Items [24]InventoryItem
}

type Equipment struct {
	Items [5]EquipmentItem // 0: helmet, 1: chest, 2: legs, 3: wep, 4: ring
}

type EquipmentItem struct {
	ItemId uint16
	// count is always 1 if itemid != 0
	Dirty bool
}

type InventoryItem struct {
	ItemId uint16
	Count  uint16
	Dirty  bool
}

func (i *Inventory) AddItem(item *grpgitem.Item) {
	firstEmptyIdx := -1

	for idx := range 24 {
		if i.Items[idx].ItemId == uint16(item.ItemId) && item.Stackable {
			i.Items[idx].Count++
			i.Items[idx].Dirty = true
			return
		}

		if i.Items[idx].ItemId == 0 && firstEmptyIdx == -1 {
			// doesnt break here cuz even if it finds an empty index then it can find a pre existing add to stackable,
			// which then insta returns making firstemptyidx irrelevant
			firstEmptyIdx = idx
		}
	}

	// if it finds a pre existing stack then it returns early anyway so np
	if firstEmptyIdx != -1 {
		i.Items[firstEmptyIdx].ItemId = uint16(item.ItemId)
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

func (i *Inventory) DecodeFromBlob(blob []byte) error {
	buf := gbuf.NewGBuf(blob)
	inv := [24]InventoryItem{}

	for idx := range 24 {
		id, err := buf.ReadUint16()
		if err != nil {
			return err
		}
		count, err := buf.ReadUint16()
		if err != nil {
			return err
		}

		inv[idx] = InventoryItem{
			ItemId: id,
			Count:  count,
			Dirty:  id != 0,
		}
	}

	i.Items = inv

	return nil
}

func (e *Equipment) EncodeToBlob() []byte {
	buf := gbuf.NewEmptyGBuf()

	for idx := range 5 {
		buf.WriteUint16(e.Items[idx].ItemId)
	}

	return buf.Bytes()
}

func (e *Equipment) DecodeFromBlob(blob []byte) error {
	if len(blob) == 0 {
		return nil
	}
	buf := gbuf.NewGBuf(blob)

	for idx := range 5 {
		itemId, err := buf.ReadUint16()
		if err != nil {
			return err
		}

		e.Items[idx] = EquipmentItem{
			ItemId: itemId,
			Dirty:  itemId != 0,
		}
	}

	return nil
}
