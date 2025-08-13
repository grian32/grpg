package grpgitem

import (
	"cmp"
	"grpg/data-go/gbuf"
)

type Header struct {
	Magic [8]byte
}

type Item struct {
	ItemId  uint16
	Texture uint16
	// maybe stackable etc in the future
	Name string
}

func WriteHeader(buf *gbuf.GBuf) {
	buf.WriteBytes([]byte{'G', 'R', 'P', 'G', 'I', 'T', 'E', 'M'})
}

func ReadHeader(buf *gbuf.GBuf) (Header, error) {
	magic, err := buf.ReadBytes(8)
	if err != nil {
		return Header{}, err
	}

	return Header{
		Magic: [8]byte(magic),
	}, nil
}

func WriteItems(buf *gbuf.GBuf, items []Item) {
	buf.WriteUint16(uint16(len(items)))

	for _, item := range items {
		buf.WriteUint16(item.ItemId)
		buf.WriteUint16(item.Texture)
		buf.WriteString(item.Name)
	}
}

func ReadItems(buf *gbuf.GBuf) ([]Item, error) {
	itemLen, err := buf.ReadUint16()
	if err != nil {
		return nil, err
	}

	itemArr := make([]Item, itemLen)

	for idx := range itemLen {
		itemId, err1 := buf.ReadUint16()
		textureId, err2 := buf.ReadUint16()
		itemName, err3 := buf.ReadString()

		if err := cmp.Or(err1, err2, err3); err != nil {
			return nil, err
		}

		itemArr[idx] = Item{
			ItemId:  itemId,
			Texture: textureId,
			Name:    itemName,
		}
	}

	return itemArr, nil
}
