package main

import (
	"errors"
	"github.com/pelletier/go-toml/v2"
	"grpg/data-go/grpgtex"
	"image/png"
	"io"
	"os"
)

type ManifestConfig struct {
	Textures []GRPGTexManifestEntry `toml:"texture"`
}

type GRPGTexManifestEntry struct {
	InternalName string `toml:"name"`
	InternalId   int    `toml:"id"`
	FilePath     string `toml:"path"`
	Type         string `toml:"type"`
}

func BuildGRPGTexFromManifest(files []GRPGTexManifestEntry) ([]grpgtex.Texture, error) {
	tex := make([]grpgtex.Texture, len(files))

	for idx, file := range files {
		f, err := os.Open(file.FilePath)
		if err != nil {
			return nil, err
		}

		pngConfig, err := png.DecodeConfig(f)
		if err != nil {
			return nil, err
		}

		if pngConfig.Width != 64 || pngConfig.Height != 64 {
			return nil, errors.New("PNG Images must be exactly 64x64")
		}

		_, err = f.Seek(0, 0)
		if err != nil {
			return nil, err
		}

		pngBytes, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}

		tex[idx] = grpgtex.Texture{
			InternalIdString: []byte(file.InternalName),
			InternalIdInt:    uint16(file.InternalId),
			PNGBytes:         pngBytes,
			Type:             getTextureType(file.Type),
		}

		f.Close()
	}

	return tex, nil
}

func ParseManifestFile(path string) ([]GRPGTexManifestEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg ManifestConfig
	err = toml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg.Textures, nil
}

var textureTypeMap = map[string]grpgtex.TextureType{
	"TILE": grpgtex.TILE,
	"OBJ":  grpgtex.OBJ,
}

func getTextureType(str string) grpgtex.TextureType {
	texType, exists := textureTypeMap[str]

	if !exists {
		return grpgtex.UNDEFINED
	}

	return texType
}
