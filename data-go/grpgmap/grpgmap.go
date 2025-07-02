package grpgmap

import (
	"grpg/data-go/gbuf"
	"log"
)

type Header struct {
	Magic   [8]byte
	Version uint16
	ChunkX  uint16
	ChunkY  uint16
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
