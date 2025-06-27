package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
	"os"
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

	buf := bytes.Buffer{}

	err := WriteGRPGTexHeader(&buf, version)
	if err != nil {
		return err
	}

	manifestData, err := ParseManifestFile(manifest)
	if err != nil {
		return err
	}
	textures, err := BuildGRPGTexFromManifest(manifestData)
	if err != nil {
		return err
	}

	err = WriteGRPGTex(&buf, textures)
	if err != nil {
		return err
	}

	f, err := os.Create(output)

	defer f.Close()

	if err != nil {
		return err
	}

	_, err = f.Write(buf.Bytes())

	if err != nil {
		return err
	}

	return nil
}
