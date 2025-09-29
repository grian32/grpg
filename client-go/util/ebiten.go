package util

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func SubImage(orig *ebiten.Image, rect image.Rectangle) *ebiten.Image {
	// per doc: The returned value is always *ebiten.Image.
	return orig.SubImage(rect).(*ebiten.Image)
}

func DrawImage(screen *ebiten.Image, img *ebiten.Image, x, y int32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, op)
}
