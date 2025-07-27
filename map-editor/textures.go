package main

import (
	"bytes"
	"fmt"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgtex"
	"image"
	"image/png"
	"io"
	"log"
	"os"

	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
)

type GiuTextureTyped struct {
	Texture           *g.Texture
	InternalIdString  string
	InternalId        uint16
	FormattedStringId string // this is mainly for perf lol so ur not computing strings in the render loop
	TextureType       grpgtex.TextureType
}

var (
	textures = make(map[int32]GiuTextureTyped)
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
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	buf := gbuf.NewGBuf(fileBytes)
	header, err := grpgtex.ReadHeader(buf)
	if err != nil {
		fmt.Printf("reading grpgtex header errored: %w. file: %s\n", err, out)
		return
	}
	correctMagic := "GRPGTEX\x00"

	// move this to some notification system or something
	if string(header.Magic[:]) != correctMagic {
		fmt.Println("File entered for texture loading has the wrong magic header.")
		return
	} else {
		fmt.Printf("Successfully loaded GRPGTex file with version %d\n", header.Version)
	}

	grpgTextures, err := grpgtex.ReadTextures(buf)
	if err != nil {
		fmt.Printf("reading grpgtex textures errored: %w. file: %s\n", err, out)
		return
	}

	for _, tex := range grpgTextures {
		pngImage, err := png.Decode(bytes.NewReader(tex.PNGBytes))
		if err != nil {
			log.Fatal(err)
		}

		internalId := string(tex.InternalIdString[:])

		g.NewTextureFromRgba(pngImage.(*image.NRGBA), func(texture *g.Texture) {
			typed := GiuTextureTyped{
				Texture:           texture,
				InternalIdString:  internalId,
				InternalId:        tex.InternalIdInt,
				FormattedStringId: fmt.Sprintf("%s(id: %d)", internalId, tex.InternalIdInt),
				TextureType:       tex.Type,
			}
			textures[int32(tex.InternalIdInt)] = typed
		})
	}

	BuildSelectorTypeMaps()
}
