package main

import (
	"bytes"
	"fmt"
	g "github.com/AllenDang/giu"
	"github.com/wizzymore/tinyfiledialogs"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgtex"
	"image"
	"image/png"
	"io"
	"log"
	"os"
)

type GiuTextureTyped struct {
	Texture     *g.Texture
	TextureType grpgtex.TextureType
}

var (
	textures = make(map[string]GiuTextureTyped)
)

func LoadTextures() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	out, got := tinyfiledialogs.OpenFileDialog("Please select a textures.pak file.", workingDir, []string{"*"}, "any", false)
	if got && out != "" {
		file, err := os.Open(out)
		if err != nil {
			log.Fatal(err)
		}

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}

		buf := gbuf.NewGBuf(fileBytes)
		header := grpgtex.ReadHeader(buf)
		correctMagic := "GRPGTEX\x00"

		// move this to some notification system or something
		if string(header.Magic[:]) != correctMagic {
			log.Fatal("File entered for texture loading has the wrong magic header.")
		} else {
			fmt.Printf("Successfully loaded GRPGTex file with version %d\n", header.Version)
		}

		grpgTextures := grpgtex.ReadTextures(buf)
		for _, tex := range grpgTextures {
			pngImage, err := png.Decode(bytes.NewReader(tex.PNGBytes))
			if err != nil {
				log.Fatal(err)
			}

			internalId := string(tex.InternalIdString[:])

			g.NewTextureFromRgba(pngImage.(*image.NRGBA), func(texture *g.Texture) {
				textures[internalId] = GiuTextureTyped{
					Texture:     texture,
					TextureType: tex.Type,
				}
			})
		}

		BuildSelectorTypeMaps()
	}
}
