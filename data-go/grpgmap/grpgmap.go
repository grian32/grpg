package grpgmap

import (
	"cmp"
	"grpg/data-go/gbuf"
)

type ObjType byte

const (
	UNDEFINED ObjType = iota
	OBJ
)

type Header struct {
	Magic  [8]byte
	ChunkX uint16
	ChunkY uint16
}

type Obj struct {
	InternalId uint16
	Type       ObjType
}

type Tile uint16

type Zone struct {
	Tiles [256]Tile
	Objs  [256]Obj
}

func WriteHeader(buf *gbuf.GBuf, header Header) {
	buf.WriteBytes(header.Magic[:])
	buf.WriteUint16(header.ChunkX)
	buf.WriteUint16(header.ChunkY)
}

func ReadHeader(buf *gbuf.GBuf) (Header, error) {
	magic, err1 := buf.ReadBytes(8)
	chunkX, err2 := buf.ReadUint16()
	chunkY, err3 := buf.ReadUint16()
	if err := cmp.Or(err1, err2, err3); err != nil {
		return Header{}, err
	}

	return Header{
		Magic:  [8]byte(magic),
		ChunkX: chunkX,
		ChunkY: chunkY,
	}, nil
}

func WriteZone(buf *gbuf.GBuf, zone Zone) {
	for _, tile := range zone.Tiles {
		buf.WriteUint16(uint16(tile))
	}

	for _, obj := range zone.Objs {
		buf.WriteUint16(obj.InternalId)
		buf.WriteByte(byte(obj.Type))
	}
}

func ReadZone(buf *gbuf.GBuf) (Zone, error) {
	tiles := [256]Tile{}
	objs := [256]Obj{}

	for idx := range 256 {
		internalId, err := buf.ReadUint16()
		if err != nil {
			return Zone{}, err
		}

		tiles[idx] = Tile(internalId)
	}

	for idx := range 256 {
		internalId, err1 := buf.ReadUint16()
		texType, err2 := buf.ReadByte()
		if err := cmp.Or(err1, err2); err != nil {
			return Zone{}, err
		}

		objs[idx] = Obj{
			InternalId: internalId,
			Type:       ObjType(texType),
		}
	}

	zone := Zone{
		Tiles: tiles,
		Objs:  objs,
	}

	return zone, nil
}
