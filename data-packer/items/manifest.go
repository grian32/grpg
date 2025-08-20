package items

import (
	"grpg/data-go/grpgitem"
	"os"

	"github.com/grian32/gcfg"
)

type ManifestConfig struct {
	Items []GRPGItemManifestEntry `gcfg:"Item"`
}

type GRPGItemManifestEntry struct {
	Name    string `gcfg:"name"`
	ItemId  uint16 `gcfg:"id"`
	Texture string `gcfg:"texture"`
}

func BuildGRPGItemFromManifest(entries []GRPGItemManifestEntry, texMap map[string]uint16) []grpgitem.Item {
	itemArr := make([]grpgitem.Item, len(entries))

	for idx, entry := range entries {
		itemArr[idx] = grpgitem.Item{
			ItemId:  entry.ItemId,
			Texture: texMap[entry.Texture],
			Name:    entry.Name,
		}
	}

	return itemArr
}

func ParseManifestFile(path string) ([]GRPGItemManifestEntry, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ManifestConfig
	err = gcfg.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg.Items, nil
}
