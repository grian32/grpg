package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image/png"
	"io"
	"os"
)

type GRPGTexHeader struct {
	Magic   [8]byte
	Version uint16
}

type GRPGTexTexture struct {
	InternalIdData []byte
	PNGBytes       []byte
}

func WriteGRPGTexHeader(buf *bytes.Buffer, version uint16) error {
	header := GRPGTexHeader{
		Magic:   [8]byte{'G', 'R', 'P', 'G', 'T', 'E', 'X', 0},
		Version: version,
	}
	err := binary.Write(buf, binary.BigEndian, header)

	return err
}

func BuildGRPGTexFromManifest(files []GRPGTexManifestEntry) ([]GRPGTexTexture, error) {
	tex := make([]GRPGTexTexture, len(files))

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

		tex[idx] = GRPGTexTexture{
			InternalIdData: []byte(file.InternalName),
			PNGBytes:       pngBytes,
		}

		f.Close()
	}

	return tex, nil
}

func WriteGRPGTex(buf *bytes.Buffer, textures []GRPGTexTexture) error {
	err := binary.Write(buf, binary.BigEndian, uint32(len(textures)))
	if err != nil {
		return err
	}

	// can add length checking for lengths being uint32 if it becomes an issue but that seems very unlikely lol..
	for _, tex := range textures {
		err = binary.Write(buf, binary.BigEndian, uint32(len(tex.InternalIdData)))

		if err != nil {
			return err
		}
		// using buf.write cuz not fixed length lol
		buf.Write(tex.InternalIdData)

		err = binary.Write(buf, binary.BigEndian, uint32(len(tex.PNGBytes)))

		if err != nil {
			return err
		}

		buf.Write(tex.PNGBytes)
	}

	return nil
}
