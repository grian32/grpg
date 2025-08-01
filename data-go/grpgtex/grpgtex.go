package grpgtex

import (
	"bytes"
	"cmp"
	"grpg/data-go/gbuf"
)

type Header struct {
	Magic   [8]byte
	Version uint16
}

type Texture struct {
	InternalIdString []byte
	InternalIdInt    uint16
	PNGBytes         []byte
	Type             TextureType
}

type TextureType byte

const (
	UNDEFINED TextureType = 0x00
	TILE      TextureType = 0x01
	OBJ       TextureType = 0x02
)

func (t Texture) Equals(other Texture) bool {
	return bytes.Equal(t.InternalIdString, other.InternalIdString) && bytes.Equal(t.PNGBytes, other.PNGBytes)
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
		buf.WriteUint32(uint32(len(tex.InternalIdString)))
		buf.WriteBytes(tex.InternalIdString)

		buf.WriteUint16(tex.InternalIdInt)

		buf.WriteUint32(uint32(len(tex.PNGBytes)))
		buf.WriteBytes(tex.PNGBytes)

		buf.WriteByte(byte(tex.Type))
	}
}

func ReadHeader(buf *gbuf.GBuf) (Header, error) {
	magic, err1 := buf.ReadBytes(8)
	version, err2 := buf.ReadUint16()

	if err := cmp.Or(err1, err2); err != nil {
		return Header{}, err
	}

	return Header{
		Magic:   [8]byte(magic),
		Version: version,
	}, nil
}

func ReadTextures(buf *gbuf.GBuf) ([]Texture, error) {
	var textures []Texture

	textureLen, err := buf.ReadUint32()
	if err != nil {
		return nil, err
	}

	for range textureLen {
		internalIdLen, err1 := buf.ReadUint32()
		internalIdString, err2 := buf.ReadBytes(int(internalIdLen))
		internalIdInt, err3 := buf.ReadUint16()
		pngBytesLen, err4 := buf.ReadUint32()
		pngBytes, err5 := buf.ReadBytes(int(pngBytesLen))
		texType, err6 := buf.ReadByte()

		if err := cmp.Or(err1, err2, err3, err4, err5, err6); err != nil {
			return nil, err
		}

		textures = append(textures, Texture{
			InternalIdString: internalIdString,
			InternalIdInt:    internalIdInt,
			PNGBytes:         pngBytes,
			Type:             TextureType(texType),
		})
	}

	return textures, nil
}
