package main

import (
	"image"
	"image/color"

	g "github.com/AllenDang/giu"
)

var (
	rgba *image.RGBA
	tex  *g.Texture
)

func BuildGrid() g.Widget {
	rgba, _ = g.LoadImage("grass_texture.png")

	g.EnqueueNewTextureFromRgba(rgba, func(t *g.Texture) {
		tex = t
	})

	if tex != nil {
		return g.Custom(func() {
			canvas := g.GetCanvas()
			pos := g.GetCursorScreenPos()
			for dx := range 15 {
				for dy := range 15 {
					minPt := image.Pt(pos.X+(dx*64), pos.Y+(dy*64))
					maxPt := image.Pt(pos.X+(dx*64)+64, pos.Y+(dy*64)+64)
					canvas.AddImage(tex, minPt, maxPt)
					canvas.AddRect(minPt, maxPt, color.RGBA{0, 0, 0, 255}, 0.0, g.DrawFlagsClosed, 1.0)
				}
			}
		})
	}
	return nil
}
