package npcs

import (
	"errors"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgnpc"
	"os"
	"texture-packer/shared"

	"github.com/spf13/cobra"
)

func RunNpcCommand(c *cobra.Command, _ []string, opts *shared.SharedOptions) error {
	manifest := opts.Manifest
	output := opts.Output
	texFilePath := opts.Textures

	if manifest == "" {
		return errors.New("no manifest file provided")
	}

	if texFilePath == "" {
		return errors.New("no manifest file provided")
	}

	texMap, err := shared.LoadTexturesToMap(texFilePath)
	if err != nil {
		return err
	}

	manifestData, err := ParseManifestFile(manifest)
	if err != nil {
		return err
	}

	npcs := BuildGRPGNpcFromManifest(manifestData, texMap)

	npcBuf := gbuf.NewEmptyGBuf()

	grpgnpc.WriteHeader(npcBuf)
	grpgnpc.WriteNpcs(npcBuf, npcs)

	outputFile, err := os.Create(output)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	_, err = outputFile.Write(npcBuf.Bytes())

	if err != nil {
		return err
	}

	return nil
}
