package util

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
)

// RGBInt64Color
// This is for use within raygui color.
func RGBInt64Color(r, g, b uint8) int64 {
	rgbaCol := color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	}
	return int64(rl.ColorToInt(rgbaCol))
}
