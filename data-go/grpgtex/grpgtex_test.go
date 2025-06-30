package grpgtex

import (
	"bytes"
	"encoding/binary"
	"grpg/data-go/gbuf"
	"io"
	"log"
	"os"
	"testing"
)

var (
	buf           = gbuf.NewEmptyGBuf()
	stone         *os.File
	grass         *os.File
	stonePngBytes []byte
	grassPngBytes []byte
)

func init() {
	var err error

	stone, err = os.Open("../testdata/stone_texture.png")
	if err != nil {
		log.Fatal("Error loading stone texture while initializing format tests")
	}
	grass, err = os.Open("../testdata/grass_texture.png")
	if err != nil {
		log.Fatal("Error loading grass texture while initializing format tests")
	}
	stonePngBytes, err = io.ReadAll(stone)
	if err != nil {
		log.Fatal("Error loading stone png bytes while initializing format tests")
	}
	grassPngBytes, err = io.ReadAll(grass)
	if err != nil {
		log.Fatal("Error loading stone png bytes while initializing format tests")
	}
}

func TestWriteHeaderVer1(t *testing.T) {
	expectedBytes := []byte{
		'G', 'R', 'P', 'G', 'T', 'E', 'X', 0, // magic
		0x00, 0x01, // ver1
	}

	WriteHeader(buf, 1)
	if !bytes.Equal(expectedBytes, buf.Bytes()) {
		t.Errorf("WriteHeader(buf, 1)= %q, want match for %#q", buf.Bytes(), expectedBytes)
	}
	buf.Clear()
}

func TestWriteHeaderVerMax(t *testing.T) {
	expectedBytes := []byte{
		'G', 'R', 'P', 'G', 'T', 'E', 'X', 0, // magic
		0xFF, 0xFF, // ver1
	}

	WriteHeader(buf, 65535)
	if !bytes.Equal(expectedBytes, buf.Bytes()) {
		t.Errorf("WriteHeader(buf, 1)= %q, want match for %#q", buf.Bytes(), expectedBytes)
	}
	buf.Clear()
}

func TestWriteTextures(t *testing.T) {
	input := []Texture{
		{
			InternalIdData: []byte("grass"),
			PNGBytes:       grassPngBytes,
			Type:           TILE,
		},
		{
			InternalIdData: []byte("stone"),
			PNGBytes:       stonePngBytes,
			Type:           OBJ,
		},
	}

	expectedBytes := []byte{0x00, 0x00, 0x00, 0x02 /* 2 textures len */}

	for _, tex := range input {
		expectedBytes = append(expectedBytes, uint32ToBytes(len(tex.InternalIdData))...)
		expectedBytes = append(expectedBytes, tex.InternalIdData...)
		expectedBytes = append(expectedBytes, uint32ToBytes(len(tex.PNGBytes))...)
		expectedBytes = append(expectedBytes, tex.PNGBytes...)
		expectedBytes = append(expectedBytes, byte(tex.Type))
	}

	WriteTextures(buf, input)

	if !bytes.Equal(expectedBytes, buf.Bytes()) {
		t.Errorf("WriteGRPGTex= %q, want match for %#q", buf.Bytes(), expectedBytes)
	}
}

func TestReadHeaderVer1(t *testing.T) {
	expectedHeader := Header{
		Magic:   [8]byte{'G', 'R', 'P', 'G', 'T', 'E', 'X', 0},
		Version: 1,
	}

	buf := gbuf.NewGBuf([]byte{
		'G', 'R', 'P', 'G', 'T', 'E', 'X', 0, // magic
		0x00, 0x01, // ver1
	})

	output := ReadHeader(buf)

	if output != expectedHeader {
		t.Errorf("ReadHeader=%q, want match for %#q", output, expectedHeader)
	}
}

func TestReadHeaderVerMax(t *testing.T) {
	expectedHeader := Header{
		Magic:   [8]byte{'G', 'R', 'P', 'G', 'T', 'E', 'X', 0},
		Version: 65535,
	}

	buf := gbuf.NewGBuf([]byte{
		'G', 'R', 'P', 'G', 'T', 'E', 'X', 0, // magic
		0xFF, 0xFF, // ver1
	})

	output := ReadHeader(buf)

	if output != expectedHeader {
		t.Errorf("ReadHeader=%q, want match for %#q", output, expectedHeader)
	}
}

func TestReadTextures(t *testing.T) {
	expected := []Texture{
		{
			InternalIdData: []byte("grass"),
			PNGBytes:       grassPngBytes,
			Type:           TILE,
		},
		{
			InternalIdData: []byte("stone"),
			PNGBytes:       stonePngBytes,
			Type:           OBJ,
		},
	}

	buf := gbuf.NewEmptyGBuf()

	buf.WriteUint32(2)

	buf.WriteUint32(5)
	buf.WriteBytes([]byte("grass"))
	buf.WriteUint32(uint32(len(grassPngBytes)))
	buf.WriteBytes(grassPngBytes)
	buf.WriteByte(byte(TILE))

	buf.WriteUint32(5)
	buf.WriteBytes([]byte("stone"))
	buf.WriteUint32(uint32(len(stonePngBytes)))
	buf.WriteBytes(stonePngBytes)
	buf.WriteByte(byte(OBJ))

	output := ReadTextures(buf)

	if !output[0].Equals(expected[0]) || !output[1].Equals(expected[1]) {
		t.Errorf("ReadHeader=%q, want match for %#q", output, expected)
	}
}
func uint32ToBytes(u int) []byte {
	arr := make([]byte, 4)
	binary.BigEndian.PutUint32(arr, uint32(u))

	return arr
}
