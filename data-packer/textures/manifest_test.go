package textures

import (
	"bytes"
	"grpg/data-go/grpgtex"
	"image/png"
	"log"
	"os"
	"testing"

	"github.com/gen2brain/jpegxl"
)

var (
	stoneJxlBytes []byte
	grassJxlBytes []byte
)

func init() {
	var err error

	stonePng, err := os.Open("../testdata/stone_texture.png")
	if err != nil {
		log.Fatal(err)
	}
	grassPng, err := os.Open("../testdata/grass_texture.png")
	if err != nil {
		log.Fatal(err)
	}

	stoneImg, err := png.Decode(stonePng)
	if err != nil {
		log.Fatal(err)
	}
	grassImg, err := png.Decode(grassPng)
	if err != nil {
		log.Fatal(err)
	}

	jxlOptions := jpegxl.Options{
		Quality: 100,
		Effort:  10,
	}

	var stoneJxlBuf bytes.Buffer
	err = jpegxl.Encode(&stoneJxlBuf, stoneImg, jxlOptions)
	stoneJxlBytes = stoneJxlBuf.Bytes()

	var grassJxlBuf bytes.Buffer
	err = jpegxl.Encode(&grassJxlBuf, grassImg, jxlOptions)
	grassJxlBytes = grassJxlBuf.Bytes()

	file, err := os.Create("grass.jxl")
	file.Write(grassJxlBytes)
	file, err = os.Create("stone.jxl")
	file.Write(stoneJxlBytes)
}

func TestParseManifestFile(t *testing.T) {
	expected := []GRPGTexManifestEntry{
		{
			InternalName: "grass_tex",
			InternalId:   1,
			FilePath:     "testdata/grass_texture.png",
		},
		{
			InternalName: "still_water",
			InternalId:   2,
			FilePath:     "testdata/stone_texture.png",
		},
	}

	filePath := "../testdata/test_tex_manifest.gcfg"

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
			InternalId:   1,
			FilePath:     "../testdata/grass_texture.png",
		},
		{
			InternalName: "stone",
			InternalId:   2,
			FilePath:     "../testdata/stone_texture.png",
		},
	}

	expected := []grpgtex.Texture{
		{
			InternalIdString: []byte("grass"),
			InternalIdInt:    1,
			ImageBytes:       grassJxlBytes,
		},
		{
			InternalIdString: []byte("stone"),
			InternalIdInt:    2,
			ImageBytes:       stoneJxlBytes,
		},
	}

	output, err := BuildGRPGTexFromManifest(manifest)

	if len(output) < 2 || !output[0].Equals(expected[0]) || !output[1].Equals(expected[1]) || err != nil {
		t.Errorf("BuildGRPGTexFromManifest(manifest)= %q, %v, want match for %#q", output, err, expected)
	}
}
