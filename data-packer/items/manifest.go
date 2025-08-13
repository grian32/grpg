package items

import (
	"grpg/data-go/grpgitem"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type ManifestConfig struct {
	Items []GRPGItemManifestEntry `toml:"item"`
}

type GRPGItemManifestEntry struct {
	Name    string `toml:"name"`
	ItemId  uint16 `toml:"id"`
	Texture string `toml:"texture"`
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

	return cfg.Items, nil
}
