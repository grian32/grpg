package objs

import (
	"errors"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgobj"
	"os"
	"texture-packer/shared"

	"github.com/spf13/cobra"
)

func RunObjCommand(c *cobra.Command, _ []string, opts *shared.SharedOptions) error {
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

	objs, err := BuildGRPGObjFromManifest(manifestData, texMap)
	if err != nil {
		return err
	}

	objBuf := gbuf.NewEmptyGBuf()

	grpgobj.WriteHeader(objBuf)
	grpgobj.WriteObjs(objBuf, objs)

	outputFile, err := os.Create(output)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	_, err = outputFile.Write(objBuf.Bytes())

	if err != nil {
		return err
	}

	return nil
}
