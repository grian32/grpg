package main

import (
	g "github.com/AllenDang/giu"
	"image"
)

var (
	rgba *image.RGBA
	tex  *g.Texture
)

func BuildGrid() []g.Widget {
	rgba, _ = g.LoadImage("grass_texture.png")

	g.EnqueueNewTextureFromRgba(rgba, func(t *g.Texture) {
		tex = t
	})

	var rows []g.Widget

	for _ = range 15 {
		var content []g.Widget

		for _ = range 15 {
			content = append(content, g.ImageButton(tex).Size(64.0, 64.0))
		}

		rows = append(rows, g.Row(content...))
	}

	return rows
}
