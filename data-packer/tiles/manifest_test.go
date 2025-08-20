package tiles

import (
	"grpg/data-go/grpgtile"
	"slices"
	"testing"
)

var manifest = []GRPGTileManifestEntry{
	{
		Name:    "grass",
		Id:      1,
		TexName: "grass_tex",
	},
	{
		Name:    "water",
		Id:      2,
		TexName: "still_water",
	},
}

var texMap = map[string]uint16{
	"grass_tex":   2,
	"still_water": 6,
}

func TestParseManifestFile(t *testing.T) {
	filepath := "../testdata/test_tile_manifest.gcfg"

	output, err := ParseManifestFile(filepath)

	if !slices.Equal(manifest, output) || err != nil {
		t.Errorf("ParseManifestFile=%v, %v, want match for %v", output, err.Error(), manifest)
	}
}

func TestBuildGRPGTileFromManifest(t *testing.T) {
	expectedTiles := []grpgtile.Tile{
		{
			Name:   "grass",
			TileId: 1,
			TexId:  2,
		},
		{
			Name:   "water",
			TileId: 2,
			TexId:  6,
		},
	}

	output := BuildGRPGTileFromManifest(manifest, texMap)

	if !slices.Equal(expectedTiles, output) {
		t.Errorf("BuildGRPGTileFromManifest=%v, want match for %v", output, expectedTiles)
	}
}
