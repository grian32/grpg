package main

import (
	"bytes"
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

func TestWriteGRPGTexHeaderVer1(t *testing.T) {
	expectedBytes := []byte{
		'G', 'R', 'P', 'G', 'T', 'E', 'X', 0, // magic
		0x00, 0x01, // ver1
	}

	err := WriteGRPGTexHeader(&buf, 1)
	if !bytes.Equal(expectedBytes, buf.Bytes()) || err != nil {
		t.Errorf("WriteGRPGTexHeader(&buf, 1)= %q, %v, want match for %#q", buf.Bytes(), err, expectedBytes)
	}
	buf.Reset()
}

func TestWriteGRPGTexHeaderVerMax(t *testing.T) {
	expectedBytes := []byte{
		'G', 'R', 'P', 'G', 'T', 'E', 'X', 0, // magic
		0xFF, 0xFF, // ver1
	}

	err := WriteGRPGTexHeader(&buf, 65535)
	if !bytes.Equal(expectedBytes, buf.Bytes()) || err != nil {
		t.Errorf("WriteGRPGTexHeader(&buf, 1)= %q, %v, want match for %#q", buf.Bytes(), err, expectedBytes)
	}
	buf.Reset()
}

func TestBuildGRPGTexFromManifest(t *testing.T) {
	manifest := []GRPGTexManifestEntry{
		{
			InternalName: "grass",
			FilePath:     "testdata/grass_texture.png",
		},
		{
			InternalName: "stone",
			FilePath:     "testdata/stone_texture.png",
		},
	}

	expected := []GRPGTexTexture{
		{
			InternalIdData: []byte("grass"),
			PNGBytes:       grassPngBytes,
		},
		{
			InternalIdData: []byte("stone"),
			PNGBytes:       stonePngBytes,
		},
	}

	output, err := BuildGRPGTexFromManifest(manifest)

	if len(output) < 2 || !output[0].Equals(expected[0]) || !output[1].Equals(expected[1]) || err != nil {
		t.Errorf("BuildGRPGTexFromManifest(manifest)= %q, %v, want match for %#q", output, err, expected)
	}
}
