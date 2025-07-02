package grpgmap

import (
	"bytes"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgtex"
	"testing"
)

func TestWriteHeader(t *testing.T) {
	header := Header{
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
	WriteHeader(buf, header)

	if !bytes.Equal(buf.Bytes(), expectedBytes) {
		t.Errorf("WriteHeader(buf, 1)= %q, want match for %#q", buf.Bytes(), expectedBytes)
	}
}

func TestWriteTiles(t *testing.T) {
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
	WriteTiles(buf, tileArr)

	if !bytes.Equal(buf.Bytes(), expectedBytes[:]) {
		t.Errorf("WriteHeader(buf, 1)= %q, want match for %#q", buf.Bytes(), expectedBytes)
	}
}
