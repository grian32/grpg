package grpgmap

import (
	"bytes"
	"grpg/data-go/gbuf"
	"testing"
)

func TestReadWriteHeader(t *testing.T) {
	expectedHeader := Header{
		Magic:  [8]byte{'G', 'R', 'P', 'G', 'M', 'A', 'P', 0x00},
		ChunkX: 1,
		ChunkY: 1,
	}

	expectedBytes := []byte{
		'G', 'R', 'P', 'G', 'M', 'A', 'P', 0x00, // magic
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

func TestWriteReadZone(t *testing.T) {
	expectedZone := Zone{}
	expectedBytes := [1024]byte{} // 2 * 256 for tile layer, 2 * 256 for obj layer

	for idx := range 128 {
		expectedZone.Tiles[uint16(idx)] = 0x00

		tileOffset := idx * 2

		expectedBytes[tileOffset] = 0x00
		expectedBytes[tileOffset+1] = 0x00

		objOffset := (idx * 2) + 512

		expectedZone.Objs[idx] = 0x00

		expectedBytes[objOffset] = 0x00
		expectedBytes[objOffset+1] = 0x00
	}

	for idx := 128; idx < 256; idx++ {
		expectedZone.Tiles[uint16(idx)] = 0x01

		tileOffset := idx * 2

		expectedBytes[tileOffset] = 0x00
		expectedBytes[tileOffset+1] = 0x01

		objOffset := (idx * 2) + 512

		expectedZone.Objs[idx] = 1

		expectedBytes[objOffset] = 0x00
		expectedBytes[objOffset+1] = 0x01
	}

	buf := gbuf.NewEmptyGBuf()

	t.Run("WriteZone", func(t *testing.T) {
		WriteZone(buf, expectedZone)

		if !bytes.Equal(buf.Bytes(), expectedBytes[:]) {
			t.Fatalf("WriteZone=%x, want match for %x", buf.Bytes(), expectedBytes[:])
		}
	})

	t.Run("ReadZone", func(t *testing.T) {
		zone, err := ReadZone(buf)

		if zone != expectedZone || err != nil {
			t.Fatalf("ReadZone=%v, want match for %#v", zone, expectedZone)
		}
	})
}
