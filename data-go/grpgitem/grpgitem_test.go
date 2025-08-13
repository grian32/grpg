package grpgitem

import (
	"bytes"
	"grpg/data-go/gbuf"
	"slices"
	"testing"
)

func TestReadWriteHeader(t *testing.T) {
	expectedBytes := []byte("GRPGITEM")
	expectedHeader := Header{Magic: [8]byte{
		'G', 'R', 'P', 'G', 'I', 'T', 'E', 'M',
	}}

	buf := gbuf.NewEmptyGBuf()

	t.Run("WriteHeader", func(t *testing.T) {
		WriteHeader(buf)
		if !bytes.Equal(buf.Bytes(), expectedBytes) {
			t.Errorf("WriteHeader=%v, want match for %v", buf.Bytes(), expectedBytes)
		}
	})

	t.Run("ReadHeader", func(t *testing.T) {
		header, err := ReadHeader(buf)
		if header != expectedHeader || err != nil {
			t.Errorf("ReadHeader=%v, %v, watch match for %v", header, err, expectedHeader)
		}
	})
}

func TestReadWriteItems(t *testing.T) {
	expectedBytes := []byte{
		0x00, 0x02, // len

		0x00, 0x01, // item id
		0x00, 0x03, // item tex
		0x00, 0x00, 0x00, 0x07, 'b', 'e', 'r', 'r', 'i', 'e', 's', // item name

		0x00, 0x02, // item id
		0x00, 0x03, // item tex
		0x00, 0x00, 0x00, 0x08, 'b', 'e', 'r', 'r', 'i', 'e', 's', 's', // item name
	}

	expectedItems := []Item{
		{
			ItemId:  1,
			Texture: 3,
			Name:    "berries",
		},
		{
			ItemId:  2,
			Texture: 3,
			Name:    "berriess",
		},
	}

	buf := gbuf.NewEmptyGBuf()

	t.Run("WriteItems", func(t *testing.T) {
		WriteItems(buf, expectedItems)

		if !bytes.Equal(buf.Bytes(), expectedBytes) {
			t.Errorf("WriteItems=%v, want match for %v", buf.Bytes(), expectedBytes)
		}
	})

	t.Run("ReadItems", func(t *testing.T) {
		items, err := ReadItems(buf)

		if !slices.Equal(items, expectedItems) || err != nil {
			t.Errorf("ReadItems=%v, %v, watch match for %v", items, err, expectedItems)
		}
	})
}
