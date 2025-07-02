package main

import (
	"grpg/data-go/grpgtex"
	"io"
	"log"
	"os"
	"testing"
)

var (
	stonePngBytes []byte
	grassPngBytes []byte
)

func init() {
	var err error

	stone, err := os.Open("./testdata/stone_texture.png")
	if err != nil {
		log.Fatal("Error loading stone texture while initializing format tests")
	}
	grass, err := os.Open("./testdata/grass_texture.png")
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

func TestParseManifestFile(t *testing.T) {
	expected := []GRPGTexManifestEntry{
		{
			InternalName: "grass",
			InternalId:   0,
			FilePath:     "testdata/grass_texture.png",
			Type:         "TILE",
		},
		{
			InternalName: "stone",
			InternalId:   1,
			FilePath:     "testdata/stone_texture.png",
			Type:         "OBJ",
		},
	}

	filePath := "./testdata/test_manifest.toml"

	output, err := ParseManifestFile(filePath)

	// ehh @ comparison
	if len(output) < 2 || output[0] != expected[0] || output[1] != expected[1] || err != nil {
		t.Errorf("ParseManifestFile = %q, %v, want match for %#q", output, err, expected)
	}
}

func TestBuildGRPGTexFromManifest(t *testing.T) {
	manifest := []GRPGTexManifestEntry{
		{
			InternalName: "grass",
			InternalId:   0,
			FilePath:     "testdata/grass_texture.png",
			Type:         "TILE",
		},
		{
			InternalName: "stone",
			InternalId:   1,
			FilePath:     "testdata/stone_texture.png",
			Type:         "OBJ",
		},
	}

	expected := []grpgtex.Texture{
		{
			InternalIdString: []byte("grass"),
			InternalIdInt:    0,
			PNGBytes:         grassPngBytes,
			Type:             grpgtex.TILE,
		},
		{
			InternalIdString: []byte("stone"),
			InternalIdInt:    1,
			PNGBytes:         stonePngBytes,
			Type:             grpgtex.OBJ,
		},
	}

	output, err := BuildGRPGTexFromManifest(manifest)

	if len(output) < 2 || !output[0].Equals(expected[0]) || !output[1].Equals(expected[1]) || err != nil {
		t.Errorf("BuildGRPGTexFromManifest(manifest)= %q, %v, want match for %#q", output, err, expected)
	}
}
