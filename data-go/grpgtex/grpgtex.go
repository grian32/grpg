package grpgtex

import (
	"bytes"
	"cmp"
	"grpg/data-go/gbuf"
)

type Header struct {
	Magic [8]byte
}

type Texture struct {
	InternalIdString []byte
	InternalIdInt    uint16
	ImageBytes       []byte
}

func (t Texture) Equals(other Texture) bool {
	return bytes.Equal(t.InternalIdString, other.InternalIdString) && bytes.Equal(t.ImageBytes, other.ImageBytes)
}

func WriteHeader(buf *gbuf.GBuf) {
	header := Header{
		Magic: [8]byte{'G', 'R', 'P', 'G', 'T', 'E', 'X', 0},
	}
	buf.WriteBytes(header.Magic[:])

}

func WriteTextures(buf *gbuf.GBuf, textures []Texture) {
	buf.WriteUint32(uint32(len(textures)))

	// can add length checking for lengths being uint32 if it becomes an issue but that seems very unlikely lol..
	for _, tex := range textures {
		buf.WriteUint32(uint32(len(tex.InternalIdString)))
		buf.WriteBytes(tex.InternalIdString)

		buf.WriteUint16(tex.InternalIdInt)

		buf.WriteUint32(uint32(len(tex.ImageBytes)))
		buf.WriteBytes(tex.ImageBytes)
	}
}

func ReadHeader(buf *gbuf.GBuf) (Header, error) {
	magic, err1 := buf.ReadBytes(8)

	if err := cmp.Or(err1); err != nil {
		return Header{}, err
	}

	return Header{
		Magic: [8]byte(magic),
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
		imageBytesLen, err4 := buf.ReadUint32()
		imageBytes, err5 := buf.ReadBytes(int(imageBytesLen))

		if err := cmp.Or(err1, err2, err3, err4, err5); err != nil {
			return nil, err
		}

		textures = append(textures, Texture{
			InternalIdString: internalIdString,
			InternalIdInt:    internalIdInt,
			ImageBytes:       imageBytes,
		})
	}

	return textures, nil
}
