package shared

import (
	"cmp"
	"errors"
	"io"
	"os"

	"grpg/data-go/gbuf"
	"grpg/data-go/grpgtex"
)

func LoadTexturesToMap(path string) (map[string]uint16, error) {
	texFile, err1 := os.Open(path)
	texBytes, err2 := io.ReadAll(texFile)

	if err := cmp.Or(err1, err2); err != nil {
		return nil, err
	}

	defer texFile.Close()

	buf := gbuf.NewGBuf(texBytes)

	header, err := grpgtex.ReadHeader(buf)
	if err != nil {
		return nil, err
	}

	if header.Magic != [8]byte{'G', 'R', 'P', 'G', 'T', 'E', 'X', 0x00} {
		return nil, errors.New("textures file inputted is not of GRPGTEX file")
	}

	textures, err := grpgtex.ReadTextures(buf)
	if err != nil {
		return nil, err
	}

	texMap := make(map[string]uint16)

	for _, tex := range textures {
		texMap[string(tex.InternalIdString)] = tex.InternalIdInt
	}

	return texMap, nil
}
