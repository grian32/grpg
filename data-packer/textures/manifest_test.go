package textures

import (
	"grpg/data-go/grpgtex"
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

	stonePngBytes, err = os.ReadFile("../testdata/stone_texture.png")
	if err != nil {
		log.Fatal("Error loading stone png bytes while initializing format tests")
	}
	grassPngBytes, err = os.ReadFile("../testdata/grass_texture.png")
	if err != nil {
		log.Fatal("Error loading grass png bytes while initializing format tests")
	}
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

	filePath := "../testdata/test_tex_manifest.toml"

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
			PNGBytes:         grassPngBytes,
		},
		{
			InternalIdString: []byte("stone"),
			InternalIdInt:    2,
			PNGBytes:         stonePngBytes,
		},
	}

	output, err := BuildGRPGTexFromManifest(manifest)

	if len(output) < 2 || !output[0].Equals(expected[0]) || !output[1].Equals(expected[1]) || err != nil {
		t.Errorf("BuildGRPGTexFromManifest(manifest)= %q, %v, want match for %#q", output, err, expected)
	}
}
