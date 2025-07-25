package main

import (
	"context"
	"errors"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgtex"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var (
	manifest string
	output   string
	version  uint16
)

func main() {
	cmd := &cobra.Command{
		Use:   "texpack",
		Short: "A texture packer for GRPG.",
		RunE:  run,
	}

	cmd.Flags().StringVarP(&manifest, "manifest", "m", "", "The path to the texture manifest.")
	cmd.Flags().StringVarP(&output, "output", "o", "textures.pak", "The output path.")
	cmd.Flags().Uint16Var(&version, "texv", 0, "The version.")

	if err := fang.Execute(context.Background(), cmd); err != nil {
		os.Exit(1)
	}
}

func run(c *cobra.Command, _ []string) error {
	if manifest == "" {
		return errors.New("no manifest file provided")
	}
	if version == 0 {
		return errors.New("either version has not been entered, or is 0, which is invalid")
	}

	buf := gbuf.NewEmptyGBuf()

	grpgtex.WriteHeader(buf, version)

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
