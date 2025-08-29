package npcs

import (
	"grpg/data-go/grpgnpc"
	"slices"
	"testing"
)

var manifest = []GRPGNpcManifestEntry{
	{
		Name:    "corey",
		NpcId:   1,
		Texture: "still_water",
	},
	{
		Name:    "grian",
		NpcId:   2,
		Texture: "grass_tex",
	},
}

var texMap = map[string]uint16{
	"grass_tex":   2,
	"still_water": 6,
}

func TestParseManifestFile(t *testing.T) {
	filepath := "../testdata/test_npc_manifest.gcfg"

	output, err := ParseManifestFile(filepath)

	if !slices.Equal(output, manifest) || err != nil {
		t.Errorf("ParseManifestFile=%v, %v, want match for %v", output, err, manifest)
	}
}

func TestBuildGRPGNpcFromManifest(t *testing.T) {
	expectedNpcs := []grpgnpc.Npc{
		{
			NpcId:     1,
			Name:      "corey",
			TextureId: 6,
		},
		{
			NpcId:     2,
			Name:      "grian",
			TextureId: 2,
		},
	}

	output := BuildGRPGNpcFromManifest(manifest, texMap)

	if !slices.Equal(output, expectedNpcs) {
		t.Errorf("BuildGRPGNpcFromManifest=%v, want match for %v", output, expectedNpcs)
	}
}
