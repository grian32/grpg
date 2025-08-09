package tiles

import (
	"errors"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgtile"
	"os"
	"texture-packer/shared"

	"github.com/spf13/cobra"
)

func RunTileCmd(c *cobra.Command, _ []string, opts *shared.SharedOptions) error {
	manifest := opts.Manifest
	output := opts.Output
	texFilePath := opts.Textures

	if manifest == "" {
		return errors.New("no manifest file provided")
	}

	if texFilePath == "" {
		return errors.New("no textures file provided")
	}

	texMap, err := shared.LoadTexturesToMap(texFilePath)
	if err != nil {
		return err
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
