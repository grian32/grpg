package grpgtex

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
	"testing"
)

var (
	buf           = bytes.Buffer{}
	stone         *os.File
	grass         *os.File
	stonePngBytes []byte
	grassPngBytes []byte
)

func init() {
	var err error

	stone, err = os.Open("./testdata/stone_texture.png")
	if err != nil {
		log.Fatal("Error loading stone texture while initializing format tests")
	}
	grass, err = os.Open("./testdata/grass_texture.png")
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

	err := WriteHeader(&buf, 1)
	if !bytes.Equal(expectedBytes, buf.Bytes()) || err != nil {
		t.Errorf("WriteHeader(&buf, 1)= %q, %v, want match for %#q", buf.Bytes(), err, expectedBytes)
	}
	buf.Reset()
}

func TestWriteHeaderVerMax(t *testing.T) {
	expectedBytes := []byte{
		'G', 'R', 'P', 'G', 'T', 'E', 'X', 0, // magic
		0xFF, 0xFF, // ver1
	}

	err := WriteHeader(&buf, 65535)
	if !bytes.Equal(expectedBytes, buf.Bytes()) || err != nil {
		t.Errorf("WriteHeader(&buf, 1)= %q, %v, want match for %#q", buf.Bytes(), err, expectedBytes)
	}
	buf.Reset()
}

func TestWriteTextures(t *testing.T) {
	input := []Texture{
		{
			InternalIdData: []byte("grass"),
			PNGBytes:       grassPngBytes,
		},
		{
			InternalIdData: []byte("stone"),
			PNGBytes:       stonePngBytes,
		},
	}

	expectedBytes := []byte{0x00, 0x00, 0x00, 0x02 /* 2 textures len */}

	for _, tex := range input {
		expectedBytes = append(expectedBytes, uint32ToBytes(len(tex.InternalIdData))...)
		expectedBytes = append(expectedBytes, tex.InternalIdData...)
		expectedBytes = append(expectedBytes, uint32ToBytes(len(tex.PNGBytes))...)
		expectedBytes = append(expectedBytes, tex.PNGBytes...)
	}

	err := WriteTextures(&buf, input)

	if !bytes.Equal(expectedBytes, buf.Bytes()) || err != nil {
		t.Errorf("WriteGRPGTex= %q, %v, want match for %#q", buf.Bytes(), err, expectedBytes)
	}
}

func uint32ToBytes(u int) []byte {
	arr := make([]byte, 4)
	binary.BigEndian.PutUint32(arr, uint32(u))

	return arr
}
