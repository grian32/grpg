package tiles

import (
	"cmp"
	"errors"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgtex"
	"grpg/data-go/grpgtile"
	"io"
	"os"

	"github.com/spf13/cobra"
)

type TilesOptions struct {
	Manifest string
	Output   string
	Textures string
}

func RunTileCmd(c *cobra.Command, _ []string, opts *TilesOptions) error {
	manifest := opts.Manifest
	output := opts.Output
	texFilePath := opts.Textures

	if manifest == "" {
		return errors.New("no manifest file provided")
	}

	if texFilePath == "" {
		return errors.New("no textures file provided")
	}

	texFile, err1 := os.Open(texFilePath)
	texBytes, err2 := io.ReadAll(texFile)

	if err := cmp.Or(err1, err2); err != nil {
		return err
	}

	defer texFile.Close()

	buf := gbuf.NewGBuf(texBytes)

	header, err := grpgtex.ReadHeader(buf)
	if err != nil {
		return err
	}

	if header.Magic != [8]byte{'G', 'R', 'P', 'G', 'T', 'E', 'X', 0x00} {
		return errors.New("textures file inputted is not of GRPGTEX file")
	}

	textures, err := grpgtex.ReadTextures(buf)
	if err != nil {
		return err
	}

	texMap := make(map[string]uint16)

	for _, tex := range textures {
		texMap[string(tex.InternalIdString)] = tex.InternalIdInt
	}

	manifestData, err := ParseManifestFile(manifest)
	if err != nil {
		return err
	}

	tiles := BuildGRPGTileFromManifest(manifestData, texMap)

	tileBuf := gbuf.NewEmptyGBuf()

	grpgtile.WriteHeader(tileBuf)
	grpgtile.WriteTiles(tileBuf, tiles)

	outputFile, err := os.Create(output)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	_, err = outputFile.Write(tileBuf.Bytes())

	if err != nil {
		return err
	}

	return nil
}
