package main

import (
	"errors"
	"grpg/data-go/grpgtex"
	"image/png"
	"io"
	"log"
	"os"
	"strings"
)

type GRPGTexManifestEntry struct {
	InternalName string
	FilePath     string
	Type         string
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
			PNGBytes:         pngBytes,
			Type:             getTextureType(file.Type),
		}

		f.Close()
	}

	return tex, nil
}

func ParseManifestFile(path string) ([]GRPGTexManifestEntry, error) {
	content, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	var lines = strings.Split(string(content), "\n")

	entries := make([]GRPGTexManifestEntry, len(lines))

	for idx, line := range lines {
		var contents = strings.Split(line, "=")

		props := strings.Split(contents[1], "@")

		// eh this is a bit shit but it's an "internal" tool anyway lol
		if !strings.HasSuffix(props[0], ".png") {
			return nil, errors.New("only .png files are allowed as textures")
		}

		entries[idx] = GRPGTexManifestEntry{
			InternalName: contents[0],
			FilePath:     props[0],
			Type:         props[1],
		}
	}

	return entries, nil
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
