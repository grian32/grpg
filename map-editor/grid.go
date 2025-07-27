package main

import (
	"fmt"
	"grpg/data-go/grpgtex"
	"image"
	"image/color"

	"github.com/AllenDang/cimgui-go/imgui"

	g "github.com/AllenDang/giu"
)

var (
	defaultTex       *g.Texture
	gridTileTextures = [256]int32{}
	gridObjTextures  = [256]int32{}
)

func LoadDefaultGridTex() {
	for idx := range 256 {
		gridTileTextures[idx] = -1
		gridObjTextures[idx] = -1
	}

	rgba, _ := g.LoadImage("./default_editor.png")

	g.NewTextureFromRgba(rgba, func(t *g.Texture) {
		defaultTex = t
	})
}

func BuildGrid() g.Widget {
	return g.Custom(func() {
		canvas := g.GetCanvas()
		pos := g.GetCursorScreenPos()
		gridMinX := float32(pos.X)
		gridMinY := float32(pos.Y)
		gridMaxX := float32(pos.X+(16*64)) - 2
		gridMaxY := float32(pos.Y+(16*64)) - 2

		for dx := range 16 {
			for dy := range 16 {
				minPt := image.Pt(pos.X+(dx*64), pos.Y+(dy*64))
				maxPt := image.Pt(pos.X+(dx*64)+64, pos.Y+(dy*64)+64)

				texTileName := gridTileTextures[dx+dy*16]
				tileTex := textures[texTileName].Texture
				if texTileName == -1 {
					tileTex = defaultTex
				}

				texObjName := gridObjTextures[dx+dy*16]
				objTex, objOk := textures[texObjName]

				canvas.AddImage(tileTex, minPt, maxPt)
				canvas.AddRect(minPt, maxPt, color.RGBA{0, 0, 0, 255}, 0.0, g.DrawFlagsClosed, 1.0)
				if objOk {
					canvas.AddImage(objTex.Texture, minPt, maxPt)
				}
			}
		}

		mousePos := imgui.MousePos()

		if textureSelected && mousePos.X >= gridMinX && mousePos.X <= gridMaxX && mousePos.Y >= gridMinY && mousePos.Y <= gridMaxY {
			if g.IsMouseDown(g.MouseButtonLeft) {
				dx := mousePos.X - gridMinX
				dy := mousePos.Y - gridMinY

				gridX := int(dx) / 64
				gridY := int(dy) / 64

				currPos := gridX + gridY*16

				if eraserEnabled {
					if gridObjTextures[currPos] != -1 {
						gridObjTextures[currPos] = -1
					} else if gridTileTextures[currPos] != -1 {
						gridTileTextures[currPos] = -1
					}
				} else {
					switch currentlySelected.val.TextureType {
					case grpgtex.OBJ:
						gridObjTextures[currPos] = int32(currentlySelected.val.InternalId)
					case grpgtex.TILE:
						gridTileTextures[currPos] = int32(currentlySelected.val.InternalId)
					case grpgtex.UNDEFINED:
					default:
						panic(fmt.Sprintf("unexpected grpgtex.TextureType: %#v", currentlySelected.val.TextureType))
					}
				}
			}
		}
	})
}
