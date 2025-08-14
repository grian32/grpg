package grpgtex

import (
	"bytes"
	"encoding/binary"
	"grpg/data-go/gbuf"
	"log"
	"os"
	"testing"
)

var (
	buf           = gbuf.NewEmptyGBuf()
	stonePngBytes []byte
	grassPngBytes []byte
)

func init() {
	var err error

	stonePngBytes, err = os.ReadFile("../testdata/stone_texture.png")
	if err != nil {
		log.Fatal("Error loading stone png bytes while initializing format tests")
	}
	grassPngBytes, err = os.ReadFile("../testdata/grass_texture.png")
	if err != nil {
		log.Fatal("Error loading grass png bytes while initializing format tests")
	}
}

func TestWriteHeader(t *testing.T) {
	expectedBytes := []byte{
		'G', 'R', 'P', 'G', 'T', 'E', 'X', 0, // magic
	}

	WriteHeader(buf)

	if !bytes.Equal(expectedBytes, buf.Bytes()) {
		t.Errorf("WriteHeader(buf, 1)= %q, want match for %#q", buf.Bytes(), expectedBytes)
	}
	buf.Clear()
}

func TestWriteTextures(t *testing.T) {
	input := []Texture{
		{
			InternalIdString: []byte("grass"),
			InternalIdInt:    0,
			PNGBytes:         grassPngBytes,
		},
		{
			InternalIdString: []byte("stone"),
			InternalIdInt:    1,
			PNGBytes:         stonePngBytes,
		},
	}

	expectedBytes := []byte{0x00, 0x00, 0x00, 0x02 /* 2 textures len */}

	for _, tex := range input {
		expectedBytes = append(expectedBytes, uint32ToBytes(len(tex.InternalIdString))...)
		expectedBytes = append(expectedBytes, tex.InternalIdString...)
		expectedBytes = append(expectedBytes, uint16ToBytes(int(tex.InternalIdInt))...)
		expectedBytes = append(expectedBytes, uint32ToBytes(len(tex.PNGBytes))...)
		expectedBytes = append(expectedBytes, tex.PNGBytes...)
	}

	WriteTextures(buf, input)

	if !bytes.Equal(expectedBytes, buf.Bytes()) {
		t.Errorf("WriteGRPGTex= %q, want match for %#q", buf.Bytes(), expectedBytes)
	}
}

func TestReadHeader(t *testing.T) {
	expectedHeader := Header{
		Magic: [8]byte{'G', 'R', 'P', 'G', 'T', 'E', 'X', 0},
	}

	buf := gbuf.NewGBuf([]byte{
		'G', 'R', 'P', 'G', 'T', 'E', 'X', 0, // magic
	})

	output, err := ReadHeader(buf)
	if err != nil {
		t.Errorf("ReadHeader errored: %v", err)
	}

	if output != expectedHeader {
		t.Errorf("ReadHeader=%q, want match for %#q", output, expectedHeader)
	}
}

func TestReadTextures(t *testing.T) {
	expected := []Texture{
		{
			InternalIdString: []byte("grass"),
			InternalIdInt:    0,
			PNGBytes:         grassPngBytes,
		},
		{
			InternalIdString: []byte("stone"),
			InternalIdInt:    1,
			PNGBytes:         stonePngBytes,
		},
	}

	buf := gbuf.NewEmptyGBuf()

	buf.WriteUint32(2)

	buf.WriteUint32(5)
	buf.WriteBytes([]byte("grass"))
	buf.WriteUint16(0)
	buf.WriteUint32(uint32(len(grassPngBytes)))
	buf.WriteBytes(grassPngBytes)

	buf.WriteUint32(5)
	buf.WriteBytes([]byte("stone"))
	buf.WriteUint16(1)
	buf.WriteUint32(uint32(len(stonePngBytes)))
	buf.WriteBytes(stonePngBytes)

	output, err := ReadTextures(buf)
	if err != nil {
		t.Errorf("ReadTextures errored: %v", err)
	}

	if !output[0].Equals(expected[0]) || !output[1].Equals(expected[1]) {
		t.Errorf("ReadHeader=%q, want match for %#q", output, expected)
	}
}

func uint32ToBytes(u int) []byte {
	arr := make([]byte, 4)
	binary.BigEndian.PutUint32(arr, uint32(u))

	return arr
}

func uint16ToBytes(u int) []byte {
	arr := make([]byte, 2)
	binary.BigEndian.PutUint16(arr, uint16(u))

	return arr
}
