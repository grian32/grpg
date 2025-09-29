package util

import "image/color"

var (
	Red    = ValuesRGB(255, 0, 0)
	Yellow = ValuesRGB(255, 255, 0)
)

func ValuesRGB(r, g, b uint8) color.RGBA {
	return color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	}
}
