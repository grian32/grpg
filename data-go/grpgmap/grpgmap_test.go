package grpgmap

import (
	"bytes"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgtex"
	"testing"
)

func TestReadWriteHeader(t *testing.T) {
	expectedHeader := Header{
		Magic:   [8]byte{'G', 'R', 'P', 'G', 'M', 'A', 'P', 0x00},
		Version: 1,
		ChunkX:  1,
		ChunkY:  1,
	}

	expectedBytes := []byte{
		'G', 'R', 'P', 'G', 'M', 'A', 'P', 0x00, // magic
		0x00, 0x01, // version
		0x00, 0x01, // chunkX
		0x00, 0x01, // chunkY
	}

	buf := gbuf.NewEmptyGBuf()

	t.Run("WriteHeader", func(t *testing.T) {
		WriteHeader(buf, expectedHeader)

		if !bytes.Equal(buf.Bytes(), expectedBytes) {
			t.Fatalf("WriteHeader=%x, want match for %x", buf.Bytes(), expectedBytes)
		}
	})

	t.Run("ReadHeader", func(t *testing.T) {
		header, err := ReadHeader(buf)

		if header != expectedHeader || err != nil {
			t.Errorf("ReadHeader=%v, %s want match for %v", header, err.Error(), expectedHeader)
		}
	})
}

func TestWriteReadTiles(t *testing.T) {
	tileArr := [256]Tile{}
	expectedBytes := [768]byte{} // (2 bytes for uint16, 1 byte for textype byte) * 256 = 768

	for idx := range 128 {
		tileArr[idx] = Tile{0, grpgtex.TILE}
		offset := idx * 3
		expectedBytes[offset] = 0x00
		expectedBytes[offset+1] = 0x00
		expectedBytes[offset+2] = 0x01
	}

	for idx := 128; idx < 256; idx++ {
		tileArr[idx] = Tile{1, grpgtex.OBJ}
		offset := idx * 3
		expectedBytes[offset] = 0x00
		expectedBytes[offset+1] = 0x01
		expectedBytes[offset+2] = 0x02
	}

	buf := gbuf.NewEmptyGBuf()

	t.Run("WriteTiles", func(t *testing.T) {
		WriteTiles(buf, tileArr)

		if !bytes.Equal(buf.Bytes(), expectedBytes[:]) {
			t.Fatalf("WriteHeader=%x, want match for %x", buf.Bytes(), expectedBytes)
		}
	})

	t.Run("ReadTiles", func(t *testing.T) {
		tiles, err := ReadTiles(buf)

		if tiles != tileArr || err != nil {
			t.Fatalf("WriteHeader=%v, want match for %#v", tiles, tileArr)
		}
	})
}
