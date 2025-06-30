package grpgtex

import (
	"bytes"
	"encoding/binary"
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

func WriteHeader(buf *bytes.Buffer, version uint16) error {
	header := Header{
		Magic:   [8]byte{'G', 'R', 'P', 'G', 'T', 'E', 'X', 0},
		Version: version,
	}
	err := binary.Write(buf, binary.BigEndian, header)

	return err
}

func WriteTextures(buf *bytes.Buffer, textures []Texture) error {
	err := binary.Write(buf, binary.BigEndian, uint32(len(textures)))
	if err != nil {
		return err
	}

	// can add length checking for lengths being uint32 if it becomes an issue but that seems very unlikely lol..
	for _, tex := range textures {
		err = binary.Write(buf, binary.BigEndian, uint32(len(tex.InternalIdData)))

		if err != nil {
			return err
		}
		// using buf.write cuz not fixed length lol
		buf.Write(tex.InternalIdData)

		err = binary.Write(buf, binary.BigEndian, uint32(len(tex.PNGBytes)))

		if err != nil {
			return err
		}

		buf.Write(tex.PNGBytes)
	}

	return nil
}
