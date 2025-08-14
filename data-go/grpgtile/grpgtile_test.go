package grpgtile

import (
	"bytes"
	"grpg/data-go/gbuf"
	"slices"
	"testing"
)

func TestReadWriteHeader(t *testing.T) {
	expectedHeader := Header{
		Magic: [8]byte{'G', 'R', 'P', 'G', 'T', 'I', 'L', 'E'},
	}

	expectedBytes := []byte{'G', 'R', 'P', 'G', 'T', 'I', 'L', 'E'}

	buf := gbuf.NewEmptyGBuf()

	t.Run("WriteHeader", func(t *testing.T) {
		WriteHeader(buf)

		if !bytes.Equal(buf.Bytes(), expectedBytes) {
			t.Errorf("WriteHeader=%v, want=%v", buf.Bytes(), expectedBytes)
		}
	})

	t.Run("ReadHeader", func(t *testing.T) {
		header, err := ReadHeader(buf)

		if header != expectedHeader || err != nil {
			t.Errorf("ReadHeader=%v,%s, want=%v", header, err.Error(), expectedHeader)
		}
	})
}

func TestReadWriteTiles(t *testing.T) {
	expectedTiles := []Tile{
		{
			Name:   "grass",
			TileId: 1,
			TexId:  2,
		},
		{
			Name:   "water",
			TileId: 2,
			TexId:  6,
		},
	}

	expectedBytes := []byte{
		0x00, 0x02, // tile arr len
		0x00, 0x00, 0x00, 0x05, 'g', 'r', 'a', 's', 's', // grass str
		0x00, 0x01, 0x00, 0x02, // tileid, texid
		0x00, 0x00, 0x00, 0x05, 'w', 'a', 't', 'e', 'r',
		0x00, 0x02, 0x00, 0x06, // tileid, texid
	}

	buf := gbuf.NewEmptyGBuf()

	t.Run("WriteTiles", func(t *testing.T) {
		WriteTiles(buf, expectedTiles)

		if !bytes.Equal(buf.Bytes(), expectedBytes) {
			t.Errorf("WriteTiles=%v, want=%v", buf.Bytes(), expectedBytes)
		}
	})

	t.Run("ReadTiles", func(t *testing.T) {
		tiles, err := ReadTiles(buf)

		if !slices.Equal(tiles, expectedTiles) || err != nil {
			t.Errorf("ReadHeader=%v,%s, want=%v", tiles, err.Error(), expectedTiles)
		}
	})
}
