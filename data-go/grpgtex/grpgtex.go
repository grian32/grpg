package grpgtex

import (
	"bytes"
	"grpg/data-go/gbuf"
)

type Header struct {
	Magic   [8]byte
	Version uint16
}

type Texture struct {
	InternalIdData []byte
	PNGBytes       []byte
}

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
	}
}
