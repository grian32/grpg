package main

import (
	"context"
	"fmt"
	"os"
	"texture-packer/textures"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "grpgpack",
		Short: "A data packer for GRPG.",
		RunE:  run,
	}

	texOpts := &textures.TexOptions{}

	texCmd := &cobra.Command{
		Use:   "tex",
		Short: "Packs textures.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return textures.RunTexCmd(cmd, args, texOpts)
		},
	}

	texCmd.Flags().StringVarP(&texOpts.Manifest, "manifest", "m", "", "The path to the texture manifest.")
	texCmd.Flags().StringVarP(&texOpts.Output, "output", "o", "textures.grpgtex", "The output path.")

	cmd.AddCommand(texCmd)

	if err := fang.Execute(context.Background(), cmd); err != nil {
		os.Exit(1)
	}
}

func run(c *cobra.Command, _ []string) error {
	fmt.Println("You've ran the grpgpack root command, this currently has no functionality.")
	return nil
}
