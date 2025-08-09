package objs

import (
	"errors"
	"io"
	"os"
	"slices"

	"grpg/data-go/grpgobj"

	"github.com/pelletier/go-toml/v2"
)

type ManifestConfig struct {
	Objs []GRPGObjManifestEntry `toml:"obj"`
}

type GRPGObjManifestEntry struct {
	Name         string   `toml:"name"`
	ObjId        uint16   `toml:"id"`
	Flags        []string `toml:"flags"`
	Textures     []string `toml:"textures"`
	InteractText string   `toml:"interact_text"`
}

// Equal only meant to be used for testing, you probably shouldn't be == this type
func (om *GRPGObjManifestEntry) Equal(other GRPGObjManifestEntry) bool {
	return om.Name == other.Name && om.ObjId == other.ObjId && slices.Equal(om.Flags, other.Flags) && slices.Equal(om.Textures, other.Textures)
}

func BuildGRPGObjFromManifest(entries []GRPGObjManifestEntry, texMap map[string]uint16) ([]grpgobj.Obj, error) {
	objArr := make([]grpgobj.Obj, len(entries))

	for idx, entry := range entries {
		flags := flagsFromStringSlice(entry.Flags)

		if !grpgobj.IsFlagSet(flags, grpgobj.INTERACT) && entry.InteractText != "" {
			return nil, errors.New("interact_text without interact flag not allowed")
		}

		texArr := make([]uint16, len(entry.Textures))

		for texIdx, tex := range entry.Textures {
			texArr[texIdx] = texMap[tex]
		}

		objArr[idx] = grpgobj.Obj{
			Name:         entry.Name,
			ObjId:        entry.ObjId,
			Flags:        flags,
			Textures:     texArr,
			InteractText: entry.InteractText,
		}
	}

	return objArr, nil
}

func flagsFromStringSlice(flags []string) grpgobj.ObjFlags {
	var data grpgobj.ObjFlags = 0

	for _, flag := range flags {
		switch flag {
		case "STATE":
			data |= grpgobj.ObjFlags(grpgobj.STATE)
		case "INTERACT":
			data |= grpgobj.ObjFlags(grpgobj.INTERACT)
		}
	}

	return data
}

func ParseManifestFile(path string) ([]GRPGObjManifestEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg ManifestConfig
	err = toml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg.Objs, nil
}
