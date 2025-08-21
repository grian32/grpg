package main

import (
	"bytes"
	"fmt"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgtex"
	"grpg/data-go/grpgtile"
	"image"
	"log"
	"os"

	"grpg/data-go/grpgobj"

	g "github.com/AllenDang/giu"
	"github.com/gen2brain/jpegxl"
	"github.com/sqweek/dialog"
)

type GiuTextureTyped struct {
	Texture    *g.Texture
	InternalId uint16
}

var (
	textures          = make(map[int32]GiuTextureTyped)
	tiles             = make(map[int32]grpgtile.Tile)
	objs              = make(map[int32]grpgobj.Obj)
	objsLoaded   bool = false
	tilesLoaded  bool = false
	assetsLoaded bool = false
)

func LoadTextures() {
	buf := loadFileToGBuf("Please select a textures.grpgtex file.")

	header, err := grpgtex.ReadHeader(buf)
	if err != nil {
		fmt.Printf("reading grpgtex header errored: %w.\n", err)
		return
	}
	correctMagic := "GRPGTEX\x00"

	// move this to some notification system or something
	if string(header.Magic[:]) != correctMagic {
		fmt.Println("File entered for texture loading has the wrong magic header.")
		return
	}

	grpgTextures, err := grpgtex.ReadTextures(buf)
	if err != nil {
		fmt.Printf("reading grpgtex textures errored: %w. file: %s\n", err)
		return
	}

	for _, tex := range grpgTextures {
		jxlImage, err := jpegxl.Decode(bytes.NewReader(tex.ImageBytes))
		if err != nil {
			log.Fatal(err)
		}

		g.NewTextureFromRgba(jxlImage.(*image.NRGBA), func(texture *g.Texture) {
			typed := GiuTextureTyped{
				Texture:    texture,
				InternalId: tex.InternalIdInt,
			}
			textures[int32(tex.InternalIdInt)] = typed
		})
	}
}

func LoadTiles() {
	buf := loadFileToGBuf("Please select a tiles.grpgtile file.")

	header, err := grpgtile.ReadHeader(buf)
	if err != nil {
		fmt.Println("error reading grpgtile header\n")
		return
	}

	if header.Magic != [8]byte{'G', 'R', 'P', 'G', 'T', 'I', 'L', 'E'} {
		fmt.Println("magic header for file isn't GRPGTILE")
		return
	}

	grpgTiles, err := grpgtile.ReadTiles(buf)
	if err != nil {
		fmt.Printf("error reading tiles %w\n", err)
	}

	for _, tile := range grpgTiles {
		tiles[int32(tile.TileId)] = tile
	}

	tilesLoaded = true

	if tilesLoaded && objsLoaded {
		assetsLoaded = true
	}
}

func LoadObjs() {
	buf := loadFileToGBuf("Please select a objs.grpgobj file.")

	header, err := grpgobj.ReadHeader(buf)
	if err != nil {
		fmt.Println("error reading grpgobj header\n")
		return
	}

	correctMagic := "GRPGOBJ\x00"
	if string(header.Magic[:]) != correctMagic {
		fmt.Println("magic header for file isn't GRPGOBJ")
		return
	}

	grpgObjs, err := grpgobj.ReadObjs(buf)
	if err != nil {
		fmt.Printf("error reading objs %w\n", err)
	}

	for _, obj := range grpgObjs {
		objs[int32(obj.ObjId)] = obj
	}

	objsLoaded = true

	if tilesLoaded && objsLoaded {
		assetsLoaded = true
	}
}

func loadFileToGBuf(dialogTitle string) *gbuf.GBuf {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	out, err := dialog.File().Title(dialogTitle).SetStartDir(workingDir).Load()
	if err != nil {
		log.Fatal(err)
	}

	fileBytes, err := os.ReadFile(out)
	if err != nil {
		log.Fatal(err)
	}

	buf := gbuf.NewGBuf(fileBytes)
	return buf
}
