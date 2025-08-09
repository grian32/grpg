package objs

import (
	"grpg/data-go/grpgobj"
	"testing"
)

var manifest = []GRPGObjManifestEntry{
	{
		Name:     "water_s",
		ObjId:    1,
		Flags:    []string{},
		Textures: []string{"still_water"},
	},
	{
		Name:     "i_water",
		ObjId:    2,
		Flags:    []string{"STATE", "INTERACT"},
		Textures: []string{"still_water", "grass_tex"},
	},
}

var texMap = map[string]uint16{
	"grass_tex":   2,
	"still_water": 6,
}

func TestParseManifestFile(t *testing.T) {
	filepath := "../testdata/test_obj_manifest.toml"

	output, err := ParseManifestFile(filepath)

	if len(output) != 2 || !manifest[0].Equal(output[0]) || !manifest[1].Equal(output[1]) || err != nil {
		t.Errorf("ParseManifestFile=%v, %v, want match for %v", output, err.Error(), manifest)
	}
}

func TestBuildGRPGObjFromManifest(t *testing.T) {
	expectedObjs := []grpgobj.Obj{
		{
			Name:     "water_s",
			ObjId:    1,
			Flags:    0,
			Textures: []uint16{6},
		},
		{
			Name:     "i_water",
			ObjId:    2,
			Flags:    grpgobj.ObjFlags(grpgobj.STATE | grpgobj.INTERACT),
			Textures: []uint16{6, 2},
		},
	}

	output, err := BuildGRPGObjFromManifest(manifest, texMap)

	if len(output) != 2 || !output[0].Equal(expectedObjs[0]) || !output[1].Equal(expectedObjs[1]) || err != nil {
		t.Errorf("BuildGRPGObjFromManifest=%v, %v want match for %v", output, err, expectedObjs)
	}
}
