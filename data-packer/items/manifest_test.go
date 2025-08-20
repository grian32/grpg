package items

import (
	"grpg/data-go/grpgitem"
	"slices"
	"testing"
)

var (
	manifest = []GRPGItemManifestEntry{
		{
			Name:    "berries",
			ItemId:  1,
			Texture: "grass_tex",
		},
		{
			Name:    "berriess",
			ItemId:  2,
			Texture: "grass_tex",
		},
	}
	texMap = map[string]uint16{
		"grass_tex": 1,
	}
)

func TestParseManifestFile(t *testing.T) {
	filepath := "../testdata/test_item_manifest.gcfg"

	output, err := ParseManifestFile(filepath)

	if !slices.Equal(output, manifest) || err != nil {
		t.Errorf("ParseManifestFile=%v, %v, wanted match for %v", output, err, manifest)
	}
}

func TestBuildGRPGItemFromManifest(t *testing.T) {
	expectedItems := []grpgitem.Item{
		{
			ItemId:  1,
			Texture: 1,
			Name:    "berries",
		},
		{
			ItemId:  2,
			Texture: 1,
			Name:    "berriess",
		},
	}

	output := BuildGRPGItemFromManifest(manifest, texMap)

	if !slices.Equal(expectedItems, output) {
		t.Errorf("BuildGRPGItemFromManifest=%v, want match for %v", output, expectedItems)
	}
}
