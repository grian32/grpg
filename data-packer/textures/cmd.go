package textures

import (
	"errors"
	"os"

	"grpg/data-go/gbuf"
	"grpg/data-go/grpgtex"

	"github.com/spf13/cobra"
)

type TexOptions struct {
	Manifest string
	Output   string
}

func RunTexCmd(c *cobra.Command, _ []string, opts *TexOptions) error {
	manifest := opts.Manifest
	output := opts.Output

	if manifest == "" {
		return errors.New("no manifest file provided")
	}

	buf := gbuf.NewEmptyGBuf()

	grpgtex.WriteHeader(buf)

	manifestData, err := ParseManifestFile(manifest)
	if err != nil {
		return err
	}
	textures, err := BuildGRPGTexFromManifest(manifestData)
	if err != nil {
		return err
	}

	grpgtex.WriteTextures(buf, textures)

	f, err := os.Create(output)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(buf.Bytes())

	if err != nil {
		return err
	}

	return nil
}
