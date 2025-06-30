package grpgtex

import (
	"bytes"
	"grpg/data-go/gbuf"
	"log"
)

type Header struct {
	Magic   [8]byte
	Version uint16
}

type Texture struct {
	InternalIdData []byte
	PNGBytes       []byte
	Type           TextureType
}

type TextureType byte

const (
	UNDEFINED TextureType = 0x00
	TILE      TextureType = 0x01
	OBJ       TextureType = 0x02
)

func (t Texture) Equals(other Texture) bool {
	return bytes.Equal(t.InternalIdData, other.InternalIdData) && bytes.Equal(t.PNGBytes, other.PNGBytes)
}

func WriteHeader(buf *gbuf.GBuf, version uint16) {
	header := Header{
		Magic:   [8]byte{'G', 'R', 'P', 'G', 'T', 'E', 'X', 0},
		Version: version,
	}
	buf.WriteBytes(header.Magic[:])
	buf.WriteUint16(version)
}

func WriteTextures(buf *gbuf.GBuf, textures []Texture) {
	buf.WriteUint32(uint32(len(textures)))

	// can add length checking for lengths being uint32 if it becomes an issue but that seems very unlikely lol..
	for _, tex := range textures {
		buf.WriteUint32(uint32(len(tex.InternalIdData)))
		buf.WriteBytes(tex.InternalIdData)

		buf.WriteUint32(uint32(len(tex.PNGBytes)))
		buf.WriteBytes(tex.PNGBytes)

		buf.WriteByte(byte(tex.Type))
	}
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
	return Header{
		Magic:   [8]byte(magic),
		Version: version,
	}
}

func ReadTextures(buf *gbuf.GBuf) []Texture {
	var textures []Texture

	textureLen, err := buf.ReadUint32()
	if err != nil {
		log.Fatal(err)
	}

	for range textureLen {
		internalIdLen, err := buf.ReadUint32()
		if err != nil {
			log.Fatal(err)
		}
		internalIdData, err := buf.ReadBytes(int(internalIdLen))
		if err != nil {
			log.Fatal(err)
		}

		pngBytesLen, err := buf.ReadUint32()
		if err != nil {
			log.Fatal(err)
		}
		pngBytes, err := buf.ReadBytes(int(pngBytesLen))
		if err != nil {
			log.Fatal(err)
		}

		texType, err := buf.ReadByte()
		if err != nil {
			log.Fatal(err)
		}

		textures = append(textures, Texture{
			InternalIdData: internalIdData,
			PNGBytes:       pngBytes,
			Type:           TextureType(texType),
		})
	}

	return textures
}
