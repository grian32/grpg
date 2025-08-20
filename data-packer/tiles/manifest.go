package tiles

import (
	"os"

	"grpg/data-go/grpgtile"

	"github.com/grian32/gcfg"
)

type ManifestConfig struct {
	Tiles []GRPGTileManifestEntry `gcfg:"Tile"`
}

type GRPGTileManifestEntry struct {
	Name    string `gcfg:"name"`
	Id      uint16 `gcfg:"id"`
	TexName string `gcfg:"tex_name"`
}

func BuildGRPGTileFromManifest(entries []GRPGTileManifestEntry, texMap map[string]uint16) []grpgtile.Tile {
	tileArr := make([]grpgtile.Tile, len(entries))

	for idx, entry := range entries {
		tileArr[idx] = grpgtile.Tile{
			Name:   entry.Name,
			TileId: entry.Id,
			TexId:  texMap[entry.TexName],
		}
	}

	return tileArr
}

func ParseManifestFile(path string) ([]GRPGTileManifestEntry, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ManifestConfig
	err = gcfg.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg.Tiles, nil
}
