package shared

import (
	"cmp"
	"grpg/data-go/gbuf"
)

type InventoryItem struct {
	ItemId uint16
	Count  uint16
	Dirty  bool
}

func EncodeInventoryToBlob(items [24]InventoryItem) []byte {
	buf := gbuf.NewEmptyGBuf()

	for idx := range 24 {
		buf.WriteUint16(items[idx].ItemId)
		buf.WriteUint16(items[idx].Count)
	}

	return buf.Bytes()
}

func DecodeInventoryFromBlob(blob []byte) ([24]InventoryItem, error) {
	buf := gbuf.NewGBuf(blob)
	inv := [24]InventoryItem{}

	for idx := range 24 {
		id, err1 := buf.ReadUint16()
		count, err2 := buf.ReadUint16()
		if err := cmp.Or(err1, err2); err != nil {
			return [24]InventoryItem{}, err
		}

		inv[idx] = InventoryItem{
			ItemId: id,
			Count:  count,
			Dirty:  false,
		}
	}

	return inv, nil
}
