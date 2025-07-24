package grpgmap

import (
	"cmp"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgtex"
)

type Header struct {
	Magic   [8]byte
	Version uint16
	ChunkX  uint16
	ChunkY  uint16
}

type Tile struct {
	InternalId uint16
	Type       grpgtex.TextureType
}

func WriteHeader(buf *gbuf.GBuf, header Header) {
	buf.WriteBytes(header.Magic[:])
	buf.WriteUint16(header.Version)
	buf.WriteUint16(header.ChunkX)
	buf.WriteUint16(header.ChunkY)
}

func ReadHeader(buf *gbuf.GBuf) (Header, error) {
	magic, err1 := buf.ReadBytes(8)
	version, err2 := buf.ReadUint16()
	chunkX, err3 := buf.ReadUint16()
	chunkY, err4 := buf.ReadUint16()
	if err := cmp.Or(err1, err2, err3, err4); err != nil {
		return Header{}, err
	}

	return Header{
		Magic:   [8]byte(magic),
		Version: version,
		ChunkX:  chunkX,
		ChunkY:  chunkY,
	}, nil
}

func WriteTiles(buf *gbuf.GBuf, tiles [256]Tile) {
	for _, tile := range tiles {
		buf.WriteUint16(tile.InternalId)
		buf.WriteByte(byte(tile.Type))
	}
}

func ReadTiles(buf *gbuf.GBuf) ([256]Tile, error) {
	arr := [256]Tile{}

	for idx := range 256 {
		internalId, err1 := buf.ReadUint16()
		texType, err2 := buf.ReadByte()
		if err := cmp.Or(err1, err2); err != nil {
			return [256]Tile{}, err
		}

		tile := Tile{
			InternalId: internalId,
			Type:       grpgtex.TextureType(texType),
		}

		arr[idx] = tile
	}

	return arr, nil
}
