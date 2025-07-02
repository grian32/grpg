package main

import (
	"bytes"
	"fmt"
	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgtex"
	"image"
	"image/png"
	"io"
	"log"
	"os"
)

type GiuTextureTyped struct {
	Texture           *g.Texture
	InternalId        uint16
	FormattedStringId string // this is mainly for perf lol so ur not computing strings in the render loop
	TextureType       grpgtex.TextureType
}

var (
	textures = make(map[string]GiuTextureTyped)
)

func LoadTextures() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	out, err := dialog.File().Title("Please select a textures.pak file").SetStartDir(workingDir).Load()
	if err != nil {
		log.Fatal(err)
	}

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
				Texture:           texture,
				InternalId:        tex.InternalIdInt,
				FormattedStringId: fmt.Sprintf("%s(id: %d)", tex.InternalIdString, tex.InternalIdInt),
				TextureType:       tex.Type,
			}
		})
	}

	BuildSelectorTypeMaps()
}
