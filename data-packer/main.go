package main

import (
	"context"
	"fmt"
	"os"
	"texture-packer/objs"
	"texture-packer/shared"
	"texture-packer/textures"
	"texture-packer/tiles"

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

	tileOpts := &shared.SharedOptions{}

	tileCmd := &cobra.Command{
		Use:   "tile",
		Short: "Packs tiles.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tiles.RunTileCmd(cmd, args, tileOpts)
		},
	}

	tileCmd.Flags().StringVarP(&tileOpts.Manifest, "manifest", "m", "", "The path to the tile manifest.")
	tileCmd.Flags().StringVarP(&tileOpts.Output, "output", "o", "tiles.grpgtile", "The output path.")
	tileCmd.Flags().StringVarP(&tileOpts.Textures, "textures", "t", "", "The path to the used grpgtex file.")

	cmd.AddCommand(tileCmd)

	objOpts := &shared.SharedOptions{}

	objCmd := &cobra.Command{
		Use:   "obj",
		Short: "Packs objects.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return objs.RunObjCommand(cmd, args, objOpts)
		},
	}

	objCmd.Flags().StringVarP(&objOpts.Manifest, "manifest", "m", "", "The path to the tile manifest.")
	objCmd.Flags().StringVarP(&objOpts.Output, "output", "o", "tiles.grpgtile", "The output path.")
	objCmd.Flags().StringVarP(&objOpts.Textures, "textures", "t", "", "The path to the used grpgtex file.")

	cmd.AddCommand(objCmd)

	if err := fang.Execute(context.Background(), cmd); err != nil {
		os.Exit(1)
	}
}

func run(c *cobra.Command, _ []string) error {
	fmt.Println("You've ran the grpgpack root command, this currently has no functionality.")
	return nil
}
