package tiles

import (
	"io"
	"os"

	"grpg/data-go/grpgtile"

	"github.com/pelletier/go-toml/v2"
)

type ManifestConfig struct {
	Tiles []GRPGTileManifestEntry `toml:"tiles"`
}

type GRPGTileManifestEntry struct {
	Name    string `toml:"name"`
	Id      uint16 `toml:"id"`
	TexName string `toml:"tex_name"`
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

	return cfg.Tiles, nil
}
