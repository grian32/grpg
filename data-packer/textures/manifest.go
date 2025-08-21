package textures

import (
	"bytes"
	"errors"
	"grpg/data-go/grpgtex"
	"image/png"
	"log"
	"os"

	"github.com/gen2brain/jpegxl"
	"github.com/grian32/gcfg"
)

type ManifestConfig struct {
	Textures []GRPGTexManifestEntry `gcfg:"Texture"`
}

type GRPGTexManifestEntry struct {
	InternalName string `gcfg:"name"`
	InternalId   int    `gcfg:"id"`
	FilePath     string `gcfg:"path"`
}

func BuildGRPGTexFromManifest(files []GRPGTexManifestEntry) ([]grpgtex.Texture, error) {
	tex := make([]grpgtex.Texture, len(files))

	for idx, file := range files {
		f, err := os.Open(file.FilePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		_, err = f.Seek(0, 0)
		if err != nil {
			return nil, err
		}

		image, err := png.Decode(f)
		if err != nil {
			return nil, errors.New("failed to decode png image")
		}

		var jpegXlBuf bytes.Buffer

		jpegXlOptions := jpegxl.Options{
			Quality: 100,
			Effort:  10,
		}

		err = jpegxl.Encode(&jpegXlBuf, image, jpegXlOptions)
		if err != nil {
			return nil, err
		}

		if file.InternalId == 0 {
			log.Fatalln("Integer ID 0 is reserved.")
		}

		tex[idx] = grpgtex.Texture{
			InternalIdString: []byte(file.InternalName),
			InternalIdInt:    uint16(file.InternalId),
			ImageBytes:       jpegXlBuf.Bytes(),
		}

		f.Close()
	}

	return tex, nil
}

func ParseManifestFile(path string) ([]GRPGTexManifestEntry, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ManifestConfig
	err = gcfg.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg.Textures, nil
}
