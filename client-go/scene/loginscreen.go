package scene

import (
	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type LoginScreen struct {
	Font         rl.Font
	ScreenWidth  int32
	ScreenHeight int32
	LoginName    string
}

func (l *LoginScreen) Cleanup() {
}

func (l *LoginScreen) Setup() {
}

func (l *LoginScreen) Loop() {
}

func (l *LoginScreen) Render() {
	rl.ClearBackground(rl.Black)

	halfWidth := l.ScreenWidth / 2
	halfHeight := l.ScreenHeight / 2

	drawTitleText(l, halfWidth, halfHeight)
	rg.SetStyle(rg.TEXTBOX, rg.BASE_COLOR_NORMAL, 1)
	rl.DrawRectangle(halfWidth-200, halfHeight-300, 400, 200, rl.NewColor(186, 109, 22, 255))
	drawEnterNameText(l, halfWidth, halfHeight)
	drawLayout(l, halfWidth, halfHeight)
}

// TODO: generalize this maybe? can pass y offset/size/text/spacing and make it draw centered text or somethingh, just trying to port rn
func drawTitleText(l *LoginScreen, halfWidth, halfHeight int32) {
	text := "GRPG Client"
	var size float32 = 48.0

	calculatedSize := rl.MeasureTextEx(l.Font, text, size, 0.0)

	textPos := rl.Vector2{
		X: float32(halfWidth) - (calculatedSize.X / 2),
		Y: float32(halfHeight) - 375,
	}

	rl.DrawTextEx(l.Font, text, textPos, size, 0, rl.White)
}

func drawEnterNameText(l *LoginScreen, halfWidth, halfHeight int32) {
	text := "Enter Name Below:"
	var size float32 = 24.0

	calculatedSize := rl.MeasureTextEx(l.Font, text, size, 0.4)

	textPos := rl.Vector2{
		X: float32(halfWidth) - (calculatedSize.X / 2),
		Y: float32(halfHeight) - 275,
	}

	rl.DrawTextEx(l.Font, text, textPos, size, 0.4, rl.White)
}

func drawLayout(l *LoginScreen, halfWidth, halfHeight int32) {
	loginTextPos := rl.Rectangle{
		X:      float32(halfWidth) - 100,
		Y:      float32(halfHeight) - 250,
		Width:  200,
		Height: 20,
	}

	rg.TextBox(
		loginTextPos,
		&l.LoginName,
		24,
		true,
	)
}
