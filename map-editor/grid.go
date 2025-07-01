package main

import (
	"fmt"
	"github.com/AllenDang/cimgui-go/imgui"
	"image"
	"image/color"

	g "github.com/AllenDang/giu"
)

var (
	rgba       *image.RGBA
	tex        *g.Texture
	wasDrawing = false
)

func BuildGrid() g.Widget {
	rgba, _ = g.LoadImage("./default_editor.png")

	g.NewTextureFromRgba(rgba, func(t *g.Texture) {
		tex = t
	})

	if tex != nil {
		return g.Custom(func() {
			canvas := g.GetCanvas()
			pos := g.GetCursorScreenPos()
			gridMinX := float32(pos.X)
			gridMinY := float32(pos.Y)
			gridMaxX := float32(pos.X + (16 * 64))
			gridMaxY := float32(pos.Y + (16 * 64))

			for dx := range 16 {
				for dy := range 16 {
					minPt := image.Pt(pos.X+(dx*64), pos.Y+(dy*64))
					maxPt := image.Pt(pos.X+(dx*64)+64, pos.Y+(dy*64)+64)
					canvas.AddImage(tex, minPt, maxPt)
					canvas.AddRect(minPt, maxPt, color.RGBA{0, 0, 0, 255}, 0.0, g.DrawFlagsClosed, 1.0)
				}
			}

			mousePos := imgui.MousePos()

			if mousePos.X >= gridMinX && mousePos.X <= gridMaxX && mousePos.Y >= gridMinY && mousePos.Y <= gridMaxY {
				if g.IsMouseDown(g.MouseButtonLeft) {
					dx := mousePos.X - gridMinX
					dy := mousePos.Y - gridMinY

					gridX := int(dx) / 64
					gridY := int(dy) / 64

					fmt.Println(gridX, gridY)
				}
			}
		})
	}

	return nil
}
