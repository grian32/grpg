package grpgmap

import (
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgtex"
	"log"
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

func ReadHeader(buf *gbuf.GBuf) Header {
	magic, err := buf.ReadBytes(8)
	if err != nil {
		log.Fatal(err)
	}
	version, err := buf.ReadUint16()
	if err != nil {
		log.Fatal(err)
	}
	chunkX, err := buf.ReadUint16()
	if err != nil {
		log.Fatal(err)
	}
	chunkY, err := buf.ReadUint16()
	if err != nil {
		log.Fatal(err)
	}

	return Header{
		Magic:   [8]byte(magic),
		Version: version,
		ChunkX:  chunkX,
		ChunkY:  chunkY,
	}
}

func WriteTiles(buf *gbuf.GBuf, tiles [256]Tile) {
	for _, tile := range tiles {
		buf.WriteUint16(tile.InternalId)
		buf.WriteByte(byte(tile.Type))
	}
}

func ReadTiles(buf *gbuf.GBuf) [256]Tile {
	arr := [256]Tile{}

	for idx := range 256 {
		internalId, err := buf.ReadUint16()
		if err != nil {
			log.Fatal(err)
		}
		texType, err := buf.ReadByte()
		if err != nil {
			log.Fatal(err)
		}

		tile := Tile{
			InternalId: internalId,
			Type:       grpgtex.TextureType(texType),
		}

		arr[idx] = tile
	}

	return arr
}
