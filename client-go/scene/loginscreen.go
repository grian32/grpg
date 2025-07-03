package scene

import (
	"client/util"
	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type LoginScreen struct {
	Font      rl.Font
	LoginName string
	Game      *Game
}

func (l *LoginScreen) Cleanup() {
	rg.LoadStyleDefault()
}

func (l *LoginScreen) Setup() {
	rg.SetFont(l.Font)
	rg.SetStyle(rg.DEFAULT, rg.TEXT_SIZE, 20)

	bgCol := util.RGBInt64Color(232, 152, 16)
	buttonHoverCol := util.RGBInt64Color(227, 160, 43)
	buttonPresedCol := util.RGBInt64Color(204, 144, 39)

	white := util.RGBInt64Color(255, 255, 255)
	rg.SetStyle(rg.TEXTBOX, rg.BASE_COLOR_PRESSED, bgCol)

	rg.SetStyle(rg.BUTTON, rg.BASE_COLOR_NORMAL, bgCol)
	rg.SetStyle(rg.BUTTON, rg.BASE_COLOR_FOCUSED, buttonHoverCol)
	rg.SetStyle(rg.BUTTON, rg.BASE_COLOR_PRESSED, buttonPresedCol)

	rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_PRESSED, white)
	rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_PRESSED, white)
	rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_NORMAL, white)
	rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_NORMAL, white)
	rg.SetStyle(rg.DEFAULT, rg.BORDER_COLOR_FOCUSED, white)
	rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_FOCUSED, white)

}

func (l *LoginScreen) Loop() {

}

func (l *LoginScreen) Render() {
	rl.ClearBackground(rl.Black)

	halfWidth := l.Game.ScreenWidth / 2
	halfHeight := l.Game.ScreenHeight / 2

	drawTitleText(l, halfWidth, halfHeight)
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
		X:      float32(halfWidth) - 50,
		Y:      float32(halfHeight) - 250,
		Width:  100,
		Height: 30,
	}

	rg.TextBox(
		loginTextPos,
		&l.LoginName,
		8,
		true,
	)

	loginButtonPos := rl.Rectangle{
		X:      float32(halfWidth) - 30,
		Y:      float32(halfHeight) - 200,
		Width:  60,
		Height: 30,
	}

	if rg.Button(loginButtonPos, "Login") {
		l.Game.SceneManager.SwitchTo(&Playground{
			Game: l.Game,
		})
	}
}
