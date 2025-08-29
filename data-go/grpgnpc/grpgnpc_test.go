package grpgnpc

import (
	"bytes"
	"grpg/data-go/gbuf"
	"slices"
	"testing"
)

func TestReadWriteHeader(t *testing.T) {
	expectedHeader := Header{
		Magic: [8]byte{'G', 'R', 'P', 'G', 'N', 'P', 'C', 0x00},
	}

	expectedBytes := []byte("GRPGNPC\x00")

	buf := gbuf.NewEmptyGBuf()

	t.Run("WriteHeader", func(t *testing.T) {
		WriteHeader(buf)

		if !bytes.Equal(expectedBytes, buf.Bytes()) {
			t.Errorf("WriteHeader=%v, wanted match for %v", buf.Bytes(), expectedBytes)
		}
	})

	t.Run("ReadHeader", func(t *testing.T) {
		header, err := ReadHeader(buf)

		if header != expectedHeader || err != nil {
			t.Errorf("ReadHeader=%v, %v, wanted match for %v", header, err, expectedHeader)
		}
	})
}

func TestReadWriteNpcs(t *testing.T) {
	expectedNpcs := []Npc{
		{
			NpcId:     1,
			Name:      "corey",
			TextureId: 8,
		},
		{
			NpcId:     2,
			Name:      "grian",
			TextureId: 9,
		},
	}

	expectedBytes := []byte{
		0x00, 0x02, // npc len

		0x00, 0x01, // npc id
		0x00, 0x00, 0x00, 0x05, 'c', 'o', 'r', 'e', 'y',
		0x00, 0x08,

		0x00, 0x02, // npc id
		0x00, 0x00, 0x00, 0x05, 'g', 'r', 'i', 'a', 'n', // npc name
		0x00, 0x09,
	}

	buf := gbuf.NewEmptyGBuf()

	t.Run("WriteNpcs", func(t *testing.T) {
		WriteNpcs(buf, expectedNpcs)

		if !bytes.Equal(buf.Bytes(), expectedBytes) {
			t.Errorf("WriteNpcs=%v, wanted match for %v", buf.Bytes(), expectedBytes)
		}
	})

	t.Run("ReadNpcs", func(t *testing.T) {
		npcs, err := ReadNpcs(buf)
		if !slices.Equal(npcs, expectedNpcs) || err != nil {
			t.Errorf("ReadNpcs=%v, %v, wanted match for %v", npcs, err, expectedNpcs)
		}
	})
}
