package grpgtile

import (
	"cmp"
	"grpg/data-go/gbuf"
)

type Header struct {
	Magic [8]byte
}

type Tile struct {
	Name   string
	TileId uint16
	TexId  uint16
}

func WriteHeader(buf *gbuf.GBuf) {
	buf.WriteBytes([]byte("GRPGTILE"))
}

func ReadHeader(buf *gbuf.GBuf) (Header, error) {
	bytes, err := buf.ReadBytes(8)
	if err != nil {
		return Header{}, err
	}

	return Header{
		Magic: [8]byte(bytes),
	}, nil
}

func WriteTiles(buf *gbuf.GBuf, tiles []Tile) {
	buf.WriteUint16(uint16(len(tiles)))

	for _, tile := range tiles {
		buf.WriteString(tile.Name)
		buf.WriteUint16(tile.TileId)
		buf.WriteUint16(tile.TexId)
	}
}

func ReadTiles(buf *gbuf.GBuf) ([]Tile, error) {
	len, err := buf.ReadUint16()
	if err != nil {
		return nil, err
	}

	arr := make([]Tile, len)

	for idx := range len {
		name, err1 := buf.ReadString()
		tileId, err2 := buf.ReadUint16()
		texId, err3 := buf.ReadUint16()

		if err := cmp.Or(err1, err2, err3); err != nil {
			return nil, err
		}

		arr[idx] = Tile{
			Name:   name,
			TileId: tileId,
			TexId:  texId,
		}
	}

	return arr, nil
}
